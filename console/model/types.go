package model

import (
	"fmt"

	"github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type (
	Model struct {
		ID      string `gorm:"column:id;type:varchar(36)"`
		Deleted bool   `gorm:"column:deleted;type:tinyint(1);default:0"`
		CTime   int64  `gorm:"column:ctime;type:bigint(13)"`
		MTime   int64  `gorm:"column:mtime;type:bigint(13)"`
	}

	Agent struct {
		*Model
		Alias       string   `gorm:"column:alias"`
		IP          []string `gorm:"-"`
		RAWColumnIP []byte   `gorm:"column:ip;type:blob"` // 请无视该字段，仅仅用于读取/设置原始的ip字段
		Offline     bool     `gorm:"column:offline;type:tinyint(1);default:1"`
		Enabled     bool     `gorm:"column:enabled;type:tinyint(1);default:1"`
	}

	Service struct {
		*Model
		Name      string `gorm:"column:name;size:45"`
		Remark    string `gorm:"column:remark;size:128"`
		Instances []*Instance
	}

	// Instance foreignkey正确用法：在创建该记录时，会同时创建agent/service记录
	Instance struct {
		*Model
		SN         string   `gorm:"column:sn"`
		ListenIP   string   `gorm:"column:listen_ip"`
		ListenPort int      `gorm:"column:listen_port"`
		Agent      *Agent   `gorm:"foreignkey:AgentID"`
		AgentID    string   `gorm:"column:agent_id"`
		Service    *Service `gorm:"foreignkey:ServiceID"`
		ServiceID  string   `gorm:"column:service_id"`
	}
)

func (a *Agent) TableName() string {
	return "agent"
}

func (a *Agent) BeforeCreate() error {
	data, err := json.Marshal(a.IP)
	if err != nil {
		return err
	}
	a.RAWColumnIP = data
	fmt.Println(string(data))
	return nil
}

func (a *Agent) AfterFind() error {
	return json.Unmarshal(a.RAWColumnIP, &a.IP)
}

func (a *Agent) BeforeUpdate() error {
	data, err := json.Marshal(a.IP)
	if err != nil {
		return err
	}
	a.RAWColumnIP = data
	return nil
}

func (s *Service) TableName() string {
	return "service"
}

func (i *Instance) TableName() string {
	return "instance"
}
