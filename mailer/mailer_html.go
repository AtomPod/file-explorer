package mailer

import (
	"bytes"
	"html/template"
	"sync"

	"github.com/phantom-atom/file-explorer/internal/log"
)

//MailTemplate 邮件模板
type MailTemplate struct {
	Subject     string
	ContentType string
	tpl         *template.Template
}

//Body 获取body
func (m *MailTemplate) Body(data interface{}) (string, error) {
	buf := buffPool.Get().(*bytes.Buffer)
	buf.Reset()

	if err := m.tpl.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

var (
	tplMap    = make(map[string]*MailTemplate)
	tplMapMux = &sync.RWMutex{}
	buffPool  = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
)

//RegisterMailTemplate 注册邮件模板
func RegisterMailTemplate(
	name string,
	filename string,
	subject, contentType string) error {

	if name == "" {
		panic("mail-template: name is empty")
	}

	tplMapMux.Lock()
	if _, ok := tplMap[name]; ok {
		tplMapMux.Unlock()
		panic("mail-template: register called twice for name " + name)
	}

	t, err := template.ParseFiles(filename)
	if err != nil {
		tplMapMux.Unlock()
		return err
	}

	tplMap[name] = &MailTemplate{
		Subject:     subject,
		ContentType: contentType,
		tpl:         t,
	}
	tplMapMux.Unlock()

	log.Info("msg", "load mail template file", "name", name, "file", filename)
	return nil
}

//GetMailTemplate 获取邮件模板
func GetMailTemplate(name string) *MailTemplate {
	tplMapMux.RLock()
	tpl, ok := tplMap[name]
	tplMapMux.RUnlock()

	if !ok {
		return nil
	}
	return tpl
}
