package service

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Base ...
type Base struct {
	Ctx context.Context
}

// SignalNotify ...
func (t *Base) SignalNotify() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		switch <-c {
		case os.Interrupt:
			log.Println("os.Signal", os.Interrupt)
			cancel()
		case syscall.SIGTERM:
			log.Println("os.Signal", syscall.SIGTERM)
			cancel()
		}
	}()

	return ctx
}
