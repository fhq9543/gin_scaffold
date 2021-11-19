package ctlUser

import (
	"baseFrame/app/db"
	"baseFrame/app/model"
	"baseFrame/pkg/auth"
	"baseFrame/pkg/logger"
	"baseFrame/pkg/response"
	"errors"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type UserCtl struct {
	*gorm.DB
	*response.Response
	*auth.Auth
}

func (ctl UserCtl) UserLogin(ctx *gin.Context) {
	bindStu := &struct {
		Mobile   string `json:"mobile" form:"mobile"`
		Password string `json:"password" form:"password"`
	}{}
	if err := ctx.ShouldBind(bindStu); !logger.Check(err) {
		response.Error(ctx, 40001, err.Error(), gin.H{})
		return
	}

	tx := db.Begin(ctl.DB)
	defer tx.RollbackIfFailed()

	// 已注册就直接登录
	user, err := model.GetUserByMobile(tx.DB, bindStu.Mobile)
	if !logger.Check(err) {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(ctx, 20005, "注册失败", nil)
			return
		}
		// 未注册
		user.Mobile = bindStu.Mobile
		user.Password = bindStu.Password
		user.EncryptionPassword()
		if !logger.Check(tx.Create(user).Error) {
			response.Error(ctx, 20006, "创建钩子执行失败", nil)
			return
		}
		if !logger.Check(user.DoAfterUserCreate(tx.DB)) {
			response.Error(ctx, 20006, "创建钩子执行失败", nil)
			return
		}
	} else {
		// 校验密码
		if !user.CheckPassword(bindStu.Password) {
			response.Error(ctx, 20010, "密码错误", nil)
			return
		}
	}

	if !logger.Check(user.DoAfterUserLogin(tx.DB)) {
		response.Error(ctx, 20007, "登录钩子执行失败", nil)
		return
	}

	jwtToken, err := ctl.Auth.AppLogin(user.ID)
	if !logger.Check(err) {
		response.Error(ctx, 20008, "生成jwt失败", nil)
		return
	}

	tx.Commit()
	response.Success(ctx, "ok", map[string]interface{}{
		"token": jwtToken,
		"user":  user,
	})
}

func (ctl UserCtl) UserList(ctx *gin.Context) {
	response.ReturnList(ctx, ctl.WithContext(ctx), &[]*model.User{})
}

func (ctl UserCtl) UserDetail(ctx *gin.Context) {
	instance := &model.User{}

	err := ctl.WithContext(ctx).Where("id = ?", ctx.Param("id")).
		Preload(clause.Associations).First(instance).Error
	if !logger.Check(err) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(ctx, 20001, "记录不存在", nil)
			return
		}
		response.Error(ctx, 20002, "查询失败", nil)
		return
	}

	response.Success(ctx, "", instance)
}

func (ctl UserCtl) UserChangePwd(ctx *gin.Context) {
	bindStu := &struct {
		Password string `json:"password" form:"password"`
	}{}
	if err := ctx.ShouldBind(bindStu); !logger.Check(err) {
		response.Error(ctx, 40001, err.Error(), gin.H{})
		return
	}
	instance := &model.User{}
	err := ctl.WithContext(ctx).Where("id = ?", ctx.Param("id")).First(instance).Error
	if !logger.Check(err) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(ctx, 20001, "记录不存在", nil)
			return
		}
		response.Error(ctx, 20002, "查询失败", nil)
		return
	}
	instance.Password = bindStu.Password
	instance.EncryptionPassword()
	if !logger.Check(ctl.WithContext(ctx).Model(instance).UpdateColumn("password", instance.Password).Error) {
		response.Error(ctx, 20004, "保存失败", nil)
		return
	}
	response.Success(ctx, "", instance)
}

func (ctl UserCtl) UserUpdate(ctx *gin.Context) {
	bindStu := &struct {
		Nickname string `json:"nickname" form:"nickname"`
	}{}
	if err := ctx.ShouldBind(bindStu); !logger.Check(err) {
		response.Error(ctx, 40001, err.Error(), gin.H{})
		return
	}

	instance := &model.User{}
	err := ctl.WithContext(ctx).Where("id = ?", ctx.Param("id")).First(instance).Error
	if !logger.Check(err) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(ctx, 20001, "记录不存在", nil)
			return
		}
		response.Error(ctx, 20002, "查询失败", nil)
		return
	}
	instance.Nickname = bindStu.Nickname
	if !logger.Check(ctl.WithContext(ctx).Save(instance).Error) {
		response.Error(ctx, 20004, "保存失败", nil)
		return
	}
	response.Success(ctx, "", instance)
}

func (ctl UserCtl) UserDelete(ctx *gin.Context) {
	if !logger.Check(ctl.WithContext(ctx).Where("id = ?", ctx.Param("id")).Delete(&model.User{}).Error) {
		response.Error(ctx, 20001, "删除失败", nil)
		return
	}
	response.Success(ctx, "ok", nil)
}
