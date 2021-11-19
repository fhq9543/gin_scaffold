package db

import (
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestDBSet(t *testing.T) {
	_, err := gorm.Open(mysql.Open(""), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	return
}
