package mmgo_test

// import (
// 	"context"
// 	"fmt"
// 	"sync"
// 	"testing"
// 	"time"

// 	"github.com/Chu16537/module_master/mmgo"
// 	"github.com/Chu16537/module_master/mmgo/hmgo"
// 	"github.com/Chu16537/module_master/myaml"
// 	"github.com/Chu16537/module_master/proto/db"
// )

// var (
// 	h   *mmgo.Handler
// 	ctx context.Context = context.Background()
// )

// func Test(t *testing.T) {

// 	c := &mmgo.Config{}

// 	myaml.Read("config.yaml", c)

// 	mmgoH, err := mmgo.New(context.Background(), c)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	h = mmgoH

// 	err = h.CreateCollection(ctx, db.ColName_Club, [][]string{[]string{"club_id"}, []string{"invitation_code"}})
// 	err = h.CreateCollection(ctx, db.ColName_Club_User_Info, [][]string{[]string{"club_id", "account"}, []string{"club_id", "nick_name"}, []string{"user_id"}, []string{"token"}})
// 	err = h.CreateCollection(ctx, db.ColName_Table, [][]string{[]string{"table_id"}})
// 	err = h.CreateCollection(ctx, db.ColName_GameWallet, [][]string{[]string{"user_id", "table_id"}})
// 	err = h.CreateCollection(ctx, db.ColName_WalletLog, [][]string{[]string{"user_id", "order_id"}})
// 	if !err.IsSuccess() {
// 		fmt.Println("err", err)
// 		t.Error(err)
// 		return
// 	}

// 	st := time.Now()
// 	defer func() {
// 		fmt.Println("run time", time.Since(st))
// 	}()

// 	// incr()
// 	clubuserInfo()

// }

// func incr() {
// 	wg := sync.WaitGroup{}
// 	mu := sync.Mutex{}
// 	a := map[int]struct{}{}

// 	wg.Add(1)
// 	for i := 0; i < 1000; i++ {
// 		go func() {
// 			wg.Add(1)
// 			defer wg.Done()

// 			c, _ := mmgo.Incr(h, ctx, "aaa")

// 			mu.Lock()
// 			defer mu.Unlock()
// 			a[c] = struct{}{}
// 		}()
// 	}

// 	time.Sleep(3 * time.Second)
// 	wg.Done()

// 	wg.Wait()

// 	fmt.Println("len(a)", len(a))
// }

// // create club user info
// func clubuserInfo() {

// 	// cui, err := hmgo.CreateClubUserInfo(h, ctx, "account_1", "password", "nickname_1", 1, "token_1")
// 	// if !err.IsSuccess() {
// 	// 	// fmt.Println("CreateClubUserInfo err", err)
// 	// 	fmt.Printf("CreateClubUserInfo 1 %+v\n", err) // 包含堆疊資訊
// 	// 	return
// 	// }

// 	// fmt.Println("cui.UserID", cui.UserID)

// 	// cui2, err := hmgo.CreateClubUserInfo(h, ctx, "account_2", "password", "nickname_2", 1, "token_2")
// 	// if !err.IsSuccess() {
// 	// 	// fmt.Println("CreateClubUserInfo err", err)
// 	// 	fmt.Printf("CreateClubUserInfo 2 %+v\n", err) // 包含堆疊資訊
// 	// 	return
// 	// }

// 	// fmt.Println("cui2.UserID", cui2.UserID)

// 	// cuiOpt := &db.ClubUserInfoOpt{
// 	// 	UserID: cui.UserID,
// 	// }
// 	// datas, err := hmgo.GetClubUserInfo(h, ctx, cuiOpt, nil)
// 	// if !err.IsSuccess() {
// 	// 	fmt.Printf("%+v\n", err) // 包含堆疊資訊
// 	// 	return
// 	// }

// 	// for _, v := range datas {
// 	// 	fmt.Println(v.UserID, v.Account)
// 	// }

// 	nowBalance, err := hmgo.TransBalanceClub(h, ctx, 1, 2, 1, 1000, fmt.Sprintf("order_test_%v", time.Now().Unix()), time.Now().Unix())
// 	if err != nil {
// 		fmt.Printf("TransBalanceClub %+v\n", err) // 包含堆疊資訊
// 		return
// 	}
// 	fmt.Println("nowBalance", nowBalance)

// }
