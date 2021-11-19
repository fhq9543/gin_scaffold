package db

import (
	"baseFrame/pkg/config"
	"baseFrame/pkg/logger"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func DBSet(cfg *config.Config) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	dsn := cfg.GetConfig("db", "dsn")
	if dsn == "" {
		fmt.Println("empty dsn")
		return nil, err
	}

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.New(),
	})
	if err != nil {
		_ = fmt.Errorf(err.Error())
		return nil, err
	}

	return db.Debug(), nil
}

type Tx struct {
	*gorm.DB
	commit bool
}

func Begin(db *gorm.DB) *Tx {
	return &Tx{
		DB: db.Begin(),
	}
}

func (tx *Tx) RollbackIfFailed() {
	if tx.commit {
		return
	}
	if err := tx.Rollback().Error; err != nil {
		log.Println("rollback failed", err)
	}
}

func (tx *Tx) RollbackIfFailedRecover() {
	if err := recover(); err != nil {
		log.Println("捕获异常:", err)
		tx.Rollback()
		return
	}
	if tx.commit {
		return
	}
	if err := tx.Rollback().Error; err != nil {
		log.Println("rollback failed", err)
	}
}

func (tx *Tx) Commit() {
	if err := tx.DB.Commit().Error; err != nil {
		log.Println("commit failed", err)
		return
	}
	tx.commit = true
}
