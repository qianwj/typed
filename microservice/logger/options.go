package logger

type Conf struct {
	Output *struct {
		Type string
		Path string
	}
	Threshold  int
	StackTrace int
}
