package response

import (
	"baseFrame/pkg/config"
	"baseFrame/pkg/logger"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type ListResp struct {
	Results   interface{} `json:"results"`
	Count     int64       `json:"count"`
	Limit     int         `json:"limit"`
	Offset    int         `json:"offset"`
	Page      int         `json:"page"`
	PageSize  int         `json:"page_size"`
	CurPage   int         `json:"cur_page"`
	TotalPage int         `json:"total_page"`
}

type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SearchField(t reflect.Type, FieldName string) (fields []reflect.StructField) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 过滤已处理的类型, 避免类型环形嵌套造成死循环
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			fs := SearchField(field.Type, FieldName)
			fields = append(fields, fs...)
		}
		if field.Name != FieldName {
			continue
		}
		fields = append(fields, field)
	}
	return fields
}

func SearchTag(t reflect.Type, TagName string) (fields []reflect.StructField) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// 过滤已处理的类型, 避免类型环形嵌套造成死循环
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			fs := SearchTag(field.Type, TagName)
			fields = append(fields, fs...)
		}
		_, ok := field.Tag.Lookup(TagName)
		if !ok {
			continue
		}
		fields = append(fields, field)
	}
	return fields
}

type ListQueryParams struct {
	All      uint     `json:"all" form:"all"` // 不做分页
	Page     int      `json:"page" form:"page"`
	PageSize int      `json:"page_size" form:"page_size"`
	Limit    int      `json:"limit" form:"limit"`
	Offset   int      `json:"offset" form:"offset"`
	Search   string   `json:"search" form:"search"`
	SortBy   []string `json:"sort_by" form:"sort_by"`
}

func GetList(ctx *gin.Context, db *gorm.DB, results interface{}, processes ...func(interface{}) error) (interface{}, error) {
	var count int64

	queryParams := &ListQueryParams{}
	if !logger.Check(ctx.ShouldBindQuery(&queryParams)) {
		queryParams.Offset = 0
		queryParams.Limit = 10
	}

	if queryParams.PageSize == 0 {
		queryParams.PageSize = 10
	}

	// 当limit，offset都为0时，则走page跟page_size
	if queryParams.Limit == 0 && queryParams.Offset == 0 {
		queryParams.Limit = queryParams.PageSize
		queryParams.Offset = (queryParams.Page - 1) * queryParams.Limit
	}

	if queryParams.Offset < 0 {
		queryParams.Offset = 0
	}

	if queryParams.Limit == 0 {
		queryParams.Limit = 10
	}

	t := reflect.TypeOf(results)
	if queryParams.Search != "" {
		fields := SearchTag(t, "search")
		var query []string
		var values []interface{}

		for _, field := range fields {
			searchField := field.Tag.Get("search")
			query = append(query, fmt.Sprintf(" (`%s` LIKE ?) ", searchField))
			values = append(values, "%"+queryParams.Search+"%")
		}
		if len(query) > 0 {
			queryStr := strings.Join(query, "OR")
			db = db.Where(queryStr, values...)
		}
	}

	fields := SearchTag(t, "json")
	if strings.Join(queryParams.SortBy, "") != "" {
		for _, sortByOrder := range queryParams.SortBy {
			parts := strings.Split(sortByOrder, " ")
			if len(parts) < 1 {
				continue
			}
			sortBy := parts[0]
			order := "desc"
			if len(parts) > 1 && parts[1] == "asc" {
				order = "asc"
			}
			for _, field := range fields {
				if field.Tag.Get("json") != sortBy {
					continue
				}
				stmt := &gorm.Statement{DB: db}
				err := stmt.Parse(results)
				if err == nil {
					db = db.Order(fmt.Sprintf("%s.`%s` %s", stmt.Table, sortBy, order))
				}
				break
			}
		}
	} else {
		hasSort := false
		hasID := false
		for _, field := range fields {
			if field.Tag.Get("json") == "id" {
				hasID = true
			}
			if field.Tag.Get("json") == "sort" {
				hasSort = true
			}
			if hasID && hasSort {
				break
			}
		}

		stmt := &gorm.Statement{DB: db}
		err := stmt.Parse(results)
		if err == nil {
			if hasSort {
				db = db.Order(fmt.Sprintf("%s.sort desc", stmt.Table))
			}
			if hasID {
				db = db.Order(fmt.Sprintf("%s.id desc", stmt.Table))
			}
		}
	}

	if queryParams.All == 1 {
		// 不走分页
		if err := db.Debug().Find(results).Count(&count).Error; err != nil {
			return nil, err
		}
	} else {
		// 走分页
		if err := db.Debug().Limit(queryParams.Limit).Offset(queryParams.Offset).Find(results).
			Limit(-1).Offset(-1).Count(&count).Error; err != nil {
			return nil, err
		}
	}

	if len(processes) > 0 {
		for _, p := range processes {
			err := p(results)
			if err != nil {
				return nil, err
			}
		}
	}

	if queryParams.All == 1 {
		// 不走分页
		return results, nil
	}

	return &ListResp{
		Results:   results,
		Count:     count,
		Limit:     queryParams.Limit,
		Offset:    queryParams.Offset,
		Page:      queryParams.Page,
		PageSize:  queryParams.PageSize,
		CurPage:   queryParams.Offset/queryParams.PageSize + 1,
		TotalPage: int(float32(count)/float32(queryParams.PageSize)) + 1,
	}, nil
}

type Response struct {
	db  *gorm.DB
	cfg *config.Config
}

func Success(ctx *gin.Context, msg string, data interface{}) {
	resp := response{
		Code: 0,
		Msg:  msg,
		Data: data,
	}
	//logger.Infof("response success:%v", resp)
	ctx.JSON(http.StatusOK, resp)
}

func Error(ctx *gin.Context, code int, msg string, data interface{}) {
	resp := response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	logger.DebugWithCtx(ctx, "response error:", resp)

	ctx.JSON(http.StatusOK, resp)
}

func ErrorBadRequest(ctx *gin.Context, code int, msg string, data interface{}) {
	resp := response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	logger.DebugWithCtx(ctx, "response error:", resp)

	ctx.JSON(http.StatusBadRequest, resp)
}

func ReturnList(ctx *gin.Context, db *gorm.DB, results interface{}, processes ...func(interface{}) error) {
	list, err := GetList(ctx, db, results, processes...)
	if !logger.Check(err) {
		Error(ctx, 20001, "查询失败", nil)
		return
	}
	Success(ctx, "", list)
}

func InitResponse(db *gorm.DB, cfg *config.Config) (*Response, error) {
	var rsp Response
	rsp.db = db
	rsp.cfg = cfg
	return &rsp, nil
}
