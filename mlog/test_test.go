package mlog_test

import (
	"testing"

	"github.com/chu16537/module_master/mlog"
)

func Test(t *testing.T) {
	config := mlog.Config{
		Name:  "test",
		Level: "debug",
	}

	err := mlog.Init(&config)
	if err != nil {
		t.Errorf("%+v", err)
		return
	}

	mlog.Info("test")

	mlog.Error("test")

	defer func() {
		if r := recover(); r != nil {
			mlog.Fatal(r)
		}
	}()

	panic("aaa")
}
