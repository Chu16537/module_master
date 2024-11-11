package consul_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Chu16537/module_master/consul"
	"github.com/Chu16537/module_master/muid"
)

func TestXxx(t *testing.T) {

	uid := muid.New(1)

	c := &consul.Config{
		NodeId: uid.CreateID(),
		Name:   "test_name",
		Addr:   "192.168.50.80",
		Port:   8083,
		Tags:   []string{"a", "b"},
	}
	fmt.Println(c.NodeId)
	consul.New()

	err := consul.RegisterService(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	a, err := consul.GetServer(c.Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range a {
		fmt.Println(v.Id, v.Ip)
	}

	time.Sleep(100 * time.Minute)
}
