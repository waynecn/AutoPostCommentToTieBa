package main

import (
	"bufio"
	"log"
	"os"
)

var RandomContent []string

func InitRandomContent() {
	RandomContent = []string{}
	f, err := os.OpenFile("tangshi300.txt", os.O_RDONLY, 0666)
	if err != nil {
		log.Println("打开文件失败:", err)
		return
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || len(line) <= 6 {
			continue
		}

		if line[0] == '0' || line[0] == '1' || line[0] == '2' || line[0] == '3' {
			continue
		}

		RandomContent = append(RandomContent, line)
	}
}
