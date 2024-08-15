package test_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/hmgo"
	"github.com/Chu16537/module_master/mmgo"
	"github.com/Chu16537/module_master/muid"
	"github.com/Chu16537/module_master/proto/db"
)

var h *hmgo.Handler

func init() {
	ctx := context.Background()

	conf := &mmgo.Config{
		// Addr: "mongodb://mongo1:27017,mongo2:27017,mongo3:27017/?replicaSet=rs0",
		Addr: "mongodb://127.0.0.1:27017/",
		// Addr:     "mongodb://127.0.0.1:27017/?directConnection=true",
		Database: "test",
		Username: "",
		Password: "",
	}

	readHandler, err := mmgo.New(ctx, conf)
	if err != nil {
		fmt.Println("readHandler", err)
		return
	}

	writeHandler, err := mmgo.New(ctx, conf)
	if err != nil {
		fmt.Println("writeHandler", err)
		return
	}

	handler, err := hmgo.New(ctx, nil, readHandler, writeHandler)
	if err != nil {
		fmt.Println("new", err)
		return
	}

	h = handler

	errr := h.Init()
	if errr != nil {
		fmt.Println("Init", err)
		return
	}

	errr = h.EveryDay(ctx, time.Now().Unix(), 7)
	if errr != nil {
		fmt.Println("EveryDay", err)
		return
	}
}

func showErr(t *testing.T, err *errorcode.Error) {
	t.Errorf("code:%v \n err:%+v", err.Code(), err.Err())
}

