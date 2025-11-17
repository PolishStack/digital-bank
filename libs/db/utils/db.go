package utils

import (
	"fmt"
	"net/url"
	"strings"
)

func SafeEncodeDSN(raw string) (string, error) {
	// Split before '@' → separates userinfo from host
	beforeAt, afterAt, found := strings.Cut(raw, "@")
	if !found {
		return "", fmt.Errorf("invalid DSN: missing '@'")
	}

	// Split after scheme://
	schemePart, userInfoPart, found := strings.Cut(beforeAt, "://")
	if !found {
		return "", fmt.Errorf("invalid DSN: missing scheme")
	}

	// Encode user and password
	user, pass, found := strings.Cut(userInfoPart, ":")
	if !found {
		pass = ""
	}
	userEnc := url.PathEscape(user)
	passEnc := url.PathEscape(pass)

	encoded := fmt.Sprintf("%s://%s:%s@%s", schemePart, userEnc, passEnc, afterAt)
	return encoded, nil
}
