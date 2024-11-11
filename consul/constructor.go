package consul

import (
	"fmt"
	"net/http"

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
	Id string
	Ip string
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
func RegisterService(config *Config) error {
	// 初始化 Consul 客戶端
	if c == nil {
		client, err := api.NewClient(api.DefaultConfig())
		if err != nil {
			return err
		}

		c = client
	}

	healthURL := fmt.Sprintf("http://%v:%v/consulhealth", config.Addr, config.Port)

	// 註冊服務到 Consul
	registration := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v.%v", config.Name, config.NodeId), // 每台服務應該有唯一 ID
		Name:    config.Name,                                      // 服務名稱（相同名稱表示是同類型服務）
		Address: config.Addr,                                      // 服務的地址
		Port:    config.Port,                                      // 服務的端口
		Tags:    config.Tags,                                      // 可選的服務標籤
		Check: &api.AgentServiceCheck{ // 健康檢查
			HTTP:     healthURL,
			Interval: "10s", // 健康檢查間隔
			Timeout:  "10s",
		},
	}

	// 把 B 服務註冊到 Consul
	err := c.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}

	go health(fmt.Sprintf("%v:%v", config.Addr, config.Port))

	return nil
}

// 啟動一個簡單的 HTTP 服務
func health(ip string) error {
	http.HandleFunc("/consulhealth", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return http.ListenAndServe(ip, nil)
}

// getHealthyServiceInstances 查詢指定服務名稱的健康實例
func GetServer(serviceName string) ([]*ServerInfo, error) {
	fmt.Println("c == nil", c == nil)
	if c == nil {
		client, err := api.NewClient(api.DefaultConfig())
		if err != nil {
			return nil, err
		}

		c = client
	}

	// "passing" 過濾條件確保返回健康的實例
	instances, _, err := c.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	ips := make([]*ServerInfo, len(instances))
	for i, v := range instances {
		ips[i] = &ServerInfo{
			Id: v.Service.ID,
			Ip: fmt.Sprintf("%v:%v", v.Service.Address, v.Service.Port),
		}
	}

	return ips, err
}
