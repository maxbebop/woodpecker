package socketmodewrap

type Logger interface {
	Printf(f string, p ...any)
}

type Log struct {
	Logger
}

func (l Log) Output(level int, msg string) error {
	l.Logger.Printf("%d %s", level, msg)

	return nil
}
