package zyaml

import "fmt"

type AA struct {
	Addr string
}
type TestStruct struct {
	AA
}

func Test() {
	t := new(TestStruct)

	// 沒有指定路徑跟 main 同一層 / 有的話相對路徑
	path := "zyaml/env.yaml"
	err := Read(path, t)

	if err != nil {
		fmt.Println("err", err)
		return
	}

	fmt.Println("t", t)
}
