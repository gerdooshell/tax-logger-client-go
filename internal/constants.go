package internal

import "fmt"

type Severity string

const (
	Debug   Severity = "debug"
	Info    Severity = "info"
	Warning Severity = "warning"
	Error   Severity = "error"
	Fatal   Severity = "fatal"
)

type Environment string

func (env Environment) ToString() string {
	return string(env)
}

func GetEnvironmentFromString(envStr string) (env Environment, err error) {
	switch envStr {
	case Prod.ToString():
		env = Prod
	case Dev.ToString():
		env = Dev
	default:
		err = fmt.Errorf("invalid environment \"%s\"", envStr)
	}
	return
}

const (
	Dev  Environment = "dev"
	Prod Environment = "prod"
)
