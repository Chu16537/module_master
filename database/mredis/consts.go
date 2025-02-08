package mredis

const (
	Mode_Singleton = 0 // 單台
	Mode_Cluster   = 1 // 群集

	Key_Node         = "node"           // 節點編號(incr)
	NodeSecond       = 600              // node 時間
	NodeUpdateSecond = NodeSecond - 100 // node 更新時間
)
