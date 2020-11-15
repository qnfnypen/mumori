package public

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	initialConf()
	initialLogger()
}

func initialConf() {
	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		panic("读取配置文件错误，请检查配置文件目录是否正确")
	}

	viper.WatchConfig()
}

func initialLogger() {
	var logFile *os.File

	// 判断当前系统类型
	if runtime.GOOS == "windows" {
		dir, _ := os.Getwd()
		str := viper.GetString("Logger.LogFile.Windows")
		fileDir := filepath.Dir(str)
		fileDir = fmt.Sprintf("%s/%s", dir, fileDir)
		os.MkdirAll(fileDir, os.ModePerm)
		file := fmt.Sprintf("%s/%s", dir, str)
		logFile, _ = os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	} else {
		str := viper.GetString("Logger.LogFile.Linux")
		fileDir := filepath.Dir(str)
		os.MkdirAll(fileDir, os.ModePerm)
		logFile, _ = os.OpenFile(str, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	}

	// 日志设置
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05 Monday"
	logLevel := viper.GetString("Logger.Level")
	if logLevel == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	log.Logger = zerolog.New(logFile).With().Timestamp().Caller().Logger()
}
