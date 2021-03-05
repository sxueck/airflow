package malio

import (
	"fmt"
	"regexp"
)

func DescExtraction(content, column, e string) string {
	reg, _ := regexp.Compile(fmt.Sprintf("%s%s", column, e))
	regCleanClo, _ := regexp.Compile(e)
	return regCleanClo.FindString(reg.FindString(content))
}

func TodayUsed(content string) string {
	e := "\\d.*?B"
	c := "今日已用: "
	return DescExtraction(content, c, e)
}

func MaxBandwidth(content string) string {
	e := "\\d.*?M"
	c := "最高带宽: "
	return DescExtraction(content, c, e)
}


