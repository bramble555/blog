package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/bramble555/blog/global"
	"github.com/sirupsen/logrus"
)

const (
	red    = 31
	green  = 32
	yellow = 33
	blue   = 36
	gray   = 37
)

type LogFormatter struct{}

// 自定义格式实现这个Formatter接口
func (l *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case logrus.TraceLevel, logrus.DebugLevel:
		levelColor = blue
	case logrus.InfoLevel, logrus.WarnLevel:
		levelColor = green
	default:
		levelColor = red
	}
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	// entry.HasCaller() 来决定是否添加调用者信息
	if entry.HasCaller() {
		funVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		// 输出格式
		fmt.Fprintf(b, "%s[%s] \x1b[%dm[%s]\x1b[0m  %s %s \x1b[%dm[%s]\x1b[0m\n",
			global.Config.Logger.Prefix,
			timestamp,
			levelColor,
			entry.Level,
			fileVal,
			funVal,
			levelColor,
			entry.Message,
		)
	} else {
		fmt.Fprintf(b, "%s[%s] \x1b[%dm[%s]\x1b[0m  %s\n",
			global.Config.Logger.Prefix,
			timestamp,
			levelColor,
			entry.Level,
			entry.Message,
		)
	}
	return b.Bytes(), nil
}

// 判断log目录是否存在
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		log.Println(err)
		return false
	}
	return true
}
func createLogFile() *os.File {
	ok := pathExists("logger")
	// 目录不存在就创建目录
	if !ok {
		err := os.MkdirAll("logger", 0777)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// 输出到哪
	logFilePath := filepath.Join("logger", global.Config.Logger.FileName)
	// 创建文件（如果不存在）
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err.Error())
	}
	return file
}
func Init() (mLog *logrus.Logger, err error) {
	mLog = logrus.New()
	// 设置输出为控制台和文件
	file := createLogFile()
	mLog.SetOutput(io.MultiWriter(os.Stdout, file))
	mLog.SetReportCaller(global.Config.Logger.ShowLine)
	// 格式
	mLog.SetFormatter(&LogFormatter{})
	// 解析level
	level, err := logrus.ParseLevel(global.Config.Logger.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	// 设置level
	mLog.SetLevel(level)
	return
}
