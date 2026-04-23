package log

import (
	"log/slog"
	"testing"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestMixtureLog(t *testing.T) {
	devZapConfig := DevConfig()
	prodZapConfig := ProdConfig()
	prodZapConfig.LevelWriters = map[zapcore.Level]*lumberjack.Logger{
		zapcore.DebugLevel: DefaultLumberjack("", "debug.log"),
		zapcore.InfoLevel:  DefaultLumberjack("", "info.log"),
		zapcore.WarnLevel:  DefaultLumberjack("", "warn.log"),
		zapcore.ErrorLevel: DefaultLumberjack("", "error.log"),
	}
	prodZapConfig.IsJSON = true
	devZapConfig.ProjectName = "111"
	devZapConfig.FieldFunc = func() []zapcore.Field {
		return []zapcore.Field{
			zap.String("name", "test"),
			zap.Int("age", 10),
		}
	}
	Init(devZapConfig, prodZapConfig)
	zap.S().Debugf("hello world")
	zap.S().Warnf("hello world")
	slog.Debug("hello world")
	slog.Warn("hello world")
	go func() {
		zap.S().Debugf("hello world")
		zap.S().Infof("hello world")
		zap.S().Warnf("hello world")
		slog.Debug("hello world")
		slog.Warn("hello world")
	}()
	for {
	}
}

func TestDevLog(t *testing.T) {
	InitZapGlobal(DevConfig())
	zap.S().Debugf("hello world")
	zap.S().Infof("hello world")
	zap.S().Warnf("hello world")
	go func() {
		zap.S().Debugf("hello world")
		zap.S().Infof("hello world")
		zap.S().Warnf("hello world")
	}()
	for {
	}
}

func TestSlog(t *testing.T) {
	devZapConfig := DevConfig()
	prodZapConfig := ProdConfig()
	prodZapConfig.IsJSON = true
	InitZapGlobal(devZapConfig, prodZapConfig)
	zap.S().Infof("hello world")
	InitSlogGlobal(devZapConfig, prodZapConfig)
	slog.Debug("hello world")
	slog.Info("Hello World")
	slog.Error("123")
}

func TestProdLog(t *testing.T) {
	InitZapGlobal(ProdConfig())
	zap.S().Debugf("hello world")
	zap.S().Infof("hello world")
	zap.S().Warnf("hello world")
	zap.S().Errorf("hello world")
}

func TestLevelConfig(t *testing.T) {
	prodZapConfig := ProdConfig()
	prodZapConfig.LevelWriters = map[zapcore.Level]*lumberjack.Logger{
		zapcore.DebugLevel: DefaultLumberjack("", "debug.log"),
		zapcore.InfoLevel:  DefaultLumberjack("", "info.log"),
		zapcore.WarnLevel:  DefaultLumberjack("", "warn.log"),
		zapcore.ErrorLevel: DefaultLumberjack("", "error.log"),
	}
	InitZapGlobal(DevConfig(), prodZapConfig)
	zap.S().Debugf("hello world")
	zap.S().Infof("hello world")
	zap.S().Warnf("hello world")
	zap.S().Errorf("hello world")
}

func TestMultiLog(t *testing.T) {
	InitZapGlobal(DevConfig(), ProdConfig())
	zap.S().Debugf("hello world")
	zap.S().Infof("hello world")
	zap.S().Warnf("hello world")
	zap.S().Errorf("hello world")
}