func TestClub(t *testing.T) {
	ctx := context.Background()

	id := uint64(time.Now().Unix())

	var (
		account        string = fmt.Sprintf("%v_%v", "test_create1_account", id)
		password       string = "password"
		nickName       string = fmt.Sprintf("%v_%v", "test_create1_nickname", id)
		clubName       string = fmt.Sprintf("%v_%v", "test_create1_clubname", id)
		expireTime     int64  = time.Now().Unix() + 864000
		invitationCode string = fmt.Sprintf("%v_%v", "test_create1_invitationCode", id)
		token          string = fmt.Sprintf("%v_%v", "test_create1_token", id)
	)

	club, cui, err := h.CreateClub(ctx, account, password, nickName, clubName, db.Club_Status_ON, expireTime, invitationCode, token)
	if err != nil {
		showErr(t, err)
		return
	}

	fmt.Println("club", club.ID, cui.UserID)

	filter := &db.ClubOpt{
		ClubID: club.ID,
	}
	updata := map[string]interface{}{
		"status": 3,
	}

	_, err = h.UpdateClub(ctx, filter, updata)
	if err != nil {
		showErr(t, err)
		return
	}

	c, err := h.GetClub(ctx, filter, nil)
	if err != nil {
		showErr(t, err)
		return
	}

	if c[0].ID != club.ID {
		fmt.Println("俱樂部錯誤")
		return
	}

	var (
		join1Account  string = fmt.Sprintf("%v_%v", "test_join1_account", id)
		join1Password string = "password"
		join1NickName string = fmt.Sprintf("%v_%v", "test_join1_nickname", id)
		join1Token    string = fmt.Sprintf("%v_%v", "test_join1_token", id)
	)

	cui1, err := h.JoinClub(ctx, join1Account, join1Password, join1NickName, club.ID, join1Token)
	if err != nil {
		showErr(t, err)
		return
	}

	// 更新權限 副會長
	cuiOpt := &db.ClubUserInfoOpt{
		UserID: cui1.UserID,
	}
	cuiUpdata := map[string]interface{}{
		"permissions": db.Club_User_Permissions_Vice_President,
	}
	_, err = h.UpdateClubUserInfo(ctx, cuiOpt, cuiUpdata)
	if err != nil {
		showErr(t, err)
		return
	}

	fmt.Println("cui1", cui1.UserID)

	var (
		join2Account  string = fmt.Sprintf("%v_%v", "test_join2_account", id)
		join2Password string = "password"
		join2NickName string = fmt.Sprintf("%v_%v", "test_join2_nickname", id)
		join2Token    string = fmt.Sprintf("%v_%v", "test_join2_token", id)
	)

	cui2, err := h.JoinClub(ctx, join2Account, join2Password, join2NickName, club.ID, join2Token)
	if err != nil {
		showErr(t, err)
		return
	}

	fmt.Println("cui2", cui2.UserID)

	var (
		join3Account  string = fmt.Sprintf("%v_%v", "test_join3_account", id)
		join3Password string = "password"
		join3NickName string = fmt.Sprintf("%v_%v", "test_join3_nickname", id)
		join3Token    string = fmt.Sprintf("%v_%v", "test_join3_token", id)
	)

	cui3, err := h.JoinClub(ctx, join3Account, join3Password, join3NickName, club.ID, join3Token)
	if err != nil {
		showErr(t, err)
		return
	}

	fmt.Println("cui3", cui3.UserID)

	// // 相同俱樂部 會長 > 副會長
	// err = h.TransBalanceClub(ctx, cui.UserID, cui1.UserID, c[0].ID, 1000, fmt.Sprintf("%v_%v_%v", "order", cui.UserID, cui1.UserID))
	// if err != nil {
	// 	showErr(t, err)
	// 	return
	// }

	// // 相同俱樂部 副會長 > 會員
	// err = h.TransBalanceClub(ctx, cui1.UserID, cui2.UserID, c[0].ID, 1000, fmt.Sprintf("%v_%v_%v", "order", cui1.UserID, cui2.UserID))
	// if err != nil {
	// 	showErr(t, err)
	// 	return
	// }

	// // 相同俱樂部  會員 > 會員
	// err = h.TransBalanceClub(ctx, cui2.UserID, cui3.UserID, c[0].ID, 1000, fmt.Sprintf("%v_%v_%v", "order", cui2.UserID, cui3.UserID))
	// if err != nil {
	// 	showErr(t, err)
	// 	return
	// }

	// var (
	// 	account2        string = fmt.Sprintf("%v_%v", "test_create2_account", id)
	// 	password2       string = "password"
	// 	nickName2       string = fmt.Sprintf("%v_%v", "test_create2_nickname", id)
	// 	clubName2       string = fmt.Sprintf("%v_%v", "test_create2_clubname", id)
	// 	expireTime2     int64  = time.Now().Unix() + 864000
	// 	invitationCode2 string = fmt.Sprintf("%v_%v", "test_create2_invitationCode", id)
	// 	token2          string = fmt.Sprintf("%v_%v", "test_create2_token", id)
	// )

	// club2, cui22, err := h.CreateClub(ctx, account2, password2, nickName2, clubName2, db.Club_Status_ON, expireTime2, invitationCode2, token2)
	// if err != nil {
	// 	showErr(t, err)
	// 	return
	// }

	// fmt.Println("club2", club2.ID, cui22.UserID)

	// // 不相同俱樂部 會長 > 會長
	// err = h.TransBalanceClub(ctx, cui.UserID, cui22.UserID, c[0].ID, 1000, fmt.Sprintf("%v_%v_%v", "order", cui.UserID, cui22.UserID))
	// if err != nil {
	// 	showErr(t, err)
	// 	return
	// }

	fmt.Println("TestCreateClub success")
}

