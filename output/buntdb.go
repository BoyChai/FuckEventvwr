package output

import (
	"FuckEventvwr/cmd"
	"FuckEventvwr/config"
	"FuckEventvwr/velocidex/evtx"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tidwall/buntdb"
)

func init() {
	var moduleConfig cmd.ModuleConfig
	moduleConfig.Description = "默认处理方式，将事件记录写入本地缓存数据库，写入和分析速度速度较快"
	moduleConfig.Params = []cmd.ModuleParam{
		{
			Name:        "path",
			Default:     "C:\\Windows\\System32\\winevt\\Logs",
			Description: "指定事件日志目录，默认为 C:\\Windows\\System32\\winevt\\Logs",
			FlagPtr:     new(string),
		},
	}
	moduleConfig.Apply = func(flags map[string]any) {
		config.Cfg.Path = *flags["path"].(*string)
		Module = NewBuntDB(*flags["path"].(*string))
	}
	cmd.RegisterModule("buntdb", moduleConfig)
}

type Buntdb struct {
	db   *buntdb.DB
	path string
}

// BuntOutput
func NewBuntDB(path string) *Buntdb {
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

	return &Buntdb{
		db:   db,
		path: dbPath,
	}
}

// 写入事件记录
func (b *Buntdb) WriteRecord(record *evtx.EventRecord) error {
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

// WriteError 方法用于写入错误信息
func (b *Buntdb) WriteError(err string) error {
	fmt.Println("[Error] BuntDB 写入错误:", err)
	return nil
}

func (b *Buntdb) Close() (int, error) {
	var keyCount int

	err := b.db.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys("*", func(key, value string) bool {
			keyCount++
			return true
		})
	})
	if err != nil {
		return 0, err
	}
	return keyCount, b.db.Close()
}

func (b *Buntdb) Count() (int, error) {
	var keyCount int

	err := b.db.View(func(tx *buntdb.Tx) error {
		return tx.AscendKeys("*", func(key, value string) bool {
			keyCount++
			return true
		})
	})
	if err != nil {
		return 0, err
	}
	return keyCount, nil
}

