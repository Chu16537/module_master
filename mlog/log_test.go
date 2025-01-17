package mlog_test

import (
	"fmt"
	"testing"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mlog"
)

func Test_A(t *testing.T) {
	conf := mlog.Config{
		FilePath: "./logs/",
		Name:     "module",
		// ElasticURL:   "http://localhost:9200", // Elasticsearch 的 URL
		// ElasticIndex: "logs_index",            // Elasticsearch 中的索引名称
	}

	err := mlog.New(&conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	l := mlog.Get("module")

	l.New(mlog.ErrorLevel, "Test_A", "tracer", nil, errorcode.New(errorcode.Code_Server_Error, fmt.Errorf("test err")))
	l.New(mlog.InfoLevel, "Test_A", "tracer2", "test info", errorcode.New(errorcode.Code_Server_Error, fmt.Errorf("test err")))
}
