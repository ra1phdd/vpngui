package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

var logger *zap.Logger
var logLevel zap.AtomicLevel

func Init() {
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeLevel = zapcore.LowercaseLevelEncoder
	config.EncodeTime = customTimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		log.Fatal("Ошибка создания папки logs", err.Error())
	}
	logFile, err := os.OpenFile("logs/main.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Ошибка создания файла main.log", err.Error())
	}
	writer := zapcore.AddSync(logFile)

	logLevel = zap.NewAtomicLevel()
	logLevel.SetLevel(zapcore.InfoLevel)

	customWriter := &customSyncer{writer: os.Stdout}

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, logLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(customWriter), logLevel),
	)
	logger = zap.New(core, zap.AddStacktrace(zapcore.FatalLevel))
}

func SetLogLevel(level string) {
	switch level {
	case "debug":
		logLevel.SetLevel(zapcore.DebugLevel)
	case "warn":
		logLevel.SetLevel(zapcore.WarnLevel)
	case "error":
		logLevel.SetLevel(zapcore.ErrorLevel)
	case "fatal":
		logLevel.SetLevel(zapcore.FatalLevel)
	case "info":
		logLevel.SetLevel(zapcore.InfoLevel)
	default:
		logLevel.SetLevel(zapcore.InfoLevel)
	}
}

var (
	OutLog []string
	MuLog  sync.Mutex
)

type customSyncer struct {
	mu     sync.Mutex
	writer io.Writer
}

func (c *customSyncer) Write(p []byte) (n int, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	line := string(p)
	OutLog = append(OutLog, line)

	return c.writer.Write(p)
}

func (c *customSyncer) Sync() error {
	return nil
}

func Debug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	logger.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}
