package main

import (
	"github.com/jacexh/polaris/console/model"
	"github.com/jinzhu/gorm"
)

func main() {
	agent := &model.Agent{
		ID:    "abcd",
		Alias: "demo-project",
	}

	db, err := gorm.Open("mysql", "root:123456@/polaris?charset=utf8")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.Create(&agent)
}
