package log

import (
	"github.com/phantom-atom/file-explorer/config"
	"go.uber.org/zap"
)

type zapLogger struct {
	logger *zap.SugaredLogger
}

//NewZapLogger 新建标准日志
func NewZapLogger(config *config.Config) Interface {
	var logger *zap.Logger
	var err error

	var conf *zap.Config
	if config.Mode == "release" {
		production := zap.NewProductionConfig()
		conf = &production
	} else {
		development := zap.NewDevelopmentConfig()
		conf = &development
	}

	logConf := &config.Log
	if logConf.Level != "" {
		level := zap.NewAtomicLevel()
		if err := level.UnmarshalText([]byte(logConf.Level)); err != nil {
			panic(err)
		} else {
			conf.Level = level
		}
	}

	if len(logConf.OutputPaths) != 0 {
		conf.OutputPaths = logConf.OutputPaths
	}

	if len(logConf.ErrorOutputPaths) != 0 {
		conf.ErrorOutputPaths = logConf.ErrorOutputPaths
	}

	if logConf.Encoding != "" {
		conf.Encoding = logConf.Encoding
	}

	logger, err = conf.Build()

	if err != nil {
		panic(err)
	}
	return &zapLogger{
		logger: logger.Sugar(),
	}
}

//Debug 使用全局日志进行打印Debug信息
func (s *zapLogger) Debug(v ...interface{}) {
	s.logger.Debugw("", v...)
}

//Info 使用全局日志进行打印Info信息
func (s *zapLogger) Info(v ...interface{}) {
	s.logger.Infow("", v...)
}

//Warn 使用全局日志进行打印Warn信息
func (s *zapLogger) Warn(v ...interface{}) {
	s.logger.Warnw("", v...)
}

//Error 使用全局日志进行打印Error信息
func (s *zapLogger) Error(v ...interface{}) {
	s.logger.Errorw("", v...)
}

//Panic 使用全局日志进行打印Panic信息
func (s *zapLogger) Panic(v ...interface{}) {
	s.logger.Panicw("", v...)
}
