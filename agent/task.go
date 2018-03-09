package agent

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/jacexh/polaris/log"
)

type (
	// SniffTask 监听任务对象，包含监听、采集（处理）部分
	SniffTask struct {
		s        *sniffer
		g        *gather
		D        time.Duration
		finished int32
	}
)

const (
	taskUnfinished int32 = iota
	taskFinished
)

// NewSniffTask SniffTask工厂函数
func NewSniffTask(ip string, port int, d time.Duration, h ...RequestHandle) (*SniffTask, error) {
	gather := newGather(gatherBuffer, h...)
	s, err := newSniffer(ip, port, gather.recv())
	if err != nil {
		return nil, err
	}
	t := &SniffTask{
		s: s,
		g: gather,
		D: d,
	}
	return t, nil
}

// Sniff SniffTask任务开始
func (st *SniffTask) Sniff() error {
	log.Logger.Info("sniffing")
	go st.g.handle()
	go st.timing()
	return st.s.run()
}

// Stop 停止嗅探任务，包括停止嗅探、采集处理
func (st *SniffTask) Stop() {
	if atomic.CompareAndSwapInt32(&st.finished, taskUnfinished, taskFinished) {
		st.s.close()
		st.g.close()
		log.Logger.Info(fmt.Sprintf("task finisehd: %s", st.s.description()))
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
	log.Logger.Info(fmt.Sprintf("task finisehd: %s", st.s.description()))
}

// WithHandle 绑定RequestHandle，用来处理采集到*http.Request对象
func (st *SniffTask) WithHandle(h ...RequestHandle) {
	st.g.with(h...)
}
