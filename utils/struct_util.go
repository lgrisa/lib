package utils

import "encoding/json"

// Struct2Map struct to map，依赖 json tab
func Struct2Map(r interface{}) (s map[string]string, err error) {
	var temp map[string]interface{}
	var result = make(map[string]string)

	bin, err := json.Marshal(r)
	if err != nil {
		return result, err
	}
	if err := json.Unmarshal(bin, &temp); err != nil {
		return nil, err
	}
	for k, v := range temp {
		result[k], err = ToStringE(v)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
