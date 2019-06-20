package security

import (
	"errors"

	"gitlab.skypicker.com/platform/security/iam/security/secrets"
)

var (
	errUnathorised = errors.New("incorrect token")
)

// VerifyToken accepts a token and a service struct and verifies if this token is accepted
func VerifyToken(secretManager secrets.SecretManager, service Service, requestToken string) error {
	token, err := secretManager.GetAppToken(service.Name, service.Environment)

	if err != nil {
		return err
	}

	if token != requestToken {
		return errUnathorised
	}

	return nil
}