package output

import (
	"FuckEventvwr/config"
	"FuckEventvwr/velocidex/evtx"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/olivere/elastic/v7"
)

type ESDataStru struct {
	Compute    string `json:"计算机"`
	Channel    string `json:"日志来源"`
	EventID    string `json:"事件ID"`
	Level      int    `json:"事件级别"`
	Version    int    `json:"事件系统版本"`
	Task       int    `json:"任务类型"`
	Opcode     int    `json:"操作码"`
	Keywords   string `json:"关键字"`
	CreateTime string `json:"事件时间"`
	RecordID   int    `json:"记录ID"`
	ActivityID string `json:"活动ID"`
	ProcessID  int    `json:"进程ID"`
	ThreadID   int    `json:"线程ID"`
	Security   string `json:"安全信息"`
	EventData  string `json:"事件数据"`
}

func (ESDataStru) Mapping() string {
	return `
{
  "mappings": {
    "properties": {
      "计算机": {
        "type": "keyword"
      },
      "日志来源": {
        "type": "keyword"
      },
      "事件ID": {
        "type": "keyword"
      },
      "事件级别": {
        "type": "integer"
      },
      "事件系统版本": {
        "type": "integer"
      },
      "任务类型": {
        "type": "integer"
      },
      "操作码": {
        "type": "integer"
      },
      "关键字": {
        "type": "keyword"
      },
      "事件时间": {
        "type": "date",
        "format": "strict_date_time"
      },
      "记录ID": {
        "type": "integer"
      },
      "活动ID": {
        "type": "keyword"
      },
      "进程ID": {
        "type": "integer"
      },
      "线程ID": {
        "type": "integer"
      },
      "安全信息": {
        "type": "text"
      },
      "事件数据": {
        "type": "text"
      }
    }
  }
}

	`
}

type Elasticsearch struct {
	Client *elastic.Client
	Index  *elastic.IndexService
}

// 缓冲
var recordBuffer []ESDataStru
var bufferLimit = 1000
var bufferLock sync.Mutex

// 创建一个新的Elasticsearch实例
func NewElasticsearch() *Elasticsearch {
	client, err := elastic.NewClient(
		elastic.SetURL(config.Cfg.EsURL),
		elastic.SetSniff(false),
	)
	if err != nil {
		fmt.Println("[Error] 连接Elasticsearch失败,请检查:", err.Error())
		os.Exit(1)
	}
	_, err = client.CreateIndex(config.Cfg.Output).BodyString(ESDataStru{}.Mapping()).Do(context.Background())
	if err != nil {
		fmt.Println("[Error] 创建Elasticsearch索引失败,请检查:", err.Error())
		os.Exit(1)
	}

	return &Elasticsearch{Client: client, Index: client.Index().Index(config.Cfg.Output)}
}
func (e *Elasticsearch) WriteRecord(record *evtx.EventRecord) error {
	var struData EventStru
	strData := fmt.Sprint(record.Event)
	err := json.Unmarshal([]byte(strData), &struData)
	if err != nil {
		fmt.Println(strData)
		return errors.New("json解析错误: " + err.Error())
	}

	data := ESDataStru{
		Compute:    struData.Event.System.Computer,
		Channel:    struData.Event.System.Channel,
		EventID:    fmt.Sprint(struData.Event.System.EventID.Value),
		Level:      struData.Event.System.Level,
		Version:    struData.Event.System.Version,
		Task:       struData.Event.System.Task,
		Opcode:     struData.Event.System.Opcode,
		Keywords:   fmt.Sprint(struData.Event.System.Keywords),
		CreateTime: convertSystemTime(struData.Event.System.TimeCreated.SystemTime),
		RecordID:   struData.Event.System.EventRecordID,
		ActivityID: struData.Event.System.Correlation.ActivityID,
		ProcessID:  struData.Event.System.Execution.ProcessID,
		ThreadID:   struData.Event.System.Execution.ThreadID,
		Security:   fmt.Sprint(struData.Event.System.Security),
		EventData:  fmt.Sprint(struData.Event.EventData),
	}

	// 将数据加入缓存
	bufferLock.Lock()
	recordBuffer = append(recordBuffer, data)

	// 当缓存达到一定大小时，执行批量写入
	if len(recordBuffer) >= bufferLimit {
		go e.do()
	} else {
		bufferLock.Unlock()
	}
	return nil
}
func (e *Elasticsearch) WriteError(err string) error {
	fmt.Println("[Error] " + err)
	return nil
}
func (e *Elasticsearch) Close() error {
	bufferLock.Lock()
	e.do()
	e.Client.Stop()
	return nil
}

func (e *Elasticsearch) do() {
	bulkRequest := e.Client.Bulk()
	for _, record := range recordBuffer {
		request := elastic.NewBulkIndexRequest().
			Index(config.Cfg.Output).
			Doc(record).
			OpType("index")
		bulkRequest = bulkRequest.Add(request)
	}

	_, err := bulkRequest.Do(context.Background())
	if err != nil {
		fmt.Println("[Error] 批量写入Elasticsearch失败, 请检查:", err.Error())
	}

	// 清空缓存
	recordBuffer = nil
	bufferLock.Unlock()
}
