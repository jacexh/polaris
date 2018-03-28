package model

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestRelatedField(t *testing.T) {
	db, err := gorm.Open("mysql", "root:123456@/polaris?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		assert.Error(t, err)
	}
	defer db.Close()

	instance := new(Instance)
	instance.Agent = new(Agent)
	err = db.Debug().Where("id = ?", "6b2477e1-0e1c-40f8-83ff-6750384ab379").First(instance).Error
	if err != nil {
		assert.Error(t, err)
	}

	err = db.Debug().Model(instance).Related(instance.Agent).Error
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, instance.Agent.ID, "6b2477e1-0e1c-40f8-83ff-6750384ab377")
}

func TestCreateRecordsByForeignKey(t *testing.T) {
	db, err := gorm.Open("mysql", "root:123456@/polaris?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		assert.Error(t, err)
	}
	defer db.Close()

	instance := &Instance{
		Model:      &Model{ID: uuid.NewV4().String()},
		SN:         "test-002",
		ListenIP:   "127.0.0.1",
		ListenPort: 8080,
		Agent:      &Agent{Model: &Model{ID: uuid.NewV4().String()}, Alias: "hehe"},
	}
	err = db.Create(instance).Error
	if err != nil {
		assert.Error(t, err)
	}
}

func TestCreateService(t *testing.T) {
	db, err := gorm.Open("mysql", "root:123456@/polaris?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		assert.Error(t, err)
	}
	defer db.Close()

	serv := &Service{Model: &Model{ID: uuid.NewV4().String()}, Name: "upay-gateway", Remark: "核心支付网关"}
	err = db.Debug().Create(serv).Error
	if err != nil {
		assert.Error(t, err)
	}

	g := &Service{Model: &Model{ID: serv.ID}}
	err = db.Debug().Where("id = ?", g.ID).Take(g).Error
	if err != nil {
		assert.Error(t, err)
	}
	assert.True(t, g.CTime > 0)
	assert.True(t, g.MTime > 0)
	assert.False(t, g.Deleted)
}

func TestFindInstancesByService(t *testing.T) {
	db, err := gorm.Open("mysql", "root:123456@/polaris?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		assert.Error(t, err)
	}
	defer db.Close()

	serv := new(Service)
	err = db.Debug().Where("name = ?", "upay-gateway").Take(serv).Error
	if err != nil {
		assert.Error(t, err)
	}

	assert.True(t, serv.Model.ID != "")

	var instances []*Instance
	db.Debug().Model(&serv).Related(&instances)
	assert.True(t, len(instances) > 0)
	for _, instance := range instances {
		assert.Equal(t, instance.ServiceID, serv.Model.ID)
	}
}
