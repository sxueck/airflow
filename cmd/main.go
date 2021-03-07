package main

import (
	"airflow/adaptation/malio"
	"airflow/net"
	"airflow/prome"
	"flag"
	"fmt"
	"log"
)

var (
	domain   = flag.String("domain", "", "")
	username = flag.String("username", "", "")
	password = flag.String("password", "", "")
	mode     = flag.String("mode", "malio", "")
	proxy    = flag.String("proxy", "", "")
)

func main() {
	flag.Parse()
	go prome.StartPromeServ()
	switch *mode {
	case "malio":
		hOption := net.New()
		hOption.ProxyURL = *proxy
		malio.Login(hOption, *domain, *username, *password)
		if userinfo,err := malio.ObtainUserInfo(hOption,*domain);err != nil {
			log.Println(err)
		} else {
			fmt.Printf("%+v\n",userinfo)
		}
	default:
		fmt.Println("please enter the correct matching pattern")
	}

}
