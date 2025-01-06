package mmgo_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Chu16537/module_master/mmgo"
	"github.com/Chu16537/module_master/mmgo/hmgo"
	"github.com/Chu16537/module_master/myaml"
	"github.com/Chu16537/module_master/proto/db"
)

var (
	h *mmgo.Handler
)

func Test(t *testing.T) {

	c := &mmgo.Config{}

	myaml.Read("config.yaml", c)

	mmgoH, err := mmgo.New(context.Background(), c)
	if err != nil {
		t.Error(err)
		return
	}

	h = mmgoH

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = h.CreateCollection(ctx, db.ColName_Club, [][]string{[]string{"club_id"}, []string{"invitation_code"}})
	err = h.CreateCollection(ctx, db.ColName_Club_User_Info, [][]string{[]string{"account"}, []string{"user_id"}, []string{"token"}})
	err = h.CreateCollection(ctx, db.ColName_Table, [][]string{[]string{"table_id"}})
	err = h.CreateCollection(ctx, db.ColName_GameWallet, [][]string{[]string{"user_id", "table_id"}})
	if err.IsNotSuccess() {
		fmt.Println("err", err)
		t.Error(err)
		return
	}

	st := time.Now()
	defer func() {
		fmt.Println("aaa time", time.Since(st))
	}()

	createClub()

	// defer h.Done()
}

// 創建俱樂部
func createClub() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cuiOpt := &db.ClubUserInfoOpt{
		UserID: 1,
	}

	err := hmgo.CreateClubUserInfo(h, ctx, &db.ClubUserInfo{
		Account:  "test",
		Password: "test",
		NickName: "test",
		ClubID:   1,
		Token:    "test",
	})

	if err.IsNotSuccess() {
		fmt.Println("CreateClubUserInfo err", err)
		return
	}

	datas, err := hmgo.GetClubUserInfo(h, ctx, cuiOpt, nil)
	if err.IsNotSuccess() {
		fmt.Println("GetClubUserInfo err", err)
		return
	}

	for i, v := range datas {
		fmt.Println(i, v)
	}
}
