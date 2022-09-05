/*
Copyright 2021 The Pixiu Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Configuration struct {
	LogFile  string
	LogLevel string

	RotateMaxSize    int
	RotateMaxAge     int
	RotateMaxBackups int
	Compress         bool
}

type LoggerInterface interface {
	Info(args ...interface{})
	Infof(f string, args ...interface{})
	Error(args ...interface{})
	Errorf(f string, args ...interface{})
	Warn(args ...interface{})
	Warnf(f string, args ...interface{})
}

var (
	Logger    LoggerInterface
	AccessLog LoggerInterface
)

func Register(logDir string, logLevel string) {

	checkDir(logDir)

	var (
		accessLogFile string
		loggerLogFile string
	)

	// 支持 标准输出，标准错误输出，和指定日志文件
	switch logDir {
	case "stdout":
		accessLogFile, loggerLogFile = "stdout", "stdout"
	case "stderr":
		accessLogFile, loggerLogFile = "stderr", "stderr"
	default:
		accessLogFile, loggerLogFile = filepath.Join(logDir, "access.log"), filepath.Join(logDir, "gopixiu.log")
	}

	// 支持 INFO, WARN 和 ERROR，默认为 INFO
	Level := "info"
	if strings.ToLower(logLevel) == "error" {
		Level = "error"
	} else if strings.ToLower(logLevel) == "warn" {
		Level = "warn"
	}

	AccessLog, _ = NewZapLogger(Configuration{
		LogFile:          accessLogFile,
		LogLevel:         "info", // access 的 log 只会有 info
		RotateMaxSize:    500,
		RotateMaxAge:     7,
		RotateMaxBackups: 3,
	})

	Logger, _ = NewZapLogger(Configuration{
		LogFile:          loggerLogFile,
		LogLevel:         Level,
		RotateMaxSize:    500,
		RotateMaxAge:     7,
		RotateMaxBackups: 3,
	})
}

// checkDir 判断log_dir指定的目录是否存在，如果不存在则直接创建
func checkDir(logPath string) {
	fmt.Println(logPath)
	stat, err := os.Stat(logPath)
	// 对象不为空说明目录存在
	if stat != nil {
		return
	}
	err = os.Mkdir(logPath, 0777)
	if err != nil {
		panic("启动失败，log目录创建异常" + err.Error())
	}
}
