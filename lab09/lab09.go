package main

import (
	"flag"
	"fmt"

	"github.com/gocolly/colly"
)

type Push struct {
	Tag      string
	UserId   string
	Content  string
	DateTime string
}

func main() {

	maxComments := flag.Int("max", 10, "Max number of comments to show")
	flag.Parse()
	count := 1

	c := colly.NewCollector()
	c.OnHTML("div.push", func(e *colly.HTMLElement) {
		if count > *maxComments {
			return
		}

		push := Push{}

		push.Tag = e.ChildText("span.push-tag")
		push.UserId = e.ChildText("span.push-userid")
		push.DateTime = e.ChildText("span.push-ipdatetime")
		push.Content = e.ChildText("span.push-content")

		fmt.Printf("%d. 名字：%s，留言%s，時間： %s\n", count, push.UserId, push.Content, push.DateTime)

		count++
	})

	c.Visit("https://www.ptt.cc/bbs/joke/M.1481217639.A.4DF.html")

	c.Wait()
}
