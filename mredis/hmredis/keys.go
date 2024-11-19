package hmredis

import "fmt"

func keyTable(tableID uint64) string {
	return fmt.Sprintf("table_%v", tableID)
}
