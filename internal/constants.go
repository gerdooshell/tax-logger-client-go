package internal

type Severity string

const (
	Debug   Severity = "debug"
	Info    Severity = "info"
	Warning Severity = "warning"
	Error   Severity = "error"
	Fatal   Severity = "fatal"
)
