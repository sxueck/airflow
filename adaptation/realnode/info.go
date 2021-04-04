package realnode

import (
	"airflow/adaptation"
	"airflow/net"
	"encoding/json"
	"fmt"
	"log"
)

func Login(hOption *net.HTTPOptions, domain string, username, password string) {
	hOption.URL = fmt.Sprintf("%s/api/v1/passport/auth/login", domain)
	hOption.ContentType = "application/x-www-form-urlencoded"
	hOption.POST(fmt.Sprintf("email=%s&password=%s", username, password))

	if err := hOption.Err; err != nil {
		fmt.Println(err)
	} else {
		hOption.ObtainCookie()
	}
}

func ObtainUserInfo(hOption *net.HTTPOptions, domain string) (*adaptation.PersonalInfo, error) {
	hOption.URL = fmt.Sprintf("%s/api/v1/user/info", domain)
	resJson := hOption.GET()
	if err := hOption.Err; err != nil {
		return nil, err
	}

	ui := &UserInfo{}
	userinfo := &adaptation.PersonalInfo{}
	err := json.Unmarshal([]byte(resJson), ui)
	if err != nil {
		return nil, err
	}

	uiData := ui.Data
	userinfo.Name = uiData.Email
	userinfo.Level = uiData.PlanID
	userinfo.Balance = uiData.Balance
	userinfo.RemainFlow = uiData.RemindExpire
	userinfo.MaxBandwidth = string(fmt.Sprintf("%dM", uiData.TransferEnable/1024/1024/1024))
	log.Println(userinfo)
	return userinfo, nil
}
