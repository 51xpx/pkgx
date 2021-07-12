package zapx

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/51xpx/pkgx/pkg/threading"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// BaseTimeFormatter base time formatter string
	BaseTimeFormatter = "2006-01-02 15:04:05.000000"
)

// LogConfig zap log config
type LogConfig struct {
	Level       int
	OutputPaths []string
	Dev         bool
}

// RotateLogConfig rotate log config
type RotateLogConfig struct {
	Level      int
	Filename   string
	MaxSize    int // megabytes
	MaxBackups int
	MaxAge     int // days
}

// K8sLogConfig k8s log config
type K8sLogConfig struct {
	Level      int
	Dev        bool
	AccessFile string
	ErrorFile  string
}

var loggers []*zap.Logger

// NewZap new zap logger
func NewZap(cfg *LogConfig) (log *zap.Logger, err error) {
	var (
		zcfg zap.Config
	)
	zcfg = zap.NewProductionConfig()
	zcfg.OutputPaths = cfg.OutputPaths
	zcfg.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.Level))
	zcfg.Development = cfg.Dev
	zcfg.EncoderConfig.EncodeTime = TimeEncoder
	log, err = zcfg.Build()
	if err != nil {
		return
	}
	newCore := zapcore.NewTee(
		log.Core(),
	)

	log = zap.New(newCore).WithOptions(zap.AddCaller())
	loggers = append(loggers, log)
	return
}

// NewRotateZap new zap logger with rotatable
func NewRotateZap(cfg *RotateLogConfig) (log *zap.Logger, err error) {
	level := zapcore.Level(cfg.Level)
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeTime = TimeEncoder
	je := zapcore.NewJSONEncoder(ec)

	ws := zapcore.Lock(os.Stdout)
	stdCore := zapcore.NewCore(je, ws, level)

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize, // megabytes
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge, // days
		LocalTime:  true,
	})
	rotateCore := zapcore.NewCore(je, w, level)

	newCore := zapcore.NewTee(
		stdCore,
		rotateCore,
	)
	log = zap.New(newCore).WithOptions(zap.AddCaller())
	loggers = append(loggers, log)
	return
}

// NewK8sZap new zap logger with k8s
func NewK8sZap(cfg *K8sLogConfig) (log *zap.Logger, err error) {
	var (
		zcfg zap.Config
	)
	zcfg = zap.NewProductionConfig()
	zcfg.OutputPaths = []string{cfg.AccessFile}
	zcfg.ErrorOutputPaths = []string{cfg.ErrorFile}
	zcfg.Level = zap.NewAtomicLevelAt(zapcore.Level(cfg.Level))
	zcfg.Development = cfg.Dev
	zcfg.EncoderConfig.EncodeTime = TimeEncoder
	log, err = zcfg.Build()
	if err != nil {
		return
	}
	newCore := zapcore.NewTee(
		log.Core(),
	)

	log = zap.New(newCore).WithOptions(zap.AddCaller())
	loggers = append(loggers, log)
	return
}

// TimeEncoder time encoder config
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(BaseTimeFormatter))
}

// Sync sync logger buffer
func Sync() {
	sync()
}

// Async sync logger buffer on background with 1s interval
func Async() {
	threading.GoSafe(func() {
		for {
			sync()
			time.Sleep(time.Second)
		}
	})
}

func sync() {
	for _, logger := range loggers {
		if logger == nil {
			continue
		}
		logger.Sync()
	}
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}
