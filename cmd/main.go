package main

import (
	"airflow/adaptation/malio"
	"airflow/net"
	"airflow/prome"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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

	handlerMutex := &sync.Mutex{}
	hOption := &net.HTTPOptions{}
	ctx, cancel := context.WithCancel(context.Background())
	sigComplete := make(chan struct{})

	switch *mode {
	case "malio":
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.NewTicker(5 * time.Hour).C: // 这里不是超时检测，不需要关心ticker内存泄漏问题
					handlerMutex.Lock()
					hOption = net.New()
					hOption.ProxyURL = *proxy
					handlerMutex.Unlock()

					log.Println("polling again...")
					time.Sleep(5 * time.Hour)
				}
			}
		}(ctx)

		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.NewTicker(40 * time.Second).C:
					handlerMutex.Lock()

					for i := 0; i >= 3; i++ {
						malio.Login(hOption, *domain, *username, *password)

						// login error
						if hOption.Err == nil {
							break
						}
						if i >= 3 {
							close(sigComplete)
						}
					}

					if userinfo, err := malio.ObtainUserInfo(hOption, *domain); err != nil {
						log.Println(err)
					} else {
						fmt.Printf("%+v\n", userinfo)
					}
					handlerMutex.Unlock()
				}
			}
		}(ctx)
	default:
		fmt.Println("please enter the correct matching pattern")
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-sigterm:
		log.Println("receive stop signal")
	case <-sigComplete:
	}

	cancel()
	time.Sleep(1 * time.Second)
}
