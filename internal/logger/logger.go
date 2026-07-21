package logger

import (
	"fmt"
	"time"
)

func timestamp() string {
	return time.Now().Format("15:04:05")
}

func Info(msg string) {
	fmt.Printf("[%s] INFO  %s\n", timestamp(), msg)
}

func Success(msg string) {
	fmt.Printf("[%s] DONE  %s\n", timestamp(), msg)
}

func Warning(msg string) {
	fmt.Printf("[%s] WARN  %s\n", timestamp(), msg)
}

func Error(msg string) {
	fmt.Printf("[%s] ERROR %s\n", timestamp(), msg)
}
