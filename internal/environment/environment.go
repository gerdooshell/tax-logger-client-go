package environment

import (
	"fmt"
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
	Dev  Environment = "Dev"
	Prod Environment = "Prod"
)
