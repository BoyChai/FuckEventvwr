package config

type config struct {
	// 路径
	Path string
	// 输出
	Output string
	// 是否覆盖
	Cover bool
}

var Cfg = &config{}
