package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/pkg/file"
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

func Init() (mLog *logrus.Logger, err error) {
	mLog = logrus.New()
	// 设置输出为控制台和文件
	logFile := file.CreateFile(global.Config.Logger.FilePath, global.Config.Logger.FileName)
	mLog.SetOutput(io.MultiWriter(os.Stdout, logFile))
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
