package utils

import "fmt"

// 函數接受一個整數切片和要從切片中刪除的整數的可變數量引數
func RemoveItems(slice []string, itemsToRemove ...string) []string {
	var s []string
	s = append(s, slice...)

	// 創建一個 map 來存儲要刪除的整數
	itemsToRemoveMap := make(map[string]bool)
	for _, item := range itemsToRemove {
		itemsToRemoveMap[item] = true
	}
	// 使用兩個索引
	i, j := 0, 0

	// 遍歷切片，將不等於要刪除的整數的元素添加到原切片中
	for i < len(s) {
		if !itemsToRemoveMap[s[i]] {
			fmt.Println(i, j)
			s[j] = s[i]
			j++
		}
		i++
	}
	return s[:j]
}
