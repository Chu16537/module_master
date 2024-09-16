package db

type FindOpt struct {
	Start uint64
	Limit uint64
}

func checkTimeUxin(start, end int64) (int64, int64) {
	// 時間 start 比 end 小
	if start > end {
		return end, start
	}
	return start, end
}

func (o *FindOpt) ToMgo() {
	if o.Start > 0 {
		o.Start--
	}

	if o.Limit == 0 {
		o.Limit = 1
	}

}

type UpdateBalanceInfo struct {
	UserID  uint64
	ClubID  uint64
	TableID uint64
	Balance uint64
}

type TotalCount struct {
	Count int64 `bson:"count"`
}
