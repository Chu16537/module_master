package mlog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/sirupsen/logrus"
)

// ElasticsearchHook 发送日志到 Elasticsearch 的 hook
type ElasticsearchHook struct {
	client *elasticsearch.Client
	index  string
}

// NewElasticsearchHook 创建一个新的 Elasticsearch hook
func newElasticsearchHook(elasticURL, index string) (*ElasticsearchHook, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{elasticURL},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ElasticsearchHook{
		client: client,
		index:  index,
	}, nil
}

// Fire 实现 logrus.Hook 接口
func (hook *ElasticsearchHook) Fire(entry *logrus.Entry) error {
	// 格式化日志条目为 JSON
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	if err := encoder.Encode(entry.Data); err != nil {
		return err
	}

	// 发送到 Elasticsearch
	req := esapi.IndexRequest{
		Index:      hook.index,
		DocumentID: entry.Data["topic"].(string) + "_" + entry.Time.Format(time.RFC3339),
		Body:       &buf,
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), hook.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

// Levels 实现 logrus.Hook 接口
func (hook *ElasticsearchHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