func TestGameWallet(t *testing.T) {
	ctx := context.Background()

	id := uint64(time.Now().Unix())

	var (
		account        string = fmt.Sprintf("%v_%v", "test_create1_account", id)
		password       string = "password"
		nickName       string = fmt.Sprintf("%v_%v", "test_create1_nickname", id)
		clubName       string = fmt.Sprintf("%v_%v", "test_create1_clubname", id)
		expireTime     int64  = time.Now().Unix() + 864000
		invitationCode string = fmt.Sprintf("%v_%v", "test_create1_invitationCode", id)
		token          string = fmt.Sprintf("%v_%v", "test_create1_token", id)
	)

	club, cui, err := h.CreateClub(ctx, account, password, nickName, clubName, db.Club_Status_ON, expireTime, invitationCode, token)
	if err != nil {
		showErr(t, err)
		return
	}

	// fmt.Println("club", club.ID)

	// fmt.Println("user", cui.UserID)

	var (
		join1Account  string = fmt.Sprintf("%v_%v", "test_join1_account", id)
		join1Password string = "password"
		join1NickName string = fmt.Sprintf("%v_%v", "test_join1_nickname", id)
		join1Token    string = fmt.Sprintf("%v_%v", "test_join1_token", id)
	)

	cui1, err := h.JoinClub(ctx, join1Account, join1Password, join1NickName, club.ID, join1Token)
	if err != nil {
		showErr(t, err)
		return
	}

	fmt.Println("user", cui1.UserID)

	_, err = h.TransBalanceClub(ctx, cui.UserID, cui1.UserID, club.ID, 1000, fmt.Sprintf("%v_%v_%v", "order", cui.UserID, cui1.UserID), time.Now().Unix())
	if err != nil {
		showErr(t, err)
		return
	}

	table := &db.Table{
		ID:     id,
		ClubID: club.ID,
	}
	err = h.CreateTable(ctx, table)
	if err != nil {
		showErr(t, err)
		return
	}

	gw, err := h.CreateGameWallet(ctx, cui1.UserID, table.ID, 500, time.Now().Unix(), fmt.Sprintf("%v_%v", "test_CreateGameWallet_order", cui1.UserID), time.Now().Unix())
	if err != nil {
		showErr(t, err)
		return
	}

	fmt.Println("gw balance", gw.Balance)

	c, g, err := h.UpdateGameWalletBalance(ctx, cui1.UserID, table.ID, 500, fmt.Sprintf("%v_%v", "test_UpdateGameWalletBalance_order", cui1.UserID), time.Now().Unix())
	if err != nil {
		showErr(t, err)
		return
	}

	fmt.Println("俱樂部金額", c, "錢包金額", g)

	gwOpt := &db.GameWalletOpt{
		// UserIDs: []uint64{cui1.UserID},
		DelExpireTime: time.Now().Unix() + 1000,
	}

	err = h.DeleteGameWallet(ctx, gwOpt, fmt.Sprintf("%v_%v", "test_DeleteGameWallet_order", cui1.UserID), time.Now().Unix())
	if err != nil {
		showErr(t, err)
		return
	}

	fmt.Println("TestGameWallet success")
}

