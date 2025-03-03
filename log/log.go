package log

import (
	"fmt"
	"github.com/Csy12139/Vesper/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"path/filepath"
	"runtime"
)

var zapLogger *zap.SugaredLogger
var atomicLevel zap.AtomicLevel

// SetLogLevel 用于在运行时修改日志等级
func SetLogLevel(level string) error {
	return atomicLevel.UnmarshalText([]byte(level))
}

// InitLog Levels ("debug", "info", "warn", "error", "dpanic", "panic", and "fatal") 大小写随意
func InitLog(logDir string, maxFileSizeMb int, maxFileNum int, logLevel string) error {
	execName, err := common.GetExecName()
	if err != nil {
		return fmt.Errorf("get exe name error [%s]", err)
	}

	logFile := filepath.Join(logDir, fmt.Sprintf("%s.log", execName))

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    maxFileSizeMb,
		MaxBackups: maxFileNum,
		Compress:   false,
		LocalTime:  true,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000000")

	encoderConfig.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		pc, file, line, ok := runtime.Caller(7)
		caller = zapcore.NewEntryCaller(pc, file, line, ok)
		if !ok {
			enc.AppendString("unknown")
			//enc.AppendString("unknown")
			return
		}
		//fn := runtime.FuncForPC(pc)
		//if fn == nil {
		//	enc.AppendString("unknown")
		//	enc.AppendString("unknown")
		//	return
		//}
		enc.AppendString(caller.TrimmedPath())
		//enc.AppendString(fn.Name())
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	atomicLevel, err = zap.ParseAtomicLevel(logLevel)
	if err != nil {
		return err
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		atomicLevel,
	)
	zapLogger = zap.New(core, zap.AddCaller()).Sugar()
	return nil
}

func Info(args ...interface{}) {
	zapLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	zapLogger.Infof(template, args...)
}

func Debug(args ...interface{}) {
	zapLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	zapLogger.Debugf(template, args...)
}

func Warn(args ...interface{}) {
	zapLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	zapLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	zapLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	zapLogger.Errorf(template, args...)
}

func Fatal(args ...interface{}) {
	zapLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	zapLogger.Fatalf(template, args...)
}
