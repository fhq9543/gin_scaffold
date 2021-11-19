package aliCloudClient

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
)

func InitClient(accessKey, secretKey, regionID string) (*Client, error) {
	client, err := sdk.NewClientWithAccessKey(regionID, accessKey, secretKey)
	if err != nil {
		fmt.Errorf(err.Error())
		return nil, err
	}
	return &Client{client, regionID}, nil
}

type Client struct {
	*sdk.Client
	RegionID string
}

//json数据解析
type Message struct {
	Message   string
	RequestId string
	BizId     string
	Code      string
}
