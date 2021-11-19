package ctlDemo

import (
	"baseFrame/app/model"
	"baseFrame/pkg/logger"
	"baseFrame/pkg/response"
	"errors"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type DemoCtl struct {
	*gorm.DB
	*response.Response
}

func (ctl DemoCtl) DemoList(ctx *gin.Context) {
	response.ReturnList(ctx, ctl.WithContext(ctx), &[]*model.Demo{})
}

func (ctl DemoCtl) DemoDetail(ctx *gin.Context) {
	instance := &model.Demo{}

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

func (ctl DemoCtl) DemoCreate(ctx *gin.Context) {
	instance := &model.Demo{}
	if !logger.Check(ctx.ShouldBind(&instance)) {
		response.Error(ctx, 20001, "获取参数失败", nil)
		return
	}
	if !logger.Check(ctl.WithContext(ctx).Create(instance).Error) {
		response.Error(ctx, 20002, "创建失败", nil)
		return
	}
	response.Success(ctx, "", instance)
}

func (ctl DemoCtl) DemoUpdate(ctx *gin.Context) {
	instance := &model.Demo{}
	err := ctl.WithContext(ctx).Where("id = ?", ctx.Param("id")).First(instance).Error
	if !logger.Check(err) {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(ctx, 20001, "记录不存在", nil)
			return
		}
		response.Error(ctx, 20002, "查询失败", nil)
		return
	}
	model := instance.Model
	if !logger.Check(ctx.ShouldBind(&instance)) {
		response.Error(ctx, 20003, "获取参数失败", nil)
		return
	}
	instance.Model = model
	if !logger.Check(ctl.WithContext(ctx).Save(instance).Error) {
		response.Error(ctx, 20004, "保存失败", nil)
		return
	}
	response.Success(ctx, "", instance)
}

func (ctl DemoCtl) DemoDelete(ctx *gin.Context) {
	if !logger.Check(ctl.WithContext(ctx).Where("id = ?", ctx.Param("id")).Delete(&model.Demo{}).Error) {
		response.Error(ctx, 20001, "删除失败", nil)
		return
	}
	response.Success(ctx, "ok", nil)
}
