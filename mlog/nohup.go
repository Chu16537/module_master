package mlog

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func InitNohup(name string) error {
	_, err := os.OpenFile(fmt.Sprintf("./%s_nohup.log", name), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	return nil
}

func NohupLog(t time.Time, s string) {
	f := logrus.Fields{}
	f["time"] = t
	f["msg"] = s
	fmt.Println(f)
}
