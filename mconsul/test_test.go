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
		NodeId: uid.CreateID(),
		Name:   "test_name",
		Addr:   "192.168.50.80",
		Port:   8083,
		Tags:   []string{"a", "b"},
	}
	fmt.Println(c.NodeId)
	mconsul.New()

	err := mconsul.RegisterService(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	a, err := mconsul.GetServer(c.Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range a {
		fmt.Println(v)
	}

	time.Sleep(100 * time.Minute)
}
