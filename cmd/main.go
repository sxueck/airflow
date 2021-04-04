package main

import (
	"airflow/adaptation"
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

	handlerMutex := &sync.Mutex{}
	hOption := &net.HTTPOptions{}
	ctx, cancel := context.WithCancel(context.Background())
	sigComplete := make(chan struct{})
	var lockPromeServer = false

	switch *mode {
	case "malio":
		go func(ctx context.Context) {
			for {
				// this is not a timeout detection, no need to care about ticker memory leaks
				handlerMutex.Lock()
				hOption = net.New()
				hOption.ProxyURL = *proxy

				for i := 0; i <= 3; i++ {
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
					case <-time.NewTicker(5 * time.Hour).C:
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
				case <-time.NewTicker(40 * time.Second).C:
					handlerMutex.Lock()

					if userinfo, err := malio.ObtainUserInfo(hOption, *domain); err != nil {
						log.Println(err)
					} else {
						err = adaptation.PassMetrics(userinfo)
						if err != nil {
							log.Println(err)
						}

						if !lockPromeServer {
							lockPromeServer = !lockPromeServer
							go prome.StartPromeServ(ctx, userinfo.Name, userinfo.Level)
						}

					}
					handlerMutex.Unlock()
				}
			}
		}(ctx)

		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case <-time.NewTicker(24 * time.Hour).C:
					handlerMutex.Lock()
					err := malio.CheckIn(hOption, *domain)
					if err != nil {
						log.Println(err)
					}
					handlerMutex.Unlock()
				}
			}
		}(ctx)

	case "realnode":

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
