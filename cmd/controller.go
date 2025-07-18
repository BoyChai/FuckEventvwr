package cmd

import (
	"flag"
	"fmt"
	"os"
)

type ModuleParam struct {
	Name        string
	Default     any
	Description string
	FlagPtr     any
}

type ModuleConfig struct {
	Params      []ModuleParam
	Description string
	Apply       func(flags map[string]any)
}

var moduleRegistry = make(map[string]ModuleConfig)

// 模块动态注册
func RegisterModule(name string, config ModuleConfig) {
	moduleRegistry[name] = config
}

// 初始化模块
func RunModule(args []string) {
	if len(args) < 2 {
		fmt.Println("未指定模块，默认采用buntdb分析")
		args[1] = "buntdb"
	}
	mode := args[1]

	module, exists := moduleRegistry[mode]
	if !exists {
		fmt.Println("未知模块:", mode)
		fmt.Println("可用子命令:", getAvailableModules())
		os.Exit(1)
	}

	fs := flag.NewFlagSet(mode, flag.ExitOnError)
	flags := make(map[string]any)

	for _, param := range module.Params {
		switch ptr := param.FlagPtr.(type) {
		case *string:
			fs.StringVar(ptr, param.Name, param.Default.(string), param.Description)
			flags[param.Name] = ptr
		case *int:
			fs.IntVar(ptr, param.Name, param.Default.(int), param.Description)
			flags[param.Name] = ptr
		}
	}

	fs.Parse(args[2:])

	module.Apply(flags)
}

// 模块信息
func getAvailableModules() string {
	modules := make([]string, 0, len(moduleRegistry))
	for name := range moduleRegistry {
		modules = append(modules, name)
	}
	return fmt.Sprint(modules)
}
