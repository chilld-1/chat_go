package tools

import (
	"fmt"
	"time"
)

func GetCurrentTime() time.Time {
	return time.Now()
}
func GetCurrentTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
func ErrorMsg(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
func Log(format string, args ...interface{}) {
	fmt.Printf("[%s] "+format+"\n", append([]interface{}{GetCurrentTimeStr()}, args...)...)
}
