package dji

import (
	"encoding/json"
	"fmt"
)

type DJIResult struct {
	key            DJIKeys
	sequenceNumber uint32
	value          DJIParamValue
	errorCode      int64
	errorDesc      string
}

func NewDJIResult(key DJIKeys) *DJIResult {
	return &DJIResult{
		key: key,
	}
}

func NewDJIResultFromJSON(jsonData []byte) *DJIResult {
	r := &DJIResult{}
	r.parseJSONData(jsonData)
	return r
}

func (r *DJIResult) Key() DJIKeys {
	return r.key
}

func (r *DJIResult) SequenceNumber() uint32 {
	return r.sequenceNumber
}

func (r *DJIResult) Value() DJIParamValue {
	return r.value
}

func (r *DJIResult) ErrorCode() int64 {
	return r.errorCode
}

func (r *DJIResult) ErrorDesc() string {
	return r.errorDesc
}

func (r *DJIResult) Succeeded() bool {
	return r.errorCode == 0
}

func (r *DJIResult) parseJSONData(jsonData []byte) {
	if len(jsonData) == 0 {
		r.errorCode = -1
		r.errorDesc = "empty or nil json data"
		return
	}

	var data map[string]interface{}
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		r.errorCode = -1
		r.errorDesc = fmt.Sprintf("invalid json data: %s", string(jsonData))
		return
	}

	r.sequenceNumber = uint32(data["Tag"].(float64))
	r.key = keyByValue(int(data["Key"].(float64)))
	r.errorCode = int64(data["Error"].(float64))
	r.value = data["Value"]
	if r.value == nil {
		r.errorCode = -1
	}
}
