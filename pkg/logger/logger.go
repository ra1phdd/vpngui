package logger

import (
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, logLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), logLevel),
	)
	logger = zap.New(core, zap.AddStacktrace(zapcore.FatalLevel))
	defer logger.Sync()
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
