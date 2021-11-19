package main

import (
	"baseFrame/app"
	"fmt"
	"log"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func main() {
	injector, injectorCleanFunc, err := app.BuildInjector()
	if err != nil || injector == nil {
		fmt.Printf("injector error:" + err.Error())
		return
	}
	db := injector.DB
	injectorCleanFunc()

	m := gormigrate.New(db, &gormigrate.Options{
		TableName:                 "migrations",
		IDColumnName:              "id",
		IDColumnSize:              100,
		UseTransaction:            false,
		ValidateUnknownMigrations: false,
	}, []*gormigrate.Migration{
		addDemo,
		addUser,
	},
	)
	if err = m.Migrate(); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}
	log.Printf("Migration did run successfully")
}

var addDemo = &gormigrate.Migration{
	ID: "202011041719", // ID用当前时间 年月日时分
	Migrate: func(db *gorm.DB) error {
		type Demo struct {
			Model
			Describe string `json:"describe"`
		}

		return db.AutoMigrate(new(Demo))
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("demos")
	},
}

var addUser = &gormigrate.Migration{
	ID: "202105151716", // ID用当前时间 年月日时分
	Migrate: func(db *gorm.DB) error {
		type User struct {
			Model
			Mobile   string `gorm:"type:char(11);index;NOT NULL;comment:手机号;" json:"mobile"`
			Nickname string `gorm:"type:varchar(50);NOT NULL;default:'';comment:昵称;" json:"nickname"`
			Password string `json:"-" gorm:"type:varchar(64);NOT NULL;default:'';comment:密码"`
		}

		return db.AutoMigrate(new(User))
	},
	Rollback: func(tx *gorm.DB) error {
		//return tx.Migrator().DropTable("users")
		return nil
	},
}
