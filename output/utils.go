package output

import (
	"crypto/md5"
	"encoding/hex"
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

// 解析Event数据
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

// KV事件数据
type KvEventData struct {
	// 索引Key
	Key string `json:"Key"`
	// 文件名称
	FileName string `json:"FileName"`
	// 事件时间
	CreateTime string `json:"CreateTime"`
	// 事件记录ID
	EventID string `json:"EventID"`
	// 主机名字
	Host string `json:"Host"`
	// 日志来源
	Source string `json:"Source"`
	// 事件类型ID
	EventTypeID string `json:"EventTypeID"`
	// 事件具体数据
	Data interface{} `json:"Data"`
	// 进程ID
	ProcessID string `json:"ProcessID"`
}

// 解析Kv事件数据
func getKvEventData(struData EventStru, file string) KvEventData {
	var KvData KvEventData
	keyStr := fmt.Sprintf("%s|%s|%s|%d|%s|%f",
		struData.Event.System.Computer,
		struData.Event.System.Provider.Name,
		struData.Event.System.Channel,
		struData.Event.System.EventRecordID,
		file,
		struData.Event.System.TimeCreated.SystemTime)

	md5Bytes := md5.Sum([]byte(keyStr))
	KvData.Key = hex.EncodeToString(md5Bytes[:])
	KvData.FileName = file
	sec := int64(struData.Event.System.TimeCreated.SystemTime)
	nsec := int64((struData.Event.System.TimeCreated.SystemTime - float64(sec)) * 1e9)
	tm := time.Unix(sec, nsec)
	KvData.CreateTime = tm.Format("2006-01-02 15:04:05")
	KvData.EventID = fmt.Sprint(struData.Event.System.EventRecordID)
	KvData.Host = struData.Event.System.Computer
	KvData.Source = struData.Event.System.Channel
	KvData.EventTypeID = fmt.Sprint(struData.Event.System.EventID.Value)
	KvData.Data = struData.Event.EventData
	KvData.ProcessID = fmt.Sprint(struData.Event.System.Execution.ProcessID)
	return KvData
}
