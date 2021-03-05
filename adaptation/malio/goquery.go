package malio

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

func CardDescExtraction(content, column string) int {
	dom, _ := goquery.NewDocumentFromReader(bytes.NewBufferString(content))
	var reValue string
	dom.Find("div .card-wrap").Each(func(i int, ele *goquery.Selection) {
		if strings.Contains(ele.Text(), column) {
			ele.Find("div .card-body").Find("span").Each(func(_ int, s *goquery.Selection) {
				if reValue != "" {
					return
				}
				reValue = s.Text()
			})
		}
	})

	value, _ := strconv.Atoi(reValue)
	return value
}

func RemainTime(content string) int {
	return CardDescExtraction(content, "会员时长")
}

func OnlineDevice(content string) int {
	return CardDescExtraction(content, "在线设备数")
}
