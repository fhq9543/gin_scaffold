package operation_log

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go_base/config"
	"go_base/utils/http"
	"go_base/utils/log"
)

//操作日志
type LogData struct {
	Module        string `json:"module"`
	Key           string `json:"key"`
	ActionFlag    int    `json:"action_flag"`
	Operator      int    `json:"operator"`
	OperatorName  string `json:"operator_name"`
	Platform      string `json:"platform"`
	ChangeMessage string `json:"change_message"`
	Remark        string `json:"remark"`
	Extras        string `json:"extras"`
}

const (
	//操作平台
	Operation = "operation"
	Supplier  = "supplier"
	Customer  = "customer"

	//集成商案例操作模块
	Aggregator = "aggregator"
	Case       = "case"

	//操作类型
	Create = 701
	Update = 702
)

func CreateOperationLog(data LogData) (success bool, error error) {
	operationHost := config.Viper.GetString("OPERATION_HOST")
	operationApi := operationHost + "/operation_log"

	r, _ := json.Marshal(data)
	b, result := http.SendRequest("POST", operationApi, "", bytes.NewBuffer(r))
	if !b {
		log.Logger.Error(fmt.Sprintf("操作日志主机连接失败，%s", result))
		return false, errors.New("操作日志主机连接失败")
	}
	return true, nil
}
