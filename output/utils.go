package output

import (
	"encoding/json"
	"fmt"
	"time"
)

// 转换 float64 类型的 SystemTime 为时间
func convertSystemTime(systemTime float64) string {
	// 将 SystemTime
	sec := int64(systemTime)
	nsec := int64((systemTime - float64(sec)) * 1e9)
	t := time.Unix(sec, nsec)

	return t.Format(time.RFC3339)
}

func getStruData(event interface{}) (EventStru, error) {
	var struData EventStru
	strData := fmt.Sprint(event)
	err := json.Unmarshal([]byte(strData), &struData)
	if err != nil {
		fmt.Println(strData)
		return struData, err
	}
	return struData, nil
}
