package utils

import (
	"fmt"
	"time"
)

func PrintWithTime(print string) {
	t := time.Now()
	fmt.Printf("[%02d:%02d:%02d] %v\n", t.Hour(), t.Minute(), t.Second(), print)
}
