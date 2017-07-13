package main

import (
	"fmt"
	_ "goinaction/chapter2/matchers"
	"goinaction/chapter2/search"
	"log"
	"os"
)

func init() {
	//将日志输出到标准输出
	log.SetOutput(os.Stdout)
}

func main() {
	fmt.Println("Hello,GoInAction")

	search.Run("president")
}
