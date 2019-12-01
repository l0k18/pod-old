package app

import (
	"github.com/p9c/pod/cmd/kopach"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/cli"
	"os"

	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/util/interrupt"
)

func kopachHandle(cx *conte.Xt) func(c *cli.Context) (err error) {
	return func(c *cli.Context) (err error) {
		log.INFO("starting up kopach standalone miner for parallelcoin")
		Configure(cx, c)
		quit := make(chan struct{})
		interrupt.AddHandler(func() {
			close(quit)
			os.Exit(0)
		})
		kopach.Main(cx, quit)
		log.DEBUG("kopach main finished")
		//<-quit
		return
	}
}
