package mlog

type ILog interface {
	Log(level int, msg interface{})
	Debug(msg interface{})
	Info(msg interface{})
	Warn(msg interface{})
	Error(msg interface{})
}

type Handler struct {
	Name string
}

func New(name string) ILog {
	return &Handler{
		Name: name,
	}
}
