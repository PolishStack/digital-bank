package jwt

import (
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "fmt"
)

// ParseRSAPrivateKeyPEM parses a PEM encoded private key (PKCS#1 or PKCS#8).
func ParseRSAPrivateKeyPEM(pemBytes []byte) (*rsa.PrivateKey, error) {
    block, _ := pem.Decode(pemBytes)
    if block == nil {
        return nil, fmt.Errorf("invalid PEM")
    }
    // try PKCS1
    if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
        return key, nil
    }
    // try PKCS8
    if parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
        if k, ok := parsed.(*rsa.PrivateKey); ok {
            return k, nil
        }
    }
    return nil, fmt.Errorf("unsupported private key format")
}

// ParseRSAPublicKeyPEM parses a PEM encoded public key (X.509 / PKIX).
func ParseRSAPublicKeyPEM(pemBytes []byte) (*rsa.PublicKey, error) {
    block, _ := pem.Decode(pemBytes)
    if block == nil {
        return nil, fmt.Errorf("invalid PEM")
    }
    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err == nil {
        if p, ok := pub.(*rsa.PublicKey); ok {
            return p, nil
        }
    }
    // maybe it's a PKCS1 public key bytes
    if pk, err2 := x509.ParsePKCS1PublicKey(block.Bytes); err2 == nil {
        return pk, nil
    }
    return nil, fmt.Errorf("unsupported public key format: %v", err)
}
