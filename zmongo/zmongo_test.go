package zmongo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Chu16537/gomodule/zmongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var h *zmongo.Handler

func init() {
	fmt.Println("init")

	conf := &zmongo.Config{
		Addr:     "mongodb://localhost:27017",
		Database: "test",
		Username: "",
		Password: "",
	}

	opt := &options.ClientOptions{}
	opt.ApplyURI(conf.Addr)

	ctx := context.Background()

	zgo, err := zmongo.New(ctx, conf, opt)

	if err != nil {
		fmt.Println("init err", err)
		return
	}

	h = zgo

	fmt.Println("init finish")
}

type Account struct {
	PlayerId uint64
	Account  string
	Password string
	Balance  uint64
}

func TestFindOne(t *testing.T) {

	c := "account"
	f := bson.M{"playerid": 1}

	a := &Account{}

	err := h.FindOne(c, a, f)
	if err != nil {
		fmt.Println("TestFindOne", err)
	}

	fmt.Println("a", a)
}
