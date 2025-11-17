package jwt

import (
    "context"
    "crypto/rsa"
    "errors"
    "net/http"
    "strings"
    "sync"
    "time"

    jwxjwt "github.com/lestrrat-go/jwx/v3/jwt"
    "github.com/lestrrat-go/jwx/v3/jwk"
)

// Verifier verifies tokens for resource servers.
type Verifier interface {
    Verify(ctx context.Context, tokenStr string) (*Claims, error)
}

var (
    ErrUnauthorized = errors.New("unauthorized")
    ErrExpired      = errors.New("token expired")
)

// Local verifier: uses a single public key
type localVerifier struct {
    pub      *rsa.PublicKey
    issuer   string
    audience string
}

func NewLocalVerifier(pub *rsa.PublicKey, issuer, audience string) Verifier {
    return &localVerifier{pub: pub, issuer: issuer, audience: audience}
}

func (v *localVerifier) Verify(ctx context.Context, tokenStr string) (*Claims, error) {
    set := jwk.NewSet()
    k, err := jwk.New(v.pub)
    if err == nil {
        _ = k.Set(jwk.KeyUsageKey, "sig")
        set.Add(k)
    }
    opts := []jwxjwt.ParseOption{jwxjwt.WithKeySet(set), jwxjwt.WithValidate(true)}
    if v.issuer != "" {
        opts = append(opts, jwxjwt.WithIssuer(v.issuer))
    }
    if v.audience != "" {
        opts = append(opts, jwxjwt.WithAudience(v.audience))
    }
    tok, err := jwxjwt.Parse([]byte(tokenStr), opts...)
    if err != nil {
        if errors.Is(err, jwxjwt.ErrTokenExpired) {
            return nil, ErrExpired
        }
        return nil, ErrUnauthorized
    }
    return extractClaimsFromToken(tok)
}

// Remote verifier: fetches JWKS and refreshes periodically
type remoteVerifier struct {
    jwksURL        string
    issuer         string
    audience       string
    refreshInterval time.Duration

    mu     sync.RWMutex
    keySet jwk.Set

    httpClient *http.Client
    quit       chan struct{}
}

func NewRemoteVerifier(jwksURL, issuer, audience string, refreshInterval time.Duration) (Verifier, error) {
    if refreshInterval <= 0 {
        refreshInterval = 5 * time.Minute
    }
    rv := &remoteVerifier{
        jwksURL:         jwksURL,
        issuer:          issuer,
        audience:         audience,
        refreshInterval: refreshInterval,
        httpClient:      &http.Client{Timeout: 10 * time.Second},
        quit:            make(chan struct{}),
    }
    if err := rv.fetchOnce(); err != nil {
        return nil, err
    }
    go rv.refreshLoop()
    return rv, nil
}

func (r *remoteVerifier) fetchOnce() error {
    set, err := jwk.Fetch(context.Background(), r.jwksURL)
    if err != nil {
        return err
    }
    r.mu.Lock()
    r.keySet = set
    r.mu.Unlock()
    return nil
}

func (r *remoteVerifier) refreshLoop() {
    ticker := time.NewTicker(r.refreshInterval)
    defer ticker.Stop()
    for {
        select {
        case <-ticker.C:
            _ = r.fetchOnce()
        case <-r.quit:
            return
        }
    }
}

func (r *remoteVerifier) Close() {
    close(r.quit)
}

func (r *remoteVerifier) getSet() jwk.Set {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return r.keySet
}

func (r *remoteVerifier) Verify(ctx context.Context, tokenStr string) (*Claims, error) {
    set := r.getSet()
    if set.Len() == 0 {
        return nil, ErrUnauthorized
    }
    opts := []jwxjwt.ParseOption{jwxjwt.WithKeySet(set), jwxjwt.WithValidate(true)}
    if r.issuer != "" {
        opts = append(opts, jwxjwt.WithIssuer(r.issuer))
    }
    if r.audience != "" {
        opts = append(opts, jwxjwt.WithAudience(r.audience))
    }
    tok, err := jwxjwt.Parse([]byte(tokenStr), opts...)
    if err != nil {
        if errors.Is(err, jwxjwt.ErrTokenExpired) {
            return nil, ErrExpired
        }
        if strings.Contains(err.Error(), "no key") || strings.Contains(err.Error(), "kid") {
            // try immediate refetch once
            _ = r.fetchOnce()
            set = r.getSet()
            tok2, err2 := jwxjwt.Parse([]byte(tokenStr), append(opts[:0:0], jwxjwt.WithKeySet(set), jwxjwt.WithValidate(true))...)
            if err2 == nil {
                return extractClaimsFromToken(tok2)
            }
        }
        return nil, ErrUnauthorized
    }
    return extractClaimsFromToken(tok)
}

func extractClaimsFromToken(tok jwxjwt.Token) (*Claims, error) {
    c := Claims{}
    if v, ok := tok.Get("uid"); ok {
        switch id := v.(type) {
        case float64:
            c.UserID = uint64(id)
        case int:
            c.UserID = uint64(id)
        case int64:
            c.UserID = uint64(id)
        case uint64:
            c.UserID = id
        }
    }
    if v, ok := tok.Get("email"); ok {
        if s, ok2 := v.(string); ok2 {
            c.Email = s
        }
    }
    if v, ok := tok.Get("role"); ok {
        if s, ok2 := v.(string); ok2 {
            c.Role = s
        }
    }
    return &c, nil
}
