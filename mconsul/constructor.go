package mconsul

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/hashicorp/consul/api"
)

const (
	consulhealthURL = "/consulhealth"
)

type Config struct {
	ConsulAddr string   // consul
	Scheme     string   // consul
	NodeId     int64    // 服務的
	Name       string   // 服務的
	Addr       string   // 服務的
	Port       int      // 服務的
	Tags       []string // 服務的
}

type ServerInfo struct {
	NodeId string
	Ip     string
}

type Handler struct {
	c      *api.Client
	config *Config
}

func New(config *Config) (*Handler, error) {
	schema := config.Scheme
	if schema == "" {
		schema = "http"
	}

	c := &api.Config{
		Address: config.ConsulAddr, // 指定 Consul 的地址和端口
		Scheme:  schema,            // 如果使用 HTTPS 則改為 "https"
	}

	// 初始化 Consul 客戶端
	client, err := api.NewClient(c)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		c:      client,
		config: config,
	}

	return h, nil
}

// 服務註冊
func (h *Handler) RegisterService() *errorcode.Error {

	ip := fmt.Sprintf("%v:%v", h.config.Addr, h.config.Port)
	healthURL := fmt.Sprintf("http://%v%v", ip, consulhealthURL)

	// 註冊服務到 Consul
	registration := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v", h.config.NodeId), // 每台服務應該有唯一 ID
		Name:    h.config.Name,                      // 服務名稱（相同名稱表示是同類型服務）
		Address: h.config.Addr,                      // 服務的地址
		Port:    h.config.Port,                      // 服務的端口
		Tags:    h.config.Tags,                      // 可選的服務標籤
		Check: &api.AgentServiceCheck{ // 健康檢查
			HTTP:     healthURL,
			Interval: "30s", // 健康檢查間隔
			Timeout:  "10s",
		},
	}

	// 把服務註冊到 Consul
	err := h.c.Agent().ServiceRegister(registration)
	if err != nil {
		return errorcode.Server(err)
	}

	err = health(ip)
	if err != nil {
		return errorcode.Server(err)
	}

	return nil
}

// 啟動一個簡單的 HTTP 服務
func health(ip string) error {
	http.HandleFunc(consulhealthURL, func(w http.ResponseWriter, r *http.Request) {
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
func (h *Handler) GetServer(serviceName string) (map[string]string, *errorcode.Error) {
	//  過濾條件確保返回健康的實例
	instances, _, err := h.c.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, errorcode.Server(err)
	}

	result := map[string]string{}
	for _, v := range instances {
		result[v.Service.ID] = v.Service.Address
	}

	return result, nil
}
