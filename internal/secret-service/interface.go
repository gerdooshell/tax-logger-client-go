package secretService

import "context"

type SecretOut struct {
	Value string
	Err   error
}

type SecretService interface {
	GetSecretValue(ctx context.Context, secretKey string) <-chan SecretOut
}
