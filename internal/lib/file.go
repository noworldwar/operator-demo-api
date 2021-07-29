package lib

import (
	"fmt"
	"os"
	"time"
)

func WriteLog(name, msg string) {
	filename := fmt.Sprintf("./log/%s%s.log", name, time.Now().Format("060102"))
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Open File Error:", err)
	}
	defer f.Close()
	if _, err := f.WriteString(msg); err != nil {
		fmt.Println("Write File Error:", err)
	}
}
