package errorcode

import (
	"github.com/pkg/errors"
)

var NewRedisError = errors.New("redis new fail")

type Error struct {
	Code int   `json:"code"`
	Err  error `json:"err"`
}

func New(code int, err error) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}

func (e *Error) GetCode() int {
	return e.Code
}

func (e *Error) GetErr() error {
	return e.Err
}

func Success() *Error {
	return &Error{
		Code: SuccessCode,
		Err:  nil,
	}
}
func Server(err error) *Error {
	return &Error{
		Code: Server_Error,
		Err:  err,
	}
}

func DataUnmarshalError(s string) *Error {
	return &Error{
		Code: Data_Unmarshal_Error,
		Err:  errors.Errorf("data unmarshal Err:%s", s),
	}
}

func DataMarshalError(s string) *Error {
	return &Error{
		Code: Data_Marshal_Error,
		Err:  errors.Errorf("data marshal Err:%s", s),
	}
}

func DataError(data interface{}) *Error {
	return &Error{
		Code: Data_Error,
		Err:  errors.Errorf("data :%v", data),
	}
}

func DataNotExist(s string) *Error {
	return &Error{
		Code: Data_Not_Exist,
		Err:  errors.Errorf("data not exist Err:%s", s),
	}
}

func NewIsExist(s string) *Error {
	return &Error{
		Code: Data_Is_Exist,
		Err:  errors.Errorf("data is exist Err:%s", s),
	}
}

func NotClubPermissions(userID uint64, clubID uint64) *Error {
	return &Error{
		Code: Club_User_Not_Permissions,
		Err:  errors.Errorf("userID:%v clubID:%v not permissions", userID, clubID),
	}
}

func NotInClub(userID uint64, clubID uint64) *Error {
	return &Error{
		Code: Club_User_Not_In_Club,
		Err:  errors.Errorf("userID:%v not in club:%v", userID, clubID),
	}
}

func ClubUserBalanceLess(userID uint64, nowBalance uint64, amount uint64) *Error {
	return &Error{
		Code: Club_User_Balacne_Less,
		Err:  errors.Errorf("userID:%v nowBalance:%v < amount:%v", userID, nowBalance, amount),
	}
}

func TableNotExist(clubID, tableID uint64) *Error {
	return &Error{
		Code: Table_Not_Exist,
		Err:  errors.Errorf("clubID:%v tableID:%v not exist", clubID, tableID),
	}
}

func GameWalletNotExist(userID, tableID uint64) *Error {
	return &Error{
		Code: GameWallet_Not_Exist,
		Err:  errors.Errorf("gameWallet user:%v table:%v not exist", userID, tableID),
	}
}

func GameWalletNotExists(userIDs []uint64, tableID uint64) *Error {
	return &Error{
		Code: GameWallet_Not_Exist,
		Err:  errors.Errorf("gameWallet users:%v table:%v not exist", userIDs, tableID),
	}
}

func TokenNotExist(token string) *Error {
	return &Error{
		Code: Token_Not_Exist,
		Err:  errors.Errorf("token:%v not exist", token),
	}
}

func GameWalleBalanceLess(userID, tableID, nowBalance, amount uint64) *Error {
	return &Error{
		Code: GameWallet_Balacne_Less,
		Err:  errors.Errorf("gameWallet user:%v table:%v nowBalance:%v < amount:%v", userID, tableID, nowBalance, amount),
	}
}

func MQPublishError(subject string, data []byte, err error) *Error {
	return &Error{
		Code: MQ_Publish_Error,
		Err:  errors.Errorf("subject:%v data:%s err:%v", subject, string(data), err.Error()),
	}
}
