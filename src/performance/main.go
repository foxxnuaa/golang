package main

import (
	"fmt"
	"log"
	"performance/generator"
	"performance/lib"
	"performance/test"
	"runtime"
	"time"
)

var printDetail = false

func main() {
	// 设置P最大数量
	runtime.GOMAXPROCS(runtime.NumCPU())

	serverAddr := "127.0.0.1:9797"

	// 初始化调用器
	comm := test.NewTcpComm(serverAddr)

	// 初始化载荷发生器
	resultCh := make(chan *lib.CallResult, 50)
	timeoutNs := 3 * time.Millisecond
	lps := uint32(200)
	durationNs := 12 * time.Second
	log.Printf("Initialize load generator (timeoutNs=%v, lps=%d, durationNs=%v)...", timeoutNs, lps, durationNs)
	gen, err := generator.NewGenerator(
		comm,
		timeoutNs,
		lps,
		durationNs,
		resultCh)
	if err != nil {
		log.Printf("Load generator initialization failing: %s.\n", err)
	}

	// 开始！
	log.Println("Start load generator...")
	gen.Start()

	// 显示结果
	countMap := make(map[lib.ResultCode]int)
	for r := range resultCh {
		countMap[r.Code] = countMap[r.Code] + 1
		if printDetail {
			log.Printf("Result: Id=%d, Code=%d, Msg=%s, Elapse=%v.\n",
				r.Id, r.Code, r.Msg, r.Elapse)
		}
	}

	var total int

	fmt.Println("Code Count:")
	for k, v := range countMap {
		codePlain := lib.GetResultCodePlain(k)
		log.Printf("  Code plain: %s (%d), Count: %d.\n",
			codePlain, k, v)
		total += v
	}

	log.Printf("Total load: %d.\n", total)
	successCount := countMap[lib.RESULT_CODE_SUCCESS]
	tps := float64(successCount) / float64(durationNs/1e9)
	log.Printf("Loads per second: %d; Treatments per second: %f.\n", lps, tps)
}
