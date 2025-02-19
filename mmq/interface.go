package mmq

type IMQHandler interface {
	// 推送
	Publish(subject string, data []byte) error

	// 訂閱
	Sub(subject string, consumer string, mode SubMode, subChan chan MQSubData) error

	// 取消訂閱
	UnSub(consumer string) error

	// 刪除指定 subject 內的所有訊息
	DelSubject(streamName, subjectName string) error
}
