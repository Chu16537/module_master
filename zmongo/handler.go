package zmongo

import "context"

// 創建 讀寫DB
func NewReadWriteDB(ctx context.Context, readConf *Config, writeConf *Config) (readDB *Handler, writeDB *Handler, err error) {

	readDB, err = New(ctx, readConf)
	if err != nil {
		return
	}

	writeDB, err = New(ctx, writeConf)
	if err != nil {
		return
	}

	return
}
