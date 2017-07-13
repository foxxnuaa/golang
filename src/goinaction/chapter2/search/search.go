package search

import (
	"log"
	"sync"
)

// 注册用于搜索的匹配器的映射
var matchers = make(map[string]Matcher)

// Run执行搜索逻辑
func Run(searchTerm string) {
	// 获取数据源列表
	feeds, err := RetrieveFeeds()

	if err != nil {
		log.Fatal(err)
	}

	//创建一个无缓冲的通道
	results := make(chan *Result)

	// 构造一个waitGroup，等待处理数据源
	var waitGroup sync.WaitGroup

	//设置等待数据源
	waitGroup.Add(len(feeds))

	// 为每个数据源启动一个goroutine查找结果
	for _, feed := range feeds {
		//获取一个匹配器进行查找
		matcher, exists := matchers[feed.Type]
		if !exists {
			matcher = matchers["default"]
		}

		go func(matcher Matcher, feed *Feed) {
			Match(matcher, feed, searchTerm, results)

			waitGroup.Done()
		}(matcher, feed)
	}

	//等待所有的查询结果
	go func() {
		waitGroup.Wait()

		close(results)
	}()

	Display(results)
}

func Register(feedType string, matcher Matcher) {
	if _, exists := matchers[feedType]; exists {
		log.Fatalln(feedType, "Matcher already registered")
	}

	log.Println("Register", feedType, "matcher")
	matchers[feedType] = matcher
}
