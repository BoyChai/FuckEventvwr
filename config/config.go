package config

type config struct {
	// 路径
	Path string
	// 输出
	Output string
	// 输出模式
	Mode int
}

var Cfg = &config{}
