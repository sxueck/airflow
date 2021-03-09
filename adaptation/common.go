package adaptation

import (
	"airflow/prome"
	"regexp"
	"strconv"
	"strings"
)

type PersonalInfo struct {
	Name         string
	RemainTime   int
	RemainFlow   string
	MaxBandwidth string
	TodayUsed    string
	OnlineDevice int
	Balance      string
	Level        string
}

func PassMetrics(info *PersonalInfo) error {
	var errChan = make(chan error)

	atof := func(s string) float32 { // unified output => MB
		var gb float32 = 1
		if strings.ToLower(s[len(s)-2:]) == "gb" {
			gb = 1024
		}

		reg, _ := regexp.Compile("^[1-9]\\d*\\.\\d*|0\\.\\d*[1-9]\\d*|[1-9]\\d*|0$")
		s = reg.FindString(s)
		var (
			iv  float64
			err error
		)
		if iv, err = strconv.ParseFloat(s, 32); err != nil {
			panic(err)
		}

		// log.Printf("%s -> %f\n", s, iv)
		return gb * float32(iv)
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				errChan <- err.(error)
			} else {
				errChan <- nil
			}
		}()

		prome.ChanTodayUsed <- atof(info.TodayUsed)
		prome.ChanMaxBandwidth <- atof(info.MaxBandwidth)
		prome.ChanBalance <- atof(info.Balance)
		prome.ChanOnlineDeviceCount <- info.OnlineDevice
		prome.ChanRemainFlow <- atof(info.RemainFlow)
		prome.ChanRemainTime <- info.RemainTime
	}()

	return <-errChan
}
