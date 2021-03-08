package main

import (
	"airflow/adaptation/malio"
	"airflow/net"
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
	//go prome.StartPromeServ()

	handlerMutex := &sync.Mutex{}
	hOption := &net.HTTPOptions{}
	ctx, cancel := context.WithCancel(context.Background())
	sigComplete := make(chan struct{})

	switch *mode {
	case "malio":
		go func(ctx context.Context) {
			for {
				// this is not a timeout detection, no need to care about ticker memory leaks
				handlerMutex.Lock()
				hOption = net.New()
				hOption.ProxyURL = *proxy

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
				handlerMutex.Unlock()

				log.Println("polling again...")
				for {
					select {
					case <-time.NewTicker(5 * time.Second).C:
						break
					case <-ctx.Done():
						return
					}
				}
			}
		}(ctx)

		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.NewTicker(4 * time.Second).C:
					handlerMutex.Lock()

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
