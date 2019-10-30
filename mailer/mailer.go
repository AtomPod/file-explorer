package mailer

import (
	"context"
	"crypto/tls"
	"errors"
	"time"

	"github.com/phantom-atom/file-explorer/internal/log"

	"gopkg.in/gomail.v2"

	"github.com/phantom-atom/file-explorer/config"
)

var (
	//ErrMailerIsClosed closed
	ErrMailerIsClosed = errors.New("mailer is closed")
)

//Mailer 邮件发送器
type Mailer struct {
	ctx       context.Context
	cancel    context.CancelFunc
	config    func() *config.Config
	messageCh chan *gomail.Message
}

//NewMailer 新建Mailer
func NewMailer(ctx context.Context,
	configFunc func() *config.Config,
	maxQueueSize int64) *Mailer {

	if ctx == nil {
		ctx = context.Background()
	}

	canceledCtx, cancel := context.WithCancel(ctx)

	mailer := &Mailer{
		ctx:       canceledCtx,
		cancel:    cancel,
		config:    configFunc,
		messageCh: make(chan *gomail.Message, maxQueueSize),
	}

	go mailer.send()
	return mailer
}

//Send 发送邮件
func (r *Mailer) Send(from string, to []string,
	subject string,
	contentType string,
	body string) error {

	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", to...)
	message.SetHeader("Subject", subject)
	message.SetBody(contentType, body)

	select {
	case <-r.ctx.Done():
		return ErrMailerIsClosed
	case r.messageCh <- message:
	}
	return nil
}

//Close Close接口
func (r *Mailer) Close() error {
	select {
	case <-r.ctx.Done():
		return ErrMailerIsClosed
	default:
	}
	r.cancel()
	return nil
}

func (r *Mailer) getDialer() (*gomail.Dialer, error) {
	emailConf := &r.config().Email
	dialer := &gomail.Dialer{
		Host:      emailConf.Host,
		Port:      emailConf.Port,
		Username:  emailConf.Username,
		Password:  emailConf.Password,
		SSL:       emailConf.SSL,
		LocalName: emailConf.LocalName,
	}

	if emailConf.SSL {
		cert, err := tls.LoadX509KeyPair(emailConf.CertPath, emailConf.KeyPath)
		if err != nil {
			return nil, err
		}
		dialer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{
				cert,
			},
		}
	}

	return dialer, nil
}

func (r *Mailer) send() {
	var sender gomail.SendCloser

	for {
		select {
		case <-r.ctx.Done():
			return
		case msg := <-r.messageCh:
			for sender == nil {
				select {
				case <-r.ctx.Done():
					return
				default:
				}

				dialer, err := r.getDialer()

				var newSender gomail.SendCloser
				if err == nil {
					newSender, err = dialer.Dial()
				}

				if err != nil {
					log.Error("msg", "occur an error when send email", "error", err.Error())
					time.Sleep(5 * time.Second)
					continue
				}
				sender = newSender
			}

			if err := gomail.Send(sender, msg); err != nil {
				log.Error("msg", "occur an error when send email", "error", err.Error())
			}
		case <-time.After(30 * time.Second):
			if sender != nil {
				if err := sender.Close(); err != nil {
					log.Error("msg", "occur an error when close email sender", "error", err.Error())
				}
				sender = nil
			}
		}
	}
}
