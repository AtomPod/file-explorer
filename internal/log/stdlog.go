package log

import (
	"log"
)

type stdLogger struct{}

//NewStdLogger 新建标准日志
func NewStdLogger() Interface {
	return &stdLogger{}
}

//Debug 使用全局日志进行打印Debug信息
func (s *stdLogger) Debug(v ...interface{}) {
	log.Println(s.appendLevel("[Debug]", v...)...)
}

//Info 使用全局日志进行打印Info信息
func (s *stdLogger) Info(v ...interface{}) {
	log.Println(s.appendLevel("[Info]", v...)...)
}

//Warn 使用全局日志进行打印Warn信息
func (s *stdLogger) Warn(v ...interface{}) {
	log.Println(s.appendLevel("[Warn]", v...)...)
}

//Error 使用全局日志进行打印Error信息
func (s *stdLogger) Error(v ...interface{}) {
	log.Println(s.appendLevel("[Error]", v...)...)
}

//Panic 使用全局日志进行打印Panic信息
func (s *stdLogger) Panic(v ...interface{}) {
	log.Println(s.appendLevel("[Panic]", v...)...)
}

func (s *stdLogger) appendLevel(l string, v ...interface{}) []interface{} {
	withLevel := make([]interface{}, 0, len(v)+1)
	withLevel = append(withLevel, l)
	withLevel = append(withLevel, s.decorate(v...))
	return withLevel
}

func (s *stdLogger) decorate(v ...interface{}) interface{} {
	kvpair := make(map[interface{}]interface{})

	l := len(v)
	pairs := l / 2

	for p := 0; p < pairs; p++ {
		kvpair[v[p*2]] = v[p*2+1]
	}

	if l%2 != 0 {
		kvpair[v[l-1]] = "non-value"
	}

	// jbyt, err := json.Marshal(kvpair)
	// if err != nil {
	// 	kvpair["json_marshal_err"] = err
	// } else {
	// 	return jbyt
	// }
	return kvpair
}
