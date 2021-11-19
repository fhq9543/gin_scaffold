package sms

import (
	"baseFrame/pkg/aliCloudClient"
	"baseFrame/pkg/logger"
	"testing"
)

func TestSMS(t *testing.T) {
	regionId := ""
	accessKey := ""
	secretKey := ""
	client, err := aliCloudClient.InitClient(accessKey, secretKey, regionId)
	if !logger.Check(err) {
		t.Failed()
	}
	ss := &SMS{client}
	err = ss.SendSms("", "", "", map[string]interface{}{"code": 123456})
	if !logger.Check(err) {
		t.Failed()
	}
}
