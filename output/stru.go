package output

// EventStru 包含一个 Event 对象，表示一个事件的结构
type EventStru struct {
	Event Event `json:"Event"` // 事件本身，包含事件的所有详细信息
}

// Provider 表示事件的提供者，包括提供者的名称和 GUID
type Provider struct {
	Name string `json:"Name"` // 提供者的名称
	GUID string `json:"Guid"` // 提供者的 GUID
}

// EventID 表示事件的 ID，包含一个值
type EventID struct {
	Value int `json:"Value"` // 事件的唯一标识符
}

// TimeCreated 表示事件的创建时间，系统时间以浮动点数表示
type TimeCreated struct {
	SystemTime float64 `json:"SystemTime"` // 事件的系统时间
}

// Correlation 包含事件之间的关联信息，标识活动 ID
type Correlation struct {
	ActivityID string `json:"ActivityID"` // 关联的活动 ID
}

// Execution 表示事件执行的上下文，包含进程 ID 和线程 ID
type Execution struct {
	ProcessID int `json:"ProcessID"` // 执行该事件的进程 ID
	ThreadID  int `json:"ThreadID"`  // 执行该事件的线程 ID
}

// Security 表示与安全相关的事件信息，如用户 ID
type Security struct {
	UserID string `json:"UserID"` // 触发事件的用户 ID
}

// System 包含事件的基本信息，如事件的提供者、事件 ID、时间等
type System struct {
	Provider      Provider    `json:"Provider"`      // 事件的提供者
	EventID       EventID     `json:"EventID"`       // 事件的 ID
	Version       int         `json:"Version"`       // 事件版本
	Level         int         `json:"Level"`         // 事件的级别
	Task          int         `json:"Task"`          // 任务标识符
	Opcode        int         `json:"Opcode"`        // 操作码
	Keywords      interface{} `json:"Keywords"`      // 事件的关键字
	TimeCreated   TimeCreated `json:"TimeCreated"`   // 事件创建时间
	EventRecordID int         `json:"EventRecordID"` // 事件记录 ID
	Correlation   Correlation `json:"Correlation"`   // 事件的关联信息
	Execution     Execution   `json:"Execution"`     // 执行信息
	Channel       string      `json:"Channel"`       // 事件通道
	Computer      string      `json:"Computer"`      // 计算机名称
	Security      Security    `json:"Security"`      // 安全信息
}

// // Data 包含事件数据的名称和值
// type Data struct {
// 	Name  string      `json:"Name"`  // 数据项的名称
// 	Value interface{} `json:"Value"` // 数据项的值，可以是任何类型
// }

// // EventData 包含事件数据，通常与事件的具体信息相关
// type EventData struct {
// 	Data interface{} `json:"Data"` // 允许多个 Data 项
// }

// Event 表示整个事件的结构，包含系统信息和事件数据
type Event struct {
	System    System      `json:"System"`    // 事件的系统信息
	EventData interface{} `json:"EventData"` // 事件的具体数据
}
