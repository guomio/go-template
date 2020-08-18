package logger

import (
	"os"
	"time"

	"github.com/guomio/go-template/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger zap
type Logger struct {
	log *zap.Logger
}

// NewLoggerOption logger 配置
type NewLoggerOption struct {
	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	Filename string
	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int
	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int
	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int
	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool
	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool
}

// NewLogger 初始化 loagger
func NewLogger(opt *NewLoggerOption) *Logger {
	out := zapcore.AddSync(&lumberjack.Logger{
		Filename:   opt.Filename,
		MaxSize:    opt.MaxSize,
		MaxBackups: opt.MaxBackups,
		MaxAge:     opt.MaxAge,
		LocalTime:  opt.LocalTime,
		Compress:   opt.Compress,
	})
	stdout := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(newProductionEncoderConfig()),
		zapcore.NewMultiWriteSyncer(stdout, out),
		zap.InfoLevel,
	)
	logger := &Logger{log: zap.New(core)}
	defer logger.log.Sync()
	return logger
}

// Sprintf uses fmt.Sprintf to log a templated message.
func (l *Logger) Sprintf(template string, args ...interface{}) {
	l.log.Sugar().Infof(template, args...)
}

// Printf logs a message with some additional context.
func (l *Logger) Printf(template string, args ...interface{}) {
	l.log.Sugar().Infow(template, args...)
}

// Sprint uses fmt.Sprint to construct and log a message.
func (l *Logger) Sprint(args ...interface{}) {
	l.log.Sugar().Info(args...)
}

// Debug logs a message at DebugLevel.
func (l *Logger) Debug(msg string, fields ...zapcore.Field) {
	l.log.Debug(msg, fields...)
}

// Info logs a message at InfoLevel.
func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	l.log.Info(msg, fields...)
}

// Warn logs a message at WarnLevel.
func (l *Logger) Warn(msg string, fields ...zapcore.Field) {
	l.log.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel.
func (l *Logger) Error(msg string, fields ...zapcore.Field) {
	l.log.Error(msg, fields...)
}

// String constructs a field with the given key and value.
func String(key string, val string) zapcore.Field {
	return zap.String(key, val)
}

// Bool constructs a field that carries a bool.
func Bool(key string, val bool) zapcore.Field {
	return zap.Bool(key, val)
}

// Int constructs a field with the given key and value.
func Int(key string, val int) zapcore.Field {
	return zap.Int(key, val)
}

// Any takes a key and an arbitrary value and chooses the best way to represent them as a field, falling back to a reflection-based approach only if necessary.
func Any(key string, value interface{}) zapcore.Field {
	return zap.Any(key, value)
}

func newProductionEncoderConfig() zapcore.EncoderConfig {
	epochTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006/01/02 15:04:05"))
	}
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     epochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// 初始化部分

// L Logger 实例
var L *Logger

// Init 初始化
func Init() {
	config := config.GetConfig()
	L = NewLogger(&NewLoggerOption{Filename: config.Log})
}

// GetLogger 获取 Logger 实例
func GetLogger() *Logger {
	return L
}
