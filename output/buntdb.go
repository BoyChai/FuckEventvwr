package output

import (
	"FuckEventvwr/velocidex/evtx"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tidwall/buntdb"
)

type BuntOutput struct {
	db   *buntdb.DB
	path string
}

// BuntOutput
func NewBuntOutput(path string) *BuntOutput {
	cacheDir := ".eventvwr"
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			fmt.Printf("创建缓存目录 %s 失败: %v\n", cacheDir, err)
			os.Exit(1)
		}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("获取绝对路径失败: %v\n", err)
		os.Exit(1)
	}
	absPath = filepath.Clean(absPath)
	hash := md5.Sum([]byte(absPath))
	dbPath := filepath.Join(cacheDir, fmt.Sprintf("%x.db", hash))

	db, err := buntdb.Open(dbPath)
	if err != nil {
		fmt.Println("打开 buntdb 数据库失败:", err)
		os.Exit(1)
	}

	// 索引
	_ = db.CreateIndex("event_type_id", "event:*", buntdb.IndexJSON("EventTypeID"))

	return &BuntOutput{
		db:   db,
		path: dbPath,
	}
}

// 写入事件记录
func (b *BuntOutput) WriteRecord(record *evtx.EventRecord) error {
	struData, err := getStruData(record.Event)
	if err != nil {
		return errors.New("json解析错误: " + err.Error())
	}

	kvData := getKvEventData(struData, record.FileName)

	value, err := json.Marshal(kvData)
	if err != nil {
		return err
	}

	key := "event:" + kvData.Key

	return b.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, string(value), nil)
		return err
	})
}

func (b *BuntOutput) WriteError(err string) error {
	fmt.Println("[Error] BuntDB 写入错误:", err)
	return nil
}

func (b *BuntOutput) Close() error {

	return b.db.Close()
}
