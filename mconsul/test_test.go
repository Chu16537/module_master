package mconsul_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/chu16537/module_master/mconsul"
	"github.com/chu16537/module_master/mjson"
	"github.com/chu16537/module_master/muid"
)

func TestXxx(t *testing.T) {

	muid.New(1)

	c := &mconsul.Config{
		ConsulAddr: "192.168.50.80:8500",
		Scheme:     "http",
		NodeID:     muid.CreateID(),
		Name:       "test_name",
		Addr:       "192.168.50.80",
		Port:       8081,
		Tags:       []string{"a", "b"},
	}

	fmt.Println(c.NodeID)

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

	time.Sleep(5 * time.Second)

	a, eC := mc.GetServer(c.Name)
	if eC != nil {
		fmt.Println(eC)
		return
	}

	for i, v := range a {
		fmt.Println(i, v)
	}

	b, err := mjson.Marshal("asdasdasda")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("UpdateKV")
	mc.UpdateKV(b)
	time.Sleep(3 * time.Second)

	a, eC = mc.GetServer(c.Name)
	if eC != nil {
		fmt.Println(eC)
		return
	}

	for i, v := range a {
		fmt.Println(i, v)
	}

	b, err = mjson.Marshal("123213132")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("UpdateKV")
	mc.UpdateKV(b)
	time.Sleep(3 * time.Second)

	a, eC = mc.GetServer(c.Name)
	if eC != nil {
		fmt.Println(eC)
		return
	}

	for i, v := range a {
		fmt.Println(i, v)
	}
}
