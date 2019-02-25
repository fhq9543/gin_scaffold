package http

import (
	"crypto/tls"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"go_base/utils/log"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

func SendRequest(method string, url string, token string, body io.Reader) (bool, interface{}) {
	req, err := http.NewRequest(method, url, body)
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	//忽略证书校验
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return false, "http请求失败:" + err.Error()
	}
	defer resp.Body.Close()

	result, err := ioutil.ReadAll(resp.Body)
	js, err := simplejson.NewJson(result)
	if js.Get("rescode").MustString() == "10000" {
		data, _ := js.Get("data").MarshalJSON()
		log.Logger.Debug(string(data))
		return true, data
	}
	return false, js.Get("msg").MustString()
}

func SetPagingHeader(c *gin.Context, offset int, len int, total int) {
	content_range := fmt.Sprintf("%d-%d", offset, offset+len)
	c.Header("X-Content-Range", content_range)
	c.Header("X-Content-Total", strconv.Itoa(total))
}
