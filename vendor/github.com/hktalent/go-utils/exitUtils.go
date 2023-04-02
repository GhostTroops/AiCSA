package go_utils

import (
	"github.com/projectdiscovery/gologger"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Exit struct {
	funcs []func() error
	Lock  *sync.Mutex
}

func (r *Exit) RegClose(fns ...func() error) {
	r.Lock.Lock()
	defer r.Lock.Unlock()
	r.funcs = append(r.funcs, fns...)
}

func NewExit() *Exit {
	x := &Exit{Lock: &sync.Mutex{}}
	// close handler
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			gologger.DefaultLogger.Info().Msg("- Ctrl+C pressed in Terminal")
			for _, fn1 := range x.funcs {
				if nil != fn1 {
					fn1()
				}
			}
			os.Exit(0)
		}()
	}()
	return x
}
