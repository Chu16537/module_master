package errorcode

import (
	"github.com/pkg/errors"
)

var NewRedisError = errors.New("redis new fail")

type Error struct {
	code int
	err  error
}

func New(code int, err error) *Error {
	return &Error{
		code: code,
		err:  err,
	}
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Err() error {
	return e.err
}

func Success() *Error {
	return &Error{
		code: SuccessCode,
		err:  nil,
	}
}
func Server(err error) *Error {
	return &Error{
		code: Server_Error,
		err:  err,
	}
}

func DataUnmarshalError(s string) *Error {
	return &Error{
		code: Data_Unmarshal_Error,
		err:  errors.Errorf("data unmarshal err:%s", s),
	}
}

func DataError(data interface{}) *Error {
	return &Error{
		code: Data_Error,
		err:  errors.Errorf("data :%v", data),
	}
}

func DataNotExist(s string) *Error {
	return &Error{
		code: Data_Not_Exist,
		err:  errors.Errorf("data not exist err:%s", s),
	}
}

func NewIsExist(s string) *Error {
	return &Error{
		code: Data_Is_Exist,
		err:  errors.Errorf("data is exist err:%s", s),
	}
}

func NotClubPermissions(userID uint64, clubID uint64) *Error {
	return &Error{
		code: Club_User_Not_Permissions,
		err:  errors.Errorf("userID:%v clubID:%v not permissions", userID, clubID),
	}
}

func NotInClub(userID uint64, clubID uint64) *Error {
	return &Error{
		code: Club_User_Not_In_Club,
		err:  errors.Errorf("userID:%v not in club:%v", userID, clubID),
	}
}

func ClubUserBalanceLess(userID uint64, nowBalance uint64, amount uint64) *Error {
	return &Error{
		code: Club_User_Balacne_Less,
		err:  errors.Errorf("userID:%v nowBalance:%v < amount:%v", userID, nowBalance, amount),
	}
}

func TableNotExist(clubID, tableID uint64) *Error {
	return &Error{
		code: Table_Not_Exist,
		err:  errors.Errorf("clubID:%v tableID:%v not exist", clubID, tableID),
	}
}

func GameWalletNotExist(userID, tableID uint64) *Error {
	return &Error{
		code: GameWallet_Not_Exist,
		err:  errors.Errorf("gameWallet user:%v table:%v not exist", userID, tableID),
	}
}

func GameWalletNotExists(userIDs []uint64, tableID uint64) *Error {
	return &Error{
		code: GameWallet_Not_Exist,
		err:  errors.Errorf("gameWallet users:%v table:%v not exist", userIDs, tableID),
	}
}

func TokenNotExist(token string) *Error {
	return &Error{
		code: Token_Not_Exist,
		err:  errors.Errorf("token:%v not exist", token),
	}
}

func GameWalleBalanceLess(userID, tableID, nowBalance, amount uint64) *Error {
	return &Error{
		code: GameWallet_Balacne_Less,
		err:  errors.Errorf("gameWallet user:%v table:%v nowBalance:%v < amount:%v", userID, tableID, nowBalance, amount),
	}
}
