package mconsul

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/hashicorp/consul/api"
)

const (
	consulHealthURL = "/consulhealth"
)

type Config struct {
	ConsulAddr  string   // Consul 地址
	Scheme      string   // http 或 https
	NodeID      int64    // 服务节点 ID
	Name        string   // 服务名称
	Addr        string   // 服务地址
	Port        int      // 服务端口
	Tags        []string // 服务标签
	HealthCheck struct { // 健康检查配置
		Interval string
		Timeout  string
	}
}

type Handler struct {
	config *Config
	c      *api.Client
}

type GetServerInfo struct {
	Addr  string
	Value []byte
}

var kv *api.KVPair

func New(config *Config) (*Handler, error) {
	if config == nil || config.ConsulAddr == "" || config.Name == "" || config.Addr == "" || config.Port <= 0 {
		return nil, fmt.Errorf("invalid config: %+v", config)
	}

	if config.Scheme == "" {
		config.Scheme = "http" // 默认协议为 HTTP
	}
	if config.HealthCheck.Interval == "" {
		config.HealthCheck.Interval = "30s" // 默认健康检查间隔
	}
	if config.HealthCheck.Timeout == "" {
		config.HealthCheck.Timeout = "10s" // 默认健康检查超时
	}

	client, err := api.NewClient(&api.Config{
		Address: config.ConsulAddr,
		Scheme:  config.Scheme,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	h := &Handler{
		config: config,
		c:      client,
	}

	kv = &api.KVPair{
		Key: fmt.Sprintf("%v", h.config.NodeID),
	}

	return h, nil
}

// 服务注册
func (h *Handler) RegisterService() *errorcode.Error {
	ip := fmt.Sprintf("%v:%v", h.config.Addr, h.config.Port)
	healthURL := fmt.Sprintf("%v://%v%v", h.config.Scheme, ip, consulHealthURL)

	registration := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v", h.config.NodeID), // 唯一 ID
		Name:    h.config.Name,                      // 服务名称
		Address: h.config.Addr,                      // 服务地址
		Port:    h.config.Port,                      // 服务端口
		Tags:    h.config.Tags,                      // 服务标签
		Check: &api.AgentServiceCheck{ // 健康检查
			HTTP:     healthURL,
			Interval: h.config.HealthCheck.Interval,
			Timeout:  h.config.HealthCheck.Timeout,
		},
	}

	// 注册服务到 Consul
	if err := h.c.Agent().ServiceRegister(registration); err != nil {
		return errorcode.New(errorcode.Code_Server_Error, err)
	}

	// 启动健康检查服务
	if err := h.startHealthCheckServer(ip); err != nil {
		return errorcode.New(errorcode.Code_Server_Error, err)
	}

	return nil
}

// 启动健康检查 HTTP 服务
func (h *Handler) startHealthCheckServer(ip string) error {
	mux := http.NewServeMux()
	mux.HandleFunc(consulHealthURL, func(w http.ResponseWriter, r *http.Request) {
		// 写入 Consul KV
		_, err := h.c.KV().Put(kv, nil)
		if err != nil {
			http.Error(w, "Failed to update KV", http.StatusInternalServerError)
			return
		}

		// 返回 HTTP 响应
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    ip,
		Handler: mux,
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- server.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(5 * time.Second):
		return nil
	}
}

// 更新kv數值
func (h *Handler) UpdateKV(v []byte) {
	kv.Value = v
	fmt.Println("UpdateKV", kv.Value)
}

// 查询指定服务名称的健康实例
func (h *Handler) GetServer(serviceName string) (map[string]GetServerInfo, *errorcode.Error) {
	instances, _, err := h.c.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, errorcode.New(errorcode.Code_Server_Error, err)
	}

	result := make(map[string]GetServerInfo)
	for _, instance := range instances {
		// 获取服务的健康状态信息（从 KV）
		kvPair, _, err := h.c.KV().Get(instance.Service.ID, nil)
		v := []byte{}
		if err == nil && kvPair != nil {
			v = kvPair.Value
		}

		gsInfo := GetServerInfo{
			Addr:  instance.Service.Address,
			Value: v,
		}
		// 格式化返回值，包括服务地址和状态
		result[instance.Service.ID] = gsInfo
	}

	return result, nil
}
