package security

// SecretManager is a interface that describes how we want to use secrets
type SecretManager interface {
	GetAppToken(string) (string, error)
	GetSetting(string) (string, error)
}
