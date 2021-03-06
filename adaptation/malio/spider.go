package malio

import (
	"airflow/adaptation"
	"airflow/net"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// ["set", "user:email", "xxx@gmail.com"] => xxx@gmail.com
// ["Class", "1"] => 1
func BracketExtraction(content, column string) string {
	reg, _ := regexp.Compile(`\[.*?]`)
	c := reg.FindAllString(content, -1)
	for _, v := range c {
		// [xx] => xx
		v = v[1:]
		v = v[:len(v)-1]

		// "key",<space>"value"
		ck := strings.Split(v, ", ")
		if len(ck) >= 2 {
			if strings.Contains(ck[len(ck)-2], column) {
				key := ck[len(ck)-1][1:]
				return key[:len(key)-1]
			}
		}
	}
	return ""
}

func Login(hOption *net.HTTPOptions, domain string, username, password string) {
	hOption.URL = fmt.Sprintf("%s/auth/login", domain)
	hOption.ContentType = "application/x-www-form-urlencoded; charset=UTF-8"
	hOption.POST(fmt.Sprintf("email=%s&passwd=%s&code=", username, password))

	if err := hOption.Err; err != nil {
		fmt.Println(err)
	} else {
		hOption.ObtainCookie()
	}
}

func ObtainUserInfo(hOption *net.HTTPOptions, domain string) (*adaptation.PersonalInfo, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	if hOption.Err != nil {
		return nil, hOption.Err
	}

	fmt.Printf("%+v\n", hOption)

	hOption.URL = fmt.Sprintf("%s/user/profile", domain)
	resBody := hOption.GET()
	if err := hOption.Err; err != nil {
		return nil, err
	}

	// fmt.Println(resBody)
	userinfo := &adaptation.PersonalInfo{}
	userinfo.Name = BracketExtraction(resBody, "user:nickname")
	userinfo.Level = BracketExtraction(resBody, "Class")
	userinfo.Balance = BracketExtraction(resBody, "Money")

	hOption.URL = fmt.Sprintf("%s/user", domain)
	resBody = hOption.GET()
	userinfo.RemainFlow = BracketExtraction(resBody, "Unused_Traffic")
	userinfo.MaxBandwidth = MaxBandwidth(resBody)
	userinfo.TodayUsed = TodayUsed(resBody)
	userinfo.OnlineDevice = OnlineDevice(resBody)
	userinfo.RemainTime = RemainTime(resBody)

	log.Println(userinfo)
	return userinfo, nil
}

func CheckIn(hOption *net.HTTPOptions, domain string) error {
	hOption.URL = fmt.Sprintf("%s/user/checkin", domain)
	res := hOption.POST("")
	fmt.Println(res)
	return hOption.Err
}
