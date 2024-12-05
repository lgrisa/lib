package mgr

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/lgrisa/lib/utils"
	"net/http"
)

func writeConfigJson(w http.ResponseWriter, data *ConfigJson, cmdOutPut string) {
	if err := jsoniter.NewEncoder(w).Encode(&response{
		Code:    200,
		Data:    data,
		Message: cmdOutPut,
	}); err != nil {
		utils.LogErrorF("writeConfigJson NewEncoder fail: %v", err)
		return
	}
}

func writeErrMsg(w http.ResponseWriter, s string) {
	_ = jsoniter.NewEncoder(w).Encode(&response{
		Code:    400,
		Message: s,
	})

	fmt.Println(s)
}

func writeCsJson(w http.ResponseWriter, data *ConfigCsJson, cmdOutPut string) {
	if err := jsoniter.NewEncoder(w).Encode(&response{
		Code:    200,
		Data:    data,
		Message: cmdOutPut,
	}); err != nil {
		return
	}
}
