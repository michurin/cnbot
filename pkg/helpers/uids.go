package helpers

import (
	"strconv"
)

func Itoa(x int64) string {
	return strconv.FormatInt(x, 10)
}

func Atoi(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
