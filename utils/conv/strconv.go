package conv

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func Atoi64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func Atou(s string) (uint, error) {
	u64, err := strconv.ParseUint(s, 10, 0)
	return uint(u64), err
}

func Atou64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

func Atob(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func Atof(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func I64toa(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Utoa(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}

func U64toa(u uint64) string {
	return strconv.FormatUint(u, 10)
}

func Btoa(b bool) string {
	return strconv.FormatBool(b)
}

func Ftoa(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

// 解析域名，去除前缀
func UrlParse(urlpath string) (res string, err error) {
	u, err := url.Parse(urlpath)
	if err != nil {
		return "", err
	}
	res = strings.Trim(u.Path, "/")
	return res, nil
}

// 解析域名，去除前缀
func UrlListParse(urlpath []string) (err error) {
	for i := 0; i < len(urlpath); i++ {
		u, err := url.Parse(urlpath[i])
		if err != nil {
			return err
		}
		urlpath[i] = strings.Trim(u.Path, "/")
	}
	return nil
}

// 转换
func Map2Str(mapData interface{}) (result string, err error) {
	resultByte, errError := json.Marshal(mapData)
	result = string(resultByte)
	err = errError
	return result, err
}

func Str2Map(jsonData string) (result map[string]interface{}, err error) {
	err = json.Unmarshal([]byte(jsonData), &result)
	return result, err
}

func Str2Stu(strData string, stu *interface{}) (err error) {
	if err := json.Unmarshal([]byte(strData), &stu); err != nil {
		return err
	}
	return nil
}

func Stu2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		//log.LogPrint(getStructTag(t.Field(i)))
		//log.LogPrint(t.Field(i).Name)
		data[getStructTag(t.Field(i))] = v.Field(i).Interface()
	}
	return data
}

func getStructTag(f reflect.StructField) string {
	return string(f.Tag.Get("json"))
}
