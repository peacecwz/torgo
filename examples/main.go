package main

import (
	"context"
	"net/http"
	"time"

	"github.com/peacecwz/torgo"
)

func main() {
	options := &torgo.Options{
		Debug: true,
		GeneralOptions: &torgo.GeneralOptions{
			SocksPort: 52507,
			Logging: []torgo.LogConfig{
				{
					SeverityRange: "notice",
					Destinations:  []string{"stdout"},
				},
			},
		},
	}
	torPrx, err := torgo.NewTorProxy(options)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	err = torPrx.Start(ctx)
	if err != nil {
		panic(err)
	}
	defer torPrx.Close()

	httpClient := http.Client{
		Transport: &http.Transport{
			DialContext: torPrx.GetProxy().DialContext,
		},
	}

	resp, err := httpClient.Get("https://check.torproject.org/api/ip")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("not ok")
	}

	bodyString := ""
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			bodyString += string(buf[:n])
		}
		if err != nil {
			break
		}
	}

	println(bodyString)
}
