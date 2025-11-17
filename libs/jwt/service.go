package jwt

import (
    "context"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "time"

    "github.com/lestrrat-go/jwx/v3/jwa"
    "github.com/lestrrat-go/jwx/v3/jwk"
    "github.com/lestrrat-go/jwx/v3/jwt"
)

// Service is the public interface used by auth-service
type Service interface {
    // Access tokens
    GenerateAccessToken(ctx context.Context, u *UserInfo) (string, error)
    ValidateAccessToken(ctx context.Context, tokenStr string) (*Claims, error)

    // Refresh tokens (opaque)
    NewOpaqueToken(nBytes int) (string, error)
    HashRefreshToken(plain string) (string, error)
    GenerateRefreshToken(ctx context.Context, u *UserInfo) (plain string, hashed string, expiresAt time.Time, err error)

    // JWKS
    JWKS() (jwk.Set, error)
}

// rsaService implements Service with RS256 signing
type rsaService struct {
    cfg    Config
    priv   *rsa.PrivateKey
    pub    *rsa.PublicKey
    jwkKey jwk.Key
}

// NewRSAService constructs the service. privKey must be non-nil and KeyID provided.
func NewRSAService(cfg Config, priv *rsa.PrivateKey, pub *rsa.PublicKey) (Service, error) {
    if priv == nil || pub == nil {
        return nil, errors.New("private/public key required")
    }
    if cfg.KeyID == "" {
        return nil, errors.New("cfg.KeyID required")
    }

    // build jwk.Key representation of public key with kid, alg, use
    k, err := jwk.New(pub)
    if err != nil {
        return nil, err
    }
    _ = k.Set(jwk.KeyIDKey, cfg.KeyID)
    _ = k.Set(jwk.AlgorithmKey, jwa.RS256.String())
    _ = k.Set(jwk.KeyUsageKey, "sig")

    return &rsaService{
        cfg:    cfg,
        priv:   priv,
        pub:    pub,
        jwkKey: k,
    }, nil
}

// GenerateAccessToken creates a signed JWT (RS256) and returns compact serialization.
func (s *rsaService) GenerateAccessToken(ctx context.Context, u *UserInfo) (string, error) {
    now := time.Now().UTC()
    tok, err := jwt.NewBuilder().
        Issuer(s.cfg.Issuer).
        Audience(s.cfg.Audience).
        IssuedAt(now).
        Expiration(now.Add(s.cfg.AccessTTL)).
        Claim("uid", u.UserID).
        Claim("email", u.Email).
        Claim("role", u.Role).
        Build()
    if err != nil {
        return "", err
    }

    // Sign token with private key. Using jwk by wrapping the private key as jwk.Key
    signKey, err := jwk.New(s.priv)
    if err != nil {
        return "", err
    }
    _ = signKey.Set(jwk.KeyIDKey, s.cfg.KeyID)          // set kid on the JWK used to sign
    _ = signKey.Set(jwk.AlgorithmKey, jwa.RS256.String()) // set alg on key

    // jwt.Sign will add the signature and produce compact JWS string
    signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256, signKey))
    if err != nil {
        return "", err
    }
    return string(signed), nil
}

// ValidateAccessToken parses and verifies a token string and returns Claims.
func (s *rsaService) ValidateAccessToken(ctx context.Context, tokenStr string) (*Claims, error) {
    // parse + verify using a jwk.Set that contains our public key
    set := jwk.NewSet()
    set.Add(s.jwkKey)

    // jwt.Parse will verify signature using the keys in the set (select by kid)
    tok, err := jwt.Parse([]byte(tokenStr), jwt.WithKeySet(set), jwt.WithValidate(true))
    if err != nil {
        // classify some errors
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, errors.New("token expired")
        }
        return nil, errors.New("invalid token")
    }

    // extract claims safely
    c := Claims{}
    if v, ok := tok.Get("uid"); ok {
        if id, ok2 := v.(float64); ok2 {
            c.UserID = uint64(id)
        }
    }
    if v, ok := tok.Get("email"); ok {
        if sEmail, ok2 := v.(string); ok2 {
            c.Email = sEmail
        }
    }
    if v, ok := tok.Get("role"); ok {
        if sRole, ok2 := v.(string); ok2 {
            c.Role = sRole
        }
    }
    return &c, nil
}

// opaque refresh helpers

func (s *rsaService) NewOpaqueToken(nBytes int) (string, error) {
    if nBytes <= 0 {
        nBytes = 32
    }
    b := make([]byte, nBytes)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *rsaService) HashRefreshToken(plain string) (string, error) {
    if plain == "" {
        return "", errors.New("empty token")
    }
    h := sha256.Sum256([]byte(plain))
    return base64.RawURLEncoding.EncodeToString(h[:]), nil
}

func (s *rsaService) GenerateRefreshToken(ctx context.Context, u *UserInfo) (plain string, hashed string, expiresAt time.Time, err error) {
    plain, err = s.NewOpaqueToken(32)
    if err != nil {
        return "", "", time.Time{}, err
    }
    hashed, err = s.HashRefreshToken(plain)
    if err != nil {
        return "", "", time.Time{}, err
    }
    expiresAt = time.Now().UTC().Add(s.cfg.RefreshTTL)
    return plain, hashed, expiresAt, nil
}

// JWKS returns a jwk.Set containing the public key (ready to JSON marshal)
func (s *rsaService) JWKS() (jwk.Set, error) {
    set := jwk.NewSet()
    set.Add(s.jwkKey)
    return set, nil
}
