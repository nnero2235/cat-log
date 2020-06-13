//log for server log requirement
//support rolling log file
//rolling log file name format: "Log-${createDate}.log"
package catlog

import (
	"fmt"
	"io"
	"os"
	"time"
)

type Logger struct {
	option        *LogOption
	consoleWriter io.WriteCloser
	fileWriter    io.WriteCloser
	today         string
	closeChan     chan struct{}
}

//TODO: 正常终止进程 如何关闭 file文件, async 未实现，maxFileSize 未实现，
//TODO: 未实现log file的 滚动 ， 未实现 定期清除 时间太久的 日志
var (
	defaultLogOption = LogOption{
		level:         INFO,
		maxFileSize:   DefaultMaxFileSize,
		logFilePath:   DefaultLogFilePath,
		fileOutput:    false,
		consoleOutput: true,
		async:         false,
	}
	logger *Logger
)

func init(){
	SetOptions()
}

func NewLogger(opts ...OptionFunc) *Logger {
	option := defaultLogOption
	for _, opt := range opts {
		opt(&option)
	}
	now := time.Now()
	lg := &Logger{
		option:  &option,
		closeChan: make(chan struct{}),
		today:   now.Format("2006-01-02"),
	}
	if option.consoleOutput {
		lg.consoleWriter = os.Stdout
	}
	if option.fileOutput {
		writer := checkAndCreateLogFiles(option.logFilePath)
		lg.fileWriter = writer
	}
	return lg
}

func SetOptions(opts ...OptionFunc) {
	if logger != nil{
		close(logger.closeChan)
	}
	logger = NewLogger(opts...)
	//another thread to check log file
	if logger.option.fileOutput {
		go logger.startCheckLogThread()
	}
}

func Debug(format string, args ...interface{}) {
	if logger.option.level <= DEBUG {
		logger.writeLogs("Debug", format, args...)
	}
}

func Trace(format string, args ...interface{}) {
	if logger.option.level <= TRACE {
		logger.writeLogs("Trace", format, args...)
	}
}

func Info(format string, args ...interface{}) {
	if logger.option.level <= INFO {
		logger.writeLogs("Info ", format, args...)
	}
}

func Error(format string, args ...interface{}) {
	if logger.option.level <= ERROR {
		logger.writeLogs("Error", format, args...)
	}
}

func Fatal(format string, args ...interface{}) {
	if logger.option.level <= FATAL {
		logger.writeLogs("Fatal", format, args...)
	}
	os.Exit(1) //fatal error cause a crash
}

func (lg *Logger) writeLogs(prefix string, format string, args ...interface{}) {
	f := fmt.Sprintf("[%s][%-s]: %s\n", prefix, time.Now().Format("2006-01-02 15:04:05.999"), format)
	if lg.consoleWriter != nil {
		if _, err := fmt.Fprintf(lg.consoleWriter, f, args...); err != nil {
			panic(err) //this must be a fatal error
		}
	}
	if lg.fileWriter != nil {
		if _, err := fmt.Fprintf(lg.fileWriter, f, args...); err != nil {
			panic(err) //this must be a fatal error
		}
	}
}

func checkAndCreateLogFiles(filePath string) io.WriteCloser {
	now := time.Now()
	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	}
	fileName := fmt.Sprintf("%s/Log-%s.log",filePath, now.Format("2006-01-02"))
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			newfile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
			if err != nil {
				panic(err)
			}
			f = newfile
		} else {
			panic(err)
		}
	}
	return f
}

func (lg *Logger)startCheckLogThread() {
	for {
		select {
		case _,ok := <-lg.closeChan:
			if ok {
				return
			}
		case <-time.Tick(time.Minute):
			newDay := time.Now().Format("2006-01-02")
			if newDay > logger.today {
				writer := checkAndCreateLogFiles(logger.option.logFilePath)
				old := logger.fileWriter
				logger.fileWriter = writer
				logger.today = newDay
				if err := old.Close();err != nil{
					Error("%v",err)
				}
			}
		}
	}
}
