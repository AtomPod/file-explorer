package log

import "sync/atomic"

//Interface 日志接口
type Interface interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Panic(v ...interface{})
}

var logger atomic.Value

//SetLogger 设置当前日志的实现作为全局使用
func SetLogger(log Interface) {
	logger.Store(&log)
}

//Logger 获取当前全局日志，如果未使用SetLogger设置，则返回nil
func Logger() Interface {
	impl := logger.Load()
	if impl == nil {
		return nil
	}
	return *impl.(*Interface)
}

//Debug 使用全局日志进行打印Debug信息
func Debug(v ...interface{}) {
	if log := Logger(); log != nil {
		log.Debug(v...)
	}
}

//Info 使用全局日志进行打印Info信息
func Info(v ...interface{}) {
	if log := Logger(); log != nil {
		log.Info(v...)
	}
}

//Warn 使用全局日志进行打印Warn信息
func Warn(v ...interface{}) {
	if log := Logger(); log != nil {
		log.Warn(v...)
	}
}

//Error 使用全局日志进行打印Error信息
func Error(v ...interface{}) {
	if log := Logger(); log != nil {
		log.Error(v...)
	}
}

//Panic 使用全局日志进行打印Panic信息
func Panic(v ...interface{}) {
	if log := Logger(); log != nil {
		log.Panic(v...)
	}
}
