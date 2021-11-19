package oss

import (
	"baseFrame/pkg/config"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// StsSignVersion sts sign version
	StsSignVersion = "1.0"
	// StsAPIVersion sts api version
	StsAPIVersion = "2015-04-01"
	// StsHost sts host, detail information:https://help.aliyun.com/document_detail/66053.html?spm=a2c4g.11186623.2.12.151338dfkKJJKh
	StsHost = "https://sts.aliyuncs.com/"
	// TimeFormat time fomrat
	TimeFormat = "2006-01-02T15:04:05Z"
	// RespBodyFormat  respone body format
	RespBodyFormat = "JSON"
	// PercentEncode '/'
	PercentEncode = "%2F"
	// HTTPGet http get method
	HTTPGet = "GET"
)

// token app config
// see: https://help.aliyun.com/knowledge_detail/39744.html
// see: https://help.aliyun.com/document_detail/116819.html?spm=a2c4g.11186623.2.24.1db17e44U1DrSX#task-2458383
type Sts struct {
	*Config
	StsTokenInfo struct {
		AccessKeyId     string `json:"access_key_id"`
		AccessKeySecret string `json:"access_key_secret"`
		BucketName      string `json:"bucket_name"`
		ClientRootName  string `json:"client_root_name"`
		SecurityToken   string `json:"security_token"`
		Expiration      string `json:"expiration"`
	}

	//StsResponse the response of sts service
	StsResponse struct {
		Credentials     Credentials     `json:"Credentials"`
		AssumedRoleUser AssumedRoleUser `json:"AssumedRoleUser"`
		RequestId       string          `json:"RequestId"`
	}
}

type Config struct {
	AccessKeyId       string `json:"AccessKeyID" ini:"accessKey"`
	AccessKeySecret   string `json:"AccessKeySecret" ini:"secretKey"`
	RoleArn           string `json:"RoleArn" ini:"roleArn"`
	TokenExpireTime   int64  `json:"TokenExpireTime" ini:"tokenExpireTime"`
	BucketReadPolicy  string `json:"BucketReadPolicy" ini:"bucketReadPolicy"`
	BucketWritePolicy string `json:"BucketWritePolicy" ini:"bucketWritePolicy"`
}

//Credentials for token app get success
type Credentials struct {
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	Expiration      string `json:"Expiration"`
	SecurityToken   string `json:"SecurityToken"`
}

//AssumedRoleUser for token app get success
type AssumedRoleUser struct {
	Arn           string `json:"Arn"`
	AssumedRoleId string `json:"AssumedRoleId"`
}

var sts Sts

// you can modify the value according to the actual situation
var stsSession string = "sts_token_server"

func InitSts(cfg *config.Config) (*Sts, error) {
	sts := new(Sts)

	var stsConfig = new(Config)
	cfg.Section("oss").MapTo(stsConfig)

	sts.newApp(stsConfig)
	fmt.Println(stsConfig)
	return sts, nil
}

func (sts Sts) GenerateSignedURL() (string, error) {
	uuidV4 := uuid.NewV4()

	queryStr := "SignatureVersion=" + StsSignVersion
	queryStr += "&Format=" + RespBodyFormat
	queryStr += "&Timestamp=" + url.QueryEscape(time.Now().UTC().Format(TimeFormat))
	queryStr += "&RoleArn=" + url.QueryEscape(sts.Config.RoleArn)
	queryStr += "&RoleSessionName=" + stsSession
	queryStr += "&AccessKeyId=" + sts.Config.AccessKeyId
	queryStr += "&SignatureMethod=HMAC-SHA1"
	queryStr += "&Version=" + StsAPIVersion
	queryStr += "&Action=AssumeRole"
	queryStr += "&SignatureNonce=" + uuidV4.String()
	queryStr += "&DurationSeconds=" + strconv.FormatUint((uint64)(sts.Config.TokenExpireTime), 10)
	if sts.Config.BucketWritePolicy != "" {
		queryStr += "&Policy=" + url.QueryEscape(sts.Config.BucketWritePolicy)
	}

	// Sort query string
	queryParams, err := url.ParseQuery(queryStr)
	if err != nil {
		return "", err
	}
	sortUrl := strings.Replace(queryParams.Encode(), "+", "%20", -1)
	strToSign := HTTPGet + "&" + PercentEncode + "&" + url.QueryEscape(sortUrl)

	// Generate signature
	hashSign := hmac.New(sha1.New, []byte(sts.Config.AccessKeySecret+"&"))
	hashSign.Write([]byte(strToSign))
	signature := base64.StdEncoding.EncodeToString(hashSign.Sum(nil))

	// Build url
	assumeURL := StsHost + "?" + queryStr + "&Signature=" + url.QueryEscape(signature)
	return assumeURL, nil
}

func (sts Sts) SendRequest(url string) ([]byte, int, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return body, resp.StatusCode, err
}

func (sts *Sts) newApp(stsConfig *Config) (alists *Sts) {
	sts.Config = stsConfig

	return sts
}
