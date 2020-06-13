package catlog

type LogLevel int

const (
	DefaultLogFilePath = "logs"
	DefaultMaxFileSize int64 = 128 * 1024 * 1024

	DEBUG LogLevel = 1
	TRACE LogLevel = 2
	INFO  LogLevel = 3
	ERROR LogLevel = 4
	FATAL LogLevel = 5
)


type LogOption struct {
	consoleOutput bool
	fileOutput    bool
	logFilePath   string
	level         LogLevel
	maxFileSize   int64
	async         bool
}

type OptionFunc func(option *LogOption)

func WithConsoleOutput(console bool) OptionFunc {
	return func(option *LogOption){
		option.consoleOutput = console
	}
}

func WithFileOutput(fileOutput bool) OptionFunc {
	return func(option *LogOption){
		option.fileOutput = fileOutput
	}
}

func WithLogFilePath(logFilePath string) OptionFunc {
	return func(option *LogOption){
		option.logFilePath = logFilePath
	}
}

func WithLevel(level LogLevel) OptionFunc {
	return func(option *LogOption){
		option.level = level
	}
}

func WithMaxFileSize(size int64) OptionFunc {
	return func(option *LogOption){
		option.maxFileSize = size
	}
}

func WithAsync(async bool) OptionFunc {
	return func(option *LogOption) {
		option.async = async
	}
}

