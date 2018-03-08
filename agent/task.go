package agent

import (
	"net/http"
	"time"

	"sync"

	"sync/atomic"

	"github.com/jacexh/polaris/log"
)

type (
	SniffTask struct {
		S        *Sniffer
		G        *Gather
		D        time.Duration
		finished int32
	}

	Gather struct {
		ch     chan *http.Request
		handle []RequestHandle
		wg     sync.WaitGroup
	}
)

const (
	taskUnfinished int32 = iota
	taskFinished
)

func NewGather(size int, h ...RequestHandle) *Gather {
	return &Gather{
		ch:     make(chan *http.Request, size),
		handle: h,
	}
}

func (g *Gather) Close() {
	close(g.ch)
	g.wg.Wait()
}

func (g *Gather) Handle() {
	for req := range g.ch {
		g.wg.Add(1)
		go func(r *http.Request) {
			defer g.wg.Done()
			for _, h := range g.handle {
				h(r)
			}
		}(req)
	}
}

func (g *Gather) Accept() chan<- *http.Request {
	return g.ch
}

func NewSniffTask(ip string, port int, d time.Duration, h ...RequestHandle) (*SniffTask, error) {
	gather := NewGather(gatherBuffer, h...)
	s, err := NewSniffer(ip, port, gather.Accept())
	if err != nil {
		return nil, err
	}
	t := &SniffTask{
		S: s,
		G: gather,
		D: d,
	}
	return t, nil
}

func (st *SniffTask) Sniff() error {
	log.Logger.Info("sniffing")
	go st.G.Handle()
	go st.timing()
	return st.S.Run()
}

func (st *SniffTask) Stop() {
	if atomic.CompareAndSwapInt32(&st.finished, taskUnfinished, taskFinished) {
		st.S.Close()
		st.G.Close()
	}
}

func (st *SniffTask) timing() {
	if st.finished == taskFinished || st.D == 0 {
		return
	}

	start := time.Now()
	for atomic.LoadInt32(&st.finished) == taskUnfinished {
		if time.Since(start) > st.D {
			st.Stop()
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
}
