package logging

import (
	"context"
	"e-klinik/config"
	"fmt"
	"time"

	zaploki "github.com/th1cha/zap-loki"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	lokiAddress = "http://localhost:3100"
	appName     = "pvsavechan"
)

var zapSugarLogger *zap.SugaredLogger

type lokiLogger struct {
	cfg    *config.Config
	logger *zap.SugaredLogger
}

// var zapLogLevelMapping = map[string]zapcore.Level{
// 	"debug": zapcore.DebugLevel,
// 	"info":  zapcore.InfoLevel,
// 	"warn":  zapcore.WarnLevel,
// 	"error": zapcore.ErrorLevel,
// 	"fatal": zapcore.FatalLevel,
// }

func newZapLokiLogger(cfg *config.Config) *lokiLogger {

	logger := &lokiLogger{cfg: cfg}
	logger.Init()
	return logger
}

func (l *lokiLogger) getLogLevel() zapcore.Level {
	level, exists := zapLogLevelMapping[l.cfg.Logger.Level]
	if !exists {
		return zapcore.DebugLevel
	}
	return level
}

func (l *lokiLogger) Init() {
	// once.Do(func() {

	// fileName := fmt.Sprintf("%s%s-%s.%s", l.cfg.Logger.FilePath, time.Now().Format("2006-01-02"), uuid.New(), "log")
	// w := zapcore.AddSync(&lumberjack.Logger{
	// 	Filename:   fileName,
	// 	MaxSize:    1,
	// 	MaxAge:     20,
	// 	LocalTime:  true,
	// 	MaxBackups: 5,
	// 	Compress:   true,
	// })

	// config := zap.NewProductionEncoderConfig()
	// config.EncodeTime = zapcore.ISO8601TimeEncoder

	// core := zapcore.NewCore(
	// 	zapcore.NewJSONEncoder(config),
	// 	w,
	// 	l.getLogLevel(),
	// )

	// logger := zap.New(core, zap.AddCaller(),
	// 	zap.AddCallerSkip(1),
	// 	zap.AddStacktrace(zapcore.ErrorLevel),
	// ).Sugar()
	zapConfig := zap.NewProductionConfig()
	loki := zaploki.New(context.Background(), zaploki.Config{
		Url:          lokiAddress,
		BatchMaxSize: 1000,
		BatchMaxWait: 10 * time.Second,
		Labels:       map[string]string{"app": appName},
	})

	logger, err := loki.WithCreateLogger(zapConfig)
	// logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		fmt.Println(err)
	}
	// zapSugarLogger = logger.Sugar()
	// })

	l.logger = logger.Sugar()
}

func (l *lokiLogger) Debug(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogInfo(cat, sub, extra)

	l.logger.Debugw(msg, params...)
}

func (l *lokiLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debug(template)
}

func (l *lokiLogger) Info(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLokiLogInfo(cat, sub, extra)
	l.logger.Infow(msg, params...)
}

func (l *lokiLogger) Infof(template string, args ...interface{}) {
	l.logger.Info(template, args)
}

func (l *lokiLogger) Warn(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLokiLogInfo(cat, sub, extra)
	l.logger.Warnw(msg, params...)
}

func (l *lokiLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args)
}

func (l *lokiLogger) Error(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLokiLogInfo(cat, sub, extra)
	l.logger.Errorw(msg, params...)
}

func (l *lokiLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args)
}

func (l *lokiLogger) Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLokiLogInfo(cat, sub, extra)
	l.logger.Fatalw(msg, params...)
}

func (l *lokiLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args)
}

func prepareLokiLogInfo(cat Category, sub SubCategory, extra map[ExtraKey]interface{}) []interface{} {
	if extra == nil {
		extra = make(map[ExtraKey]interface{})
	}
	extra["Category"] = cat
	extra["SubCategory"] = sub

	return logParamsToZapParams(extra)
}
