package main

import (
	"context"
	"time"

	"github.com/peacecwz/torgo"
)

func main() {
	options := &torgo.Options{
		GeneralOptions: &torgo.GeneralOptions{
			SocksPort: 52507,
			Logging: []torgo.LogConfig{
				{
					SeverityRange: "notice",
					Destinations:  []string{"stderr"},
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

}
