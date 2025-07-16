package output

import (
	"FuckEventvwr/velocidex/evtx"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/blevesearch/bleve"
)

type indexEventStr struct {
	// 索引Key
	Key string `json:"Key"`
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

type Bleve struct {
	index bleve.Index
	batch *bleve.Batch
	count int
	mu    sync.Mutex
}

func NewBleve(path string) *Bleve {

	cacheDir := ".eventvwr"
	// 检查缓存目录是否存在,不存在创建
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			fmt.Println("创建缓存目录 %s 失败: %w", cacheDir, err)
			os.Exit(1)
		}
	} else if err != nil {
		fmt.Println("检查缓存目录 %s 失败: %w", cacheDir, err)
		os.Exit(1)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("获取绝对路径失败: %w", err)
		os.Exit(1)
	}
	absPath = filepath.Clean(absPath)
	hash := md5.Sum([]byte(absPath))
	var outBleve Bleve

	index, err := bleve.Open(".eventvwr/" + fmt.Sprintf("%x", hash))
	if err == bleve.ErrorIndexPathDoesNotExist {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(".eventvwr/"+fmt.Sprintf("%x", hash), mapping)
		if err != nil {
			fmt.Println("[Error] Bleve 新建索引失败:", err.Error())
			os.Exit(1)
		}
	} else if err != nil {
		fmt.Println("[Error] Bleve 打开索引失败:", err.Error())
		os.Exit(1)
	}

	outBleve.index = index
	outBleve.batch = index.NewBatch()
	outBleve.count = 0
	return &outBleve
}

func (b *Bleve) WriteRecord(record *evtx.EventRecord) error {
	struData, err := getStruData(record.Event)
	if err != nil {
		return errors.New("json解析错误: " + err.Error())
	}
	b.mu.Lock()
	defer b.mu.Unlock()

	var indexData indexEventStr
	keyStr := fmt.Sprintf("%s_%s_%d",
		struData.Event.System.Computer,
		struData.Event.System.Channel,
		struData.Event.System.EventRecordID)
	md5Bytes := md5.Sum([]byte(keyStr))
	indexData.Key = hex.EncodeToString(md5Bytes[:])
	indexData.EventID = fmt.Sprint(struData.Event.System.EventRecordID)
	indexData.Host = struData.Event.System.Computer
	indexData.Source = struData.Event.System.Channel
	indexData.EventTypeID = fmt.Sprint(struData.Event.System.EventID.Value)
	indexData.Data = struData.Event.EventData
	indexData.ProcessID = fmt.Sprint(struData.Event.System.Execution.ProcessID)

	b.batch.Index(fmt.Sprint(indexData.EventID), indexData)
	b.count++
	// 1000条批量提交
	if b.count >= 1000 {
		if err := b.index.Batch(b.batch); err != nil {
			return err
		}
		b.batch = b.index.NewBatch()
		b.count = 0
	}
	return nil
}

func (b *Bleve) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.count > 0 {
		if err := b.index.Batch(b.batch); err != nil {
			return err
		}
	}
	return b.index.Close()
}

// 写入错误
func (b *Bleve) WriteError(err string) error {
	fmt.Println("[Error] Bleve 写入错误: " + err)
	return nil
}
