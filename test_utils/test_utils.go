package test_utils

import (
	"server/utils"
	"sort"
)

// CreateErrorMessage はテスト用のエラーメッセージ構造体を生成する関数
func CreateErrorMessage(field string, message string) map[string]interface{} {
	return map[string]interface{}{
		"field":   field,
		"message": message,
	}
}

func SortErrorMessages(sortData []utils.ErrorMessages) {
	sort.SliceStable(
		sortData, func(i, j int) bool {
			return sortData[i].Field < sortData[j].Field
		},
	)
}
