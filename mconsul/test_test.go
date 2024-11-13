package mconsul_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Chu16537/module_master/mconsul"
	"github.com/Chu16537/module_master/muid"
)

func TestXxx(t *testing.T) {

	uid := muid.New(1)

	c := &mconsul.Config{
		ConsulAddr: "192.168.50.80:8500",
		Scheme:     "http",
		NodeId:     uid.CreateID(),
		Name:       "test_name",
		Addr:       "192.168.50.80",
		Port:       8081,
		Tags:       []string{"a", "b"},
	}

	fmt.Println(c.NodeId)

	mc, err := mconsul.New(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	eC := mc.RegisterService()
	if eC != nil {
		fmt.Println(eC)
		return
	}

	a, eC := mc.GetServer(c.Name)
	if eC != nil {
		fmt.Println(eC)
		return
	}

	for _, v := range a {
		fmt.Println(v)
	}

	time.Sleep(100 * time.Minute)
}
