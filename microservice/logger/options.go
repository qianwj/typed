package logger

type Type int

const (
	Std Type = iota
	Zap
)

type Conf struct {
	Output *struct {
		Type string
		Path string
	}
	Threshold  int
	StackTrace int
}

func Bootstrap(loggerType Type, conf *Conf) {
	switch loggerType {
	case Zap:

	}
}
