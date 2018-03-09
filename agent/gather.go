package agent

import (
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

type gather struct {
	ch      chan *http.Request
	handles []RequestHandle
	wg      sync.WaitGroup
}

// newGather 实例化一个Gather对象
func newGather(size int, h ...RequestHandle) *gather {
	return &gather{
		ch:      make(chan *http.Request, size),
		handles: h,
	}
}

func (g *gather) close() {
	close(g.ch)
	g.wg.Wait()
}

// Handle Gather主函数，开始处理采集到的*http.Request对象
func (g *gather) handle() {
	for req := range g.ch {
		g.wg.Add(1)
		go func(r *http.Request) {
			defer g.wg.Done()
			for _, h := range g.handles {
				h(r)
			}
			// todo: 是否应该关闭Request.Body
			io.Copy(ioutil.Discard, r.Body)
			r.Body.Close()
		}(req)
	}
}

// recv 返回一个chan，用于接收*http.Request
func (g *gather) recv() chan<- *http.Request {
	return g.ch
}

func (g *gather) with(h ...RequestHandle) {
	g.handles = append(g.handles, h...)
}
