package mlog

import "fmt"

func (h *Handler) Log(level int, msg interface{}) {
	switch level {
	case Level_Debug:
		h.Debug(msg)
	case Level_Info:
		h.Info(msg)
	case Level_Warn:
		h.Warn(msg)
	case Level_Error:
		h.Error(msg)
	default:
		h.Debug(msg)
	}
}

func (h *Handler) Debug(msg interface{}) {
	fmt.Println("Debug", msg)
}

func (h *Handler) Info(msg interface{}) {
	fmt.Println("Info", msg)
}

func (h *Handler) Warn(msg interface{}) {
	fmt.Println("Warn", msg)
}

func (h *Handler) Error(msg interface{}) {
	fmt.Println("Error", msg)
}