func TestGetGameRecord(t *testing.T) {
	ctx := context.Background()

	id := uint64(time.Now().Unix())

	var (
		account        string = fmt.Sprintf("%v_%v", "test_create1_account", id)
		password       string = "password"
		nickName       string = fmt.Sprintf("%v_%v", "test_create1_nickname", id)
		clubName       string = fmt.Sprintf("%v_%v", "test_create1_clubname", id)
		expireTime     int64  = time.Now().Unix() + 864000
		invitationCode string = fmt.Sprintf("%v_%v", "test_create1_invitationCode", id)
		token          string = fmt.Sprintf("%v_%v", "test_create1_token", id)
	)

	club, cui, err := h.CreateClub(ctx, account, password, nickName, clubName, db.Club_Status_ON, expireTime, invitationCode, token)
	if err != nil {
		showErr(t, err)
		return
	}

	var (
		join1Account  string = fmt.Sprintf("%v_%v", "test_join1_account", id)
		join1Password string = "password"
		join1NickName string = fmt.Sprintf("%v_%v", "test_join1_nickname", id)
		join1Token    string = fmt.Sprintf("%v_%v", "test_join1_token", id)
	)

	cui1, err := h.JoinClub(ctx, join1Account, join1Password, join1NickName, club.ID, join1Token)
	if err != nil {
		showErr(t, err)
		return
	}

	_, err = h.TransBalanceClub(ctx, cui.UserID, cui1.UserID, club.ID, 10000, fmt.Sprintf("%v_%v_%v", "order", cui.UserID, cui1.UserID), time.Now().Unix())
	if err != nil {
		showErr(t, err)
		return
	}

	table := &db.Table{
		ID:     id,
		ClubID: club.ID,
	}
	err = h.CreateTable(ctx, table)
	if err != nil {
		showErr(t, err)
		return
	}

	_, err = h.CreateGameWallet(ctx, cui.UserID, table.ID, 10000, time.Now().Unix(), fmt.Sprintf("%v_%v", "test_CreateGameWallet_order", cui1.UserID), time.Now().Unix())
	if err != nil {
		showErr(t, err)
		return
	}

	_, err = h.CreateGameWallet(ctx, cui1.UserID, table.ID, 10000, time.Now().Unix(), fmt.Sprintf("%v_%v", "test_CreateGameWallet_order", cui1.UserID), time.Now().Unix())
	if err != nil {
		showErr(t, err)
		return
	}

	userIds := []uint64{cui.UserID, cui1.UserID}

	count := 7
	opt := &db.GameRecordOpt{
		StartTimeUnix: time.Now().UTC().Unix(),
		EndTimeUnix:   time.Now().UTC().Unix() + int64(86400*count),
	}

	for i := 0; i < count; i++ {
		gr := &db.GameRecord{
			GameRecordID: muid.GetUID(),
			CreateTime:   time.Now().UTC().Unix() + int64(86400*i),
			ClubID:       uint64(i),
			TableID:      table.ID,
			GameID:       i + 1,
			GameType:     i + 1,
		}

		urs := make([]*db.UserRecord, len(userIds))
		infoMap := map[uint64]*db.UserRecordInfodMultiple{}

		for j, v := range userIds {
			infoMap[v] = &db.UserRecordInfodMultiple{
				BetZone:           []uint64{v},
				TotalBet:          2 * v,
				WinZone:           []uint64{3 * v},
				TotalWin:          4 * v,
				EffectiveBet:      []float64{5.5 * float64(v)},
				TotalEffectiveBet: 6.6 * float64(v),
			}

			infoByte, err := json.Marshal(infoMap[v])
			if err != nil {
				fmt.Println("Marshal", err)
				return
			}

			urs[j] = &db.UserRecord{
				UserID:        v,
				GameRecordID:  gr.GameRecordID,
				CreateTime:    gr.CreateTime,
				ClubID:        gr.ClubID,
				TableID:       gr.TableID,
				GameID:        gr.GameID,
				GameType:      gr.GameType,
				ResultBalance: int64(infoMap[v].TotalWin - infoMap[v].TotalBet),
				GameResult:    gr.GameResult,
				Info:          infoByte,
			}

		}

		h.CreateGameRecord(ctx, gr, urs, infoMap)
	}

	findOpt := &db.FindOpt{
		Limit: 1000,
	}

	grs, total, err := h.GetGameRecord(ctx, opt, findOpt)
	if err != nil {
		showErr(t, err)
		return
	}

	fmt.Println("total", total, len(grs))

	for _, v := range grs {
		fmt.Println("TestGetGameRecord", v.GameRecordID)
	}

	urOpt := &db.UserRecordOpt{
		UserID:        userIds[0],
		StartTimeUnix: time.Now().Unix(),
		EndTimeUnix:   time.Now().Unix() + int64(86400*count),
	}

	urFindOpt := &db.FindOpt{
		Start: 2,
		Limit: 5,
	}

	urs, total2, err := h.GetUserRecord(ctx, urOpt, urFindOpt)
	if err != nil {
		showErr(t, err)
		return
	}

	fmt.Println("total2", total2)

	for _, v := range urs {
		fmt.Println(v.GameRecordID)
	}

	res, err := h.GetUserRecordTotalResultBalance(ctx, urOpt)
	if err != nil {
		showErr(t, err)
		return
	}

	fmt.Println(res)

	fmt.Println("TestGetGameRecord success")
}

func TestCreateUserRecord(t *testing.T) {
	ctx := context.Background()

	count := 4
	opt := &db.UserRecordOpt{
		UserID:        1,
		StartTimeUnix: time.Now().Unix(),
		EndTimeUnix:   time.Now().Unix() + int64(86400*count),
	}

	for i := 1; i < count; i++ {
		for j := 1; j < count; j++ {
			ur := &db.UserRecord{
				UserID:        1,
				GameRecordID:  fmt.Sprintf("%v_%v", i, j),
				CreateTime:    time.Now().Unix() + int64(86400*(i-1)),
				ClubID:        1,
				TableID:       uint64(j),
				GameID:        1,
				GameType:      1,
				ResultBalance: int64(i*100 + j),
				GameResult:    nil,
				Info:          nil,
			}

			h.CreateUserRecord(ctx, []*db.UserRecord{ur})
		}
	}

	findOpt := &db.FindOpt{
		Start: 2,
		Limit: 5,
	}
	urs, total, err := h.GetUserRecord(ctx, opt, findOpt)
	if err != nil {
		showErr(t, err)
		return
	}

	fmt.Println("total", total)
	for _, v := range urs {
		fmt.Println(v.GameRecordID)
	}

	res, err := h.GetUserRecordTotalResultBalance(ctx, opt)
	if err != nil {
		showErr(t, err)
		return
	}

	fmt.Println(res)

	fmt.Println("TestCreateUserRecord success")
}
