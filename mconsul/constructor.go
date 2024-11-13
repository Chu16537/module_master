package mconsul

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/hashicorp/consul/api"
)

type Config struct {
	NodeId int64
	Name   string
	Addr   string
	Port   int
	Tags   []string
}

type ServerInfo struct {
	NodeId string
	Ip     string
}

var c *api.Client

func New() error {
	// 初始化 Consul 客戶端
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}

	c = client
	return nil
}

// 服務註冊
func RegisterService(config *Config) *errorcode.Error {
	// 初始化 Consul 客戶端
	if c == nil {
		client, err := api.NewClient(api.DefaultConfig())
		if err != nil {
			return errorcode.Server(err)
		}

		c = client
	}

	healthURL := fmt.Sprintf("http://%v:%v/consulhealth", config.Addr, config.Port)

	// 註冊服務到 Consul
	registration := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v", config.NodeId), // 每台服務應該有唯一 ID
		Name:    config.Name,                      // 服務名稱（相同名稱表示是同類型服務）
		Address: config.Addr,                      // 服務的地址
		Port:    config.Port,                      // 服務的端口
		Tags:    config.Tags,                      // 可選的服務標籤
		Check: &api.AgentServiceCheck{ // 健康檢查
			HTTP:     healthURL,
			Interval: "30s", // 健康檢查間隔
			Timeout:  "10s",
		},
	}

	// 把服務註冊到 Consul
	err := c.Agent().ServiceRegister(registration)
	if err != nil {
		return errorcode.Server(err)
	}

	err = health(fmt.Sprintf("%v:%v", config.Addr, config.Port))
	if err != nil {
		return errorcode.Server(err)
	}

	return nil
}

// 啟動一個簡單的 HTTP 服務
func health(ip string) error {
	http.HandleFunc("/consulhealth", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	errChan := make(chan error, 1)

	go func() {
		err := http.ListenAndServe(ip, nil)

		if err != nil {
			errChan <- err
		}
	}()

	// 等待n秒判斷是否有錯
	select {
	case err := <-errChan:
		return err

	case <-time.After(5 * time.Second):
		// 等待5秒發現沒有錯誤
		return nil
	}
}

// 查詢指定服務名稱的健康實例
// 回傳 map[nodeId]Addr
func GetServer(serviceName string) (map[string]string, *errorcode.Error) {
	if c == nil {
		client, err := api.NewClient(api.DefaultConfig())
		if err != nil {
			return nil, errorcode.Server(err)
		}

		c = client
	}

	//  過濾條件確保返回健康的實例
	instances, _, err := c.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, errorcode.Server(err)
	}

	result := map[string]string{}
	for _, v := range instances {
		result[v.Service.ID] = v.Service.Address
	}

	return result, nil
}
