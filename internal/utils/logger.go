package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// LogLevel 日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// Logger 日志记录器
type Logger struct {
	level    LogLevel
	logger   *log.Logger
	file     *os.File
	filePath string
}

// NewLogger 创建新的日志记录器
func NewLogger(level LogLevel, logToFile bool) (*Logger, error) {
	logger := &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}

	if logToFile {
		// 创建日志目录
		logDir := "logs"
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// 创建日志文件
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		logger.filePath = filepath.Join(logDir, fmt.Sprintf("genpulse_%s.log", timestamp))

		file, err := os.OpenFile(logger.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		logger.file = file
		logger.logger = log.New(file, "", log.LstdFlags)
	}

	return logger, nil
}

// Close 关闭日志记录器
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// logInternal 内部日志方法
func (l *Logger) logInternal(level LogLevel, levelStr string, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	message := fmt.Sprintf(format, args...)
	l.logger.Printf("[%s] %s", levelStr, message)
}

// Debug 调试日志
func (l *Logger) Debug(format string, args ...interface{}) {
	l.logInternal(DEBUG, "DEBUG", format, args...)
}

// Info 信息日志
func (l *Logger) Info(format string, args ...interface{}) {
	l.logInternal(INFO, "INFO", format, args...)
}

// Warn 警告日志
func (l *Logger) Warn(format string, args ...interface{}) {
	l.logInternal(WARN, "WARN", format, args...)
}

// Error 错误日志
func (l *Logger) Error(format string, args ...interface{}) {
	l.logInternal(ERROR, "ERROR", format, args...)
}

// Fatal 致命错误日志
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.logInternal(FATAL, "FATAL", format, args...)
	os.Exit(1)
}

// GetLogFilePath 获取日志文件路径
func (l *Logger) GetLogFilePath() string {
	return l.filePath
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// Global logger instance
var globalLogger *Logger

// InitGlobalLogger 初始化全局日志记录器
func InitGlobalLogger(level LogLevel, logToFile bool) error {
	logger, err := NewLogger(level, logToFile)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// GetGlobalLogger 获取全局日志记录器
func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		// 默认使用控制台日志
		logger, _ := NewLogger(INFO, false)
		globalLogger = logger
	}
	return globalLogger
}

// Debug 全局调试日志
func Debug(format string, args ...interface{}) {
	GetGlobalLogger().Debug(format, args...)
}

// Info 全局信息日志
func Info(format string, args ...interface{}) {
	GetGlobalLogger().Info(format, args...)
}

// Warn 全局警告日志
func Warn(format string, args ...interface{}) {
	GetGlobalLogger().Warn(format, args...)
}

// Error 全局错误日志
func Error(format string, args ...interface{}) {
	GetGlobalLogger().Error(format, args...)
}

// Fatal 全局致命错误日志
func Fatal(format string, args ...interface{}) {
	GetGlobalLogger().Fatal(format, args...)
}
