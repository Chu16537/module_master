package mgracefulshutdown

// shutdown時由低到高 依序關閉
const (
	GF_Level_Controller = 1
	GF_Level_Server     = 2
	GF_Level_MQ         = 3
	GF_Level_DB         = 4 // db 跟 redis

	GF_Level_Log = 99 // log 最後
)
