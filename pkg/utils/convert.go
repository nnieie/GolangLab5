package utils

import "strconv"

func I64ToStr(i64 int64) string {
	return strconv.FormatInt(i64, 10)
}
