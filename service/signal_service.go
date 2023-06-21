package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var SignalInterruptHandler *InterruptHandler

func InitInterruptHandler() {
	SignalInterruptHandler = NewInterruptHandle()
}

func InterruptHandleAddFunc(f func()) {
	if SignalInterruptHandler != nil {
		SignalInterruptHandler.FuncList = append(SignalInterruptHandler.FuncList, f)
	}
}

var SignalTermHandler *TermHandler

func InitTermHandler() {
	SignalTermHandler = NewTermHandle()
}

func TermHandleAddFunc(f func()) {
	if SignalTermHandler != nil {
		SignalTermHandler.FuncList = append(SignalTermHandler.FuncList, f)
	}
}

// ---------------------------------------------------
// SIGINT信号监听，func列表执行后从容退出
func NewInterruptHandle() *InterruptHandler {
	h := &InterruptHandler{
		make([]func(), 0),
	}
	go h.handle()
	return h
}

type InterruptHandler struct {
	FuncList []func()
}

func (h *InterruptHandler) registerFunc(f func()) {
	h.FuncList = append(h.FuncList, f)
}

func (h *InterruptHandler) handle() {
	ch := make(chan os.Signal, 0)
	signal.Notify(ch, syscall.SIGINT)
	if s, ok := <-ch; ok {
		fmt.Println("捕捉到" + s.String() + "信号")
		for _, f := range h.FuncList {
			f()
		}
	}
	signal.Stop(ch)
	close(ch)
	os.Exit(0)
}

// ---------------------------------------------------
// SIGTERM信号监听，func列表执行后从容退出
func NewTermHandle() *TermHandler {
	h := &TermHandler{
		make([]func(), 0),
	}
	go h.handle()
	return h
}

type TermHandler struct {
	FuncList []func()
}

func (h *TermHandler) registerFunc(f func()) {
	h.FuncList = append(h.FuncList, f)
}

func (h *TermHandler) handle() {
	ch := make(chan os.Signal, 0)
	signal.Notify(ch, syscall.SIGTERM)
	if s, ok := <-ch; ok {
		fmt.Println("捕捉到" + s.String() + "信号")
		for _, f := range h.FuncList {
			f()
		}
	}
	signal.Stop(ch)
	close(ch)
	os.Exit(0)
}
