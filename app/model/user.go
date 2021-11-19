package model

import (
	"baseFrame/pkg/encryption"
	"baseFrame/pkg/logger"
	"baseFrame/pkg/util"

	"gorm.io/gorm"
)

type User struct {
	Model
	Mobile   string `gorm:"type:char(11);index;NOT NULL;comment:手机号;" json:"mobile"`
	Nickname string `gorm:"type:varchar(50);NOT NULL;default:'';comment:昵称;" json:"nickname"`
	Password string `json:"-" gorm:"type:varchar(64);NOT NULL;default:'';comment:密码"`
}

func (s *User) EncryptionPassword() {
	if s.Password == "" && len(s.Mobile) >= 4 {
		s.Password = "pass" + s.Mobile[len(s.Mobile)-4:len(s.Mobile)]
	}
	s.Password = encryption.GeneratePassword(s.Password)
}

func (s *User) CheckPassword(pwd string) bool {
	return encryption.ComparePassword(s.Password, pwd)
}

func GetUserByMobile(db *gorm.DB, mobile string) (*User, error) {
	user := new(User)
	err := db.Where(User{Mobile: mobile}).First(user).Error
	if !logger.Check(err) {
		return &User{}, err
	}
	return user, nil
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	if user.Nickname == "" {
		user.Nickname = "用户" + util.RandDigitStr(8)
	}
	return nil
}

func (user *User) DoAfterUserCreate(db *gorm.DB) error {
	logger.Debug("user create success")
	if user.Nickname == "" {
		user.Nickname = "用户" + util.RandStr(8)
	}
	return nil
}

func (user *User) DoAfterUserLogin(db *gorm.DB) error {
	logger.Debug("user login success")
	return nil
}
