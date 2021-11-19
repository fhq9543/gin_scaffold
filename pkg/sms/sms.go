package sms

import (
	"baseFrame/pkg/aliCloudClient"
	"baseFrame/pkg/config"
	"baseFrame/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

func InitSMS(cfg *config.Config, db *gorm.DB) (*SMS, error) {
	accessKey := cfg.GetConfig("dysms", "accessKey")
	secretKey := cfg.GetConfig("dysms", "secretKey")
	regionID := cfg.GetConfig("dysms", "regionId")
	client, err := aliCloudClient.InitClient(accessKey, secretKey, regionID)
	if !logger.Check(err) {
		return nil, err
	}
	env := cfg.GetConfig("", "env")
	return &SMS{client, db, env}, nil
}

type SMS struct {
	*aliCloudClient.Client
	*gorm.DB
	Env string
}

func (client *SMS) SendSms(phone, signName, templateCode string, templateParam map[string]interface{}) (err error) {
	if client.Env != "prd" {
		// 测试环境走白名单，避免影响正式用户
		count := int64(0)
		client.Table("dev_notice_white_lists").Where("mobile = ?", phone).
			Count(&count)
		if count == 0 {
			return nil
		}
	}
	request := requests.NewCommonRequest()             // 构造一个公共请求
	request.Method = "POST"                            // 设置请求方式
	request.Product = "Ecs"                            // 指定产品
	request.Scheme = "https"                           // https | http
	request.Domain = "dysmsapi.aliyuncs.com"           // 指定域名则不会寻址，如认证方式为 Bearer Token 的服务则需要指定
	request.Version = "2017-05-25"                     // 指定产品版本
	request.ApiName = "SendSms"                        // 指定接口名
	request.QueryParams["RegionId"] = client.RegionID  // 地区
	request.QueryParams["SignName"] = signName         //阿里云验证过的项目名 自己设置
	request.QueryParams["TemplateCode"] = templateCode //阿里云的短信模板号 自己设置
	request.QueryParams["PhoneNumbers"] = phone        //手机号

	templateParamString, err := json.MarshalIndent(templateParam, "", "	")
	if err != nil {
		fmt.Errorf(err.Error())
		return err
	}
	request.QueryParams["TemplateParam"] = string(templateParamString) //短信模板中的参数

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		fmt.Errorf(err.Error())
		return err
	}
	logger.Debug(response.String())
	var message aliCloudClient.Message //阿里云返回的json信息对应的类
	//判断错误信息
	err = json.Unmarshal(response.GetHttpContentBytes(), &message)
	if err != nil {
		fmt.Errorf(err.Error())
		return err
	}
	if message.Message != "OK" {
		fmt.Errorf(message.Message)
		err = errors.New(message.Message)
		logger.Debug(request.QueryParams)
		return
	}
	return nil
}
