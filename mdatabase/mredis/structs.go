package mredis

type ZsetRangeOpt struct {
	Min    string
	Max    string
	Offset int64
	Count  int64
}
