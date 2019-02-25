package models

import (
	"fmt"
	"go_base/db"
)

func Migrate() {
	err := db.Xorm.Sync2(new(AggregatorInfo), new(Application), new(Qualification))
	if err != nil {
		panic(fmt.Sprintf("mysql migrate error. %s", err.Error()))
	}
}
