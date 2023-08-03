package goev

import (
	"syscall"
)

const (
	// EPOLLET Refer to sys/epoll.h
	EPOLLET = 1 << 31

	// EvIn is readable event
	EvIn uint32 = syscall.EPOLLIN | syscall.EPOLLRDHUP

	// EvOut is writeable event
	EvOut uint32 = syscall.EPOLLOUT | syscall.EPOLLRDHUP

	// EvInET is readable event in EPOLLET mode
	EvInET uint32 = EvIn | EPOLLET

	// EvOutET is readable event in EPOLLET mode
	EvOutET uint32 = EvOut | EPOLLET

	// EvEventfd used for eventfd
	EvEventfd uint32 = syscall.EPOLLIN | syscall.EPOLLRDHUP // Not ET mode

	// EvAccept used for acceptor
	// 用水平触发, 循环Accept有可能会导致不可控
	EvAccept uint32 = syscall.EPOLLIN | syscall.EPOLLRDHUP

	// EvConnect used for connector
	EvConnect uint32 = syscall.EPOLLIN | syscall.EPOLLOUT | syscall.EPOLLRDHUP
)

// EvHandler is the event handling interface of the Reactor core
//
// The same EvHandler is repeatedly registered with the Reactor
type EvHandler interface {
	setEvPoll(ep *evPoll)
	getEvPoll() *evPoll

	setReactor(r *Reactor)
	GetReactor() *Reactor

	setTimerItem(ti *timerItem)
	getTimerItem() *timerItem

	// Call by acceptor on `accept` a new fd or connector on `connect` successful
	//
	// Call OnClose() when return false
	OnOpen(fd int) bool

	// EvPoll catch readable i/o event
	//
	// Call OnClose() when return false
	OnRead(fd int) bool

	// EvPoll catch writeable i/o event
	//
	// Call OnClose() when return false
	OnWrite(fd int) bool

	// EvPoll catch connect result
	// Only be asynchronously called after connector.Connect() returns nil
	//
	// Will not call OnClose() after OnConnectFail() (So you don't need to manually release the fd)
	// The param err Refer to ev_handler.go: ErrConnect*
	OnConnectFail(err error)

	// EvPoll catch timeout event
	// The parameter 'millisecond' represents the time of batch retrieval of epoll events, not the current
	// precise time. Use it with caution (as it can reduce the frequency of obtaining the current
	// time to some extent).
	//
	// Remove timer when return false
	OnTimeout(millisecond int64) bool

	// Call by reactor(OnOpen must have been called before calling OnClose.)
	//
	// You need to manually release the fd resource call fd.Close()
	// You'd better only call fd.Close() here.
	OnClose(fd int)
}

// Detecting illegal struct copies using `go vet`
type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
