package agent

import (
	"errors"
	"net"
	"os"
	"sync"
	"time"

	"github.com/satori/go.uuid"
)

type (
	// Agent agent服务顶级对象
	Agent struct {
		ID    string              // 全局唯一id，用于标识出不同的agent
		Alias string              // 别名，可以自定义，或者取系统hostname
		ips   map[string][]string // {"eth0": ["192.168.10.1", "172.16.85.1"]}
		tasks *taskManager
	}

	taskManager struct {
		interval time.Duration
		tasks    sync.Map
	}
)

// NewAgent 实例化一个Agent对象
func NewAgent(alias string) *Agent {
	if alias == "" {
		var err error
		alias, err = os.Hostname()
		if err != nil {
			alias = "unknown"
		}
	}
	return &Agent{ID: uuid.NewV4().String(), Alias: alias, ips: getIPAddr()}
}

func getIPAddr() map[string][]string {
	ret := map[string][]string{}
	ips, err := net.Interfaces()
	if err != nil {
		return ret
	}
	for _, ifa := range ips {
		address, err := ifa.Addrs()
		if err != nil {
			continue
		}
		ret[ifa.Name] = []string{}
		for _, addr := range address {
			ret[ifa.Name] = append(ret[ifa.Name], addr.String())
		}
	}
	return ret
}

// newTask 根据ip、端口实例化一个SniffTask对象
func (tm *taskManager) newTask(ip string, port int, duration time.Duration) error {
	st, err := NewSniffTask(ip, port, duration)
	if err != nil {
		return err
	}
	tm.tasks.Store(st.ID, st)
	return nil
}

// run 执行的指定的SniffTask
func (tm *taskManager) run(id string) error {
	var err error
	var founded bool
	tm.tasks.Range(func(k, v interface{}) bool {
		if k.(string) == id {
			founded = true
			err = v.(*SniffTask).Sniff()
			return false
		}
		return true
	})
	if !founded {
		return errors.New("cannot found SniffTask with id " + id)
	}
	return err
}

// stop 手动停止SniffTask
func (tm *taskManager) stop(id string) error {
	var founded bool
	tm.tasks.Range(func(key, value interface{}) bool {
		if key.(string) == id {
			founded = true
			value.(*SniffTask).Stop()
			return false
		}
		return true
	})
	tm.tasks.Delete(id)
	if !founded {
		return errors.New("cannot found SniffTask with id " + id)
	}
	return nil
}

// autoRemove 自动清理已经停止的任务
func (tm *taskManager) autoRemove() {
	var ret []string
	tm.tasks.Range(func(key, value interface{}) bool {
		if value.(*SniffTask).Status() == StatusFinished {
			ret = append(ret, key.(string))
		}
		return true
	})
	for _, t := range ret {
		tm.tasks.Delete(t)
	}
}
