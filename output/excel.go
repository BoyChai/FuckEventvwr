package output

import (
	"FuckEventvwr/config"
	"FuckEventvwr/velocidex/evtx"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/xuri/excelize/v2"
)

var logIndexLock sync.Mutex

type Excel struct {
	file           *excelize.File
	logWriter      *excelize.StreamWriter
	originalWriter *excelize.StreamWriter
	errorWriter    *excelize.StreamWriter
	logIndex       int
}

func NewExcel() *Excel {
	f := excelize.NewFile()
	err := f.SaveAs(config.Cfg.Output)
	if err != nil {
		fmt.Println("[Error] 指定输出文件有错误,请检查:", err.Error())
		os.Exit(1)
	}
	// 设置Log工作表
	index, err := f.NewSheet("Log")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	f.SetActiveSheet(index)
	logWrite, err := f.NewStreamWriter("Log")
	if err != nil {
		fmt.Println("[Error] 创建工作簿出现错误,请检查:", err.Error())
		os.Exit(1)
	}
	// 设置列名
	if err := logWrite.SetColWidth(1, 16, 15); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	header := []interface{}{}
	for _, call := range []string{
		"计算机", "日志来源", "事件ID", "版本", "级别", "任务类型", "操作码", "关键字", "事件时间", "事件记录ID", "关联活动ID", "执行的进程ID", "执行的线程ID", "事件通道", "安全信息", "事件数据",
	} {
		header = append(header, call)
	}
	if err := logWrite.SetRow("A1", header); err != nil {
		fmt.Println("[Error] 设置Log工作簿出现错误,请检查:", err.Error())
		os.Exit(1)
	}

	// 设置Original工作表
	index, err = f.NewSheet("Original")
	if err != nil {
		fmt.Println(err)
	}
	f.SetActiveSheet(index)
	originalWrite, err := f.NewStreamWriter("Original")
	if err != nil {
		fmt.Println("[Error] 创建工作簿出现错误,请检查:", err.Error())
		os.Exit(1)
	}

	// 设置Error工作表
	index, err = f.NewSheet("Error")
	if err != nil {
		fmt.Println(err)
	}
	f.SetActiveSheet(index)
	errWrite, err := f.NewStreamWriter("Error")
	if err != nil {
		fmt.Println("[Error] 创建工作簿出现错误,请检查:", err.Error())
		os.Exit(1)
	}

	err = f.SaveAs(config.Cfg.Output)
	if err != nil {
		fmt.Println("[Error] 指定输出文件有错误,请检查:", err.Error())
		os.Exit(1)
	}
	return &Excel{
		file:           f,
		logWriter:      logWrite,
		originalWriter: originalWrite,
		errorWriter:    errWrite,
		logIndex:       2,
	}
}

func (e *Excel) WriteRecord(record *evtx.EventRecord) error {
	var struData EventStru

	strData := fmt.Sprint(record.Event)
	err := json.Unmarshal([]byte(strData), &struData)
	if err != nil {
		fmt.Println(strData)
		return errors.New("json解析错误: " + err.Error())
	}
	logIndexLock.Lock()
	logData := []interface{}{}
	for _, call := range []string{
		struData.Event.System.Computer,
		struData.Event.System.Provider.Name,
		fmt.Sprint(struData.Event.System.EventID.Value),
		fmt.Sprint(struData.Event.System.Version),
		fmt.Sprint(struData.Event.System.Level),
		fmt.Sprint(struData.Event.System.Task),
		fmt.Sprint(struData.Event.System.Opcode),
		fmt.Sprint(struData.Event.System.Keywords),
		convertSystemTime(struData.Event.System.TimeCreated.SystemTime),
		fmt.Sprint(struData.Event.System.EventRecordID),
		fmt.Sprint(struData.Event.System.Correlation.ActivityID),
		fmt.Sprint(struData.Event.System.Execution.ProcessID),
		fmt.Sprint(struData.Event.System.Execution.ThreadID),
		fmt.Sprint(struData.Event.System.Channel),
		fmt.Sprint(struData.Event.System.Security),
		fmt.Sprint(struData.Event.EventData.Data),
	} {
		logData = append(logData, call)
	}
	if err := e.logWriter.SetRow(fmt.Sprint("A", e.logIndex), logData); err != nil {
		fmt.Println("[Error] 写入Log工作簿出现错误,请检查:", err.Error())
	}
	e.logIndex++
	defer logIndexLock.Unlock()
	return nil
}

func (e *Excel) WriteError(err string) error {
	return nil
}

func (e *Excel) Close() error {
	// 刷新所有流式写入器
	if err := e.logWriter.Flush(); err != nil {
		fmt.Println("[Error] 刷新Log工作簿流式写入器时出现错误:", err.Error())
		return err
	}
	// 设置筛选
	if err := e.file.AddTable("Log", &excelize.Table{
		Range: fmt.Sprint("A1:P", e.logIndex),
		Name:  "table",
	}); err != nil {
		fmt.Println("[Error] 设置Log工作簿出现错误,请检查:", err.Error())
	}

	if err := e.originalWriter.Flush(); err != nil {
		fmt.Println("[Error] 刷新Original工作簿流式写入器时出现错误:", err.Error())
		return err
	}

	if err := e.errorWriter.Flush(); err != nil {
		fmt.Println("[Error] 刷新Error工作簿流式写入器时出现错误:", err.Error())
		return err
	}

	// 最终保存文件
	if err := e.file.SaveAs(config.Cfg.Output); err != nil {
		fmt.Println("[Error] 保存文件时出现错误:", err.Error())
		return err
	}

	return nil
}
