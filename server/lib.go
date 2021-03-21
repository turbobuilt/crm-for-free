package main

import (
	"strconv"
	"time"
)

func GetUnixMillisecondsString() string {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(now, 10)
}
func GetLoginExpirationMillisecondsString() string {
	t := time.Now().Add(time.Hour*24).UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(t, 10)
}
