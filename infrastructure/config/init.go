package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/task-done/infrastructure/constants"
	"gopkg.in/yaml.v3"
)

var app *AppConfig

func InitConfig() error {
	absolutePath := getTheAbsolutePath()
	configFilePath := filepath.Join(absolutePath, constants.ConfigFilePath)
	appConfig, err := parseConfig(configFilePath)
	if err != nil {
		return err
	}

	appConfig.Log.ErrLogPath = absolutePath + appConfig.Log.ErrLogPath
	appConfig.Log.SysLogPath = absolutePath + appConfig.Log.SysLogPath
	app = appConfig

	return nil
}

func GetConfig() *AppConfig {
	ret := copy(app)
	appConfig, ok := ret.(*AppConfig)
	if !ok {
		return nil
	}
	return appConfig
}

func copy(ptr interface{}) interface{} {
	// 获取指针指向的值
	value := reflect.ValueOf(ptr).Elem()

	// 分配内存
	newValue := reflect.New(value.Type())

	// 复制值到新分配的内存
	newValue.Elem().Set(value)

	// 转换并返回指针
	return newValue.Interface()
}

// 获取当前函数运行的绝对路径
func getTheAbsolutePath() string {
	var absolutePath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		absolutePath = path.Dir(filename)
	}

	absolutePath = filepath.Join(absolutePath, "../../")
	return absolutePath
}

func parseConfig(confFilePath string) (*AppConfig, error) {
	file, err := os.Open(confFilePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	decoder := yaml.NewDecoder(file)
	cfg := new(AppConfig)
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
