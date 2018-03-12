package agent

import (
	"fmt"
	"time"

	"sync"

	"github.com/jacexh/polaris/log"
	"github.com/satori/go.uuid"
)

type (
	// SniffTask 监听任务对象，包含监听、采集（处理）部分
	SniffTask struct {
		ID     string // 全局唯一id
		s      *sniffer
		g      *gather
		D      time.Duration // 任务持续时间
		Since  time.Time     // 嗅探开始时间
		status TaskStatus    // 任务状态
		mu     sync.RWMutex
	}

	TaskStatus int32
)

const (
	// StatusNew 任务尚未开始
	StatusNew TaskStatus = iota
	// StatusSniffing 嗅探中，表示任务进行中
	StatusSniffing
	// StatusFinished 任务已完成
	StatusFinished
)

// NewSniffTask SniffTask工厂函数
func NewSniffTask(ip string, port int, d time.Duration, h ...RequestHandle) (*SniffTask, error) {
	gather := newGather(gatherBuffer, h...)
	s, err := newSniffer(ip, port, gather.recv())
	if err != nil {
		return nil, err
	}
	t := &SniffTask{
		ID: uuid.NewV4().String(),
		s:  s,
		g:  gather,
		D:  d,
	}
	return t, nil
}

// Sniff SniffTask任务开始
func (st *SniffTask) Sniff() error {
	st.mu.Lock()
	st.status = StatusSniffing
	st.mu.Unlock()

	go st.g.handle()
	go st.timing()
	return st.s.run()
}

// Stop 停止嗅探任务，包括停止嗅探、采集处理
func (st *SniffTask) Stop() {
	if st.Status() == StatusSniffing {
		st.set(StatusFinished)

		st.s.close()
		st.g.close()
		log.Logger.Info(fmt.Sprintf("task finisehd: %s", st.s.description()))
	}
}

func (st *SniffTask) timing() {
	if st.Status() != StatusSniffing || st.D == 0 {
		return
	}

	st.Since = time.Now()
	for st.Status() == StatusSniffing {
		if st.Duration() >= st.D {
			st.Stop()
			return
		}
		time.Sleep(200 * time.Millisecond)
	}
	log.Logger.Info(fmt.Sprintf("task finished: %s", st.s.description()))
}

// WithHandle 绑定RequestHandle，用来处理采集到*http.Request对象
func (st *SniffTask) WithHandle(h ...RequestHandle) {
	st.g.with(h...)
}

// Status 返回该SniffTask的状态，分别为未开始、进行中、已完成
func (st *SniffTask) Status() TaskStatus {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.status
}

func (st *SniffTask) set(s TaskStatus) {
	st.mu.Lock()
	defer st.mu.Unlock()
	st.status = s
}

func (st *SniffTask) Duration() time.Duration {
	if st.Status() == StatusSniffing {
		return time.Since(st.Since)
	}
	return 0
}
