package main

import (
	"fmt"
	"time"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mlog"
	"github.com/Chu16537/module_master/mtime"
	"github.com/pkg/errors"
)

func main() {
	conf := mlog.Config{
		LogFilePath: "./logs/",
		LogFileName: "admin",
		// ElasticURL:   "http://localhost:9200", // Elasticsearch 的 URL
		// ElasticIndex: "logs_index",            // Elasticsearch 中的索引名称
	}

	t := mtime.GetZero()
	err := mlog.New(&conf, t)
	if err != nil {
		fmt.Println(err)
		return
	}

	opt := &mlog.LogData{
		Tracer: fmt.Sprintf("%d", time.Now().Unix()),
	}

	l := mlog.Get("測試")

	// opt.FnName = "d"
	l.Debug(opt)

	opt.FnName = "i"
	l.Info(opt)

	opt.Err = errorcode.New(errorcode.Data_Unmarshal_Error, errors.Wrap(fmt.Errorf("aa"), "asdadada"))
	opt.FnName = "w"
	l.Warn(opt)

	opt.FnName = "e"
	l.Error(opt)

	mlog.Done()
}
