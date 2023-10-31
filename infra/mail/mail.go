package mail

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ChangSZ/blog/infra/conf"
	"github.com/ChangSZ/blog/infra/log"
	"gopkg.in/gomail.v2"
)

type EmailType string

type EmailParam struct {
	User        EmailType `json:"user"`
	Password    EmailType `json:"password"`
	Host        EmailType `json:"host"`
	Port        EmailType `json:"port"`
	To          EmailType `json:"to"`
	Subject     EmailType `json:"subject"`
	Body        EmailType `json:"body"`
	MailType    EmailType `json:"mail_type"`
	Description EmailType `json:"description"`
	Attaches    map[string]string
}

var mailParam *EmailParam

var mailAddr string

type EM func(*EmailParam) (interface{}, error)

func (et EmailType) CheckIsNull() error {
	if string(et) == "" {
		log.Error("value can not be null")
		return errors.New("value can not be null")
	}
	return nil
}

func (ep *EmailParam) SetMailUser(user EmailType) EM {
	return func(e *EmailParam) (interface{}, error) {
		u := e.User
		err := user.CheckIsNull()
		if err != nil {
			return nil, err
		}
		e.User = user
		return u, nil
	}
}

func (ep *EmailParam) SetMailPwd(pwd EmailType) EM {
	return func(ep *EmailParam) (interface{}, error) {
		p := ep.Password
		err := pwd.CheckIsNull()
		if err != nil {
			return nil, err
		}
		ep.Password = pwd
		return p, nil
	}
}

func (ep *EmailParam) SetMailHost(host EmailType) EM {
	return func(ep *EmailParam) (interface{}, error) {
		h := ep.Host
		err := host.CheckIsNull()
		if err != nil {
			return nil, err
		}
		ep.Host = host
		return h, nil
	}
}

func (ep *EmailParam) SetMailPort(port EmailType) EM {
	return func(ep *EmailParam) (interface{}, error) {
		h := ep.Port
		err := port.CheckIsNull()
		if err != nil {
			return nil, err
		}
		_, err = strconv.Atoi(string(port))
		if err != nil {
			return nil, fmt.Errorf("端口非数字：%v, err: %v", port, err)
		}
		ep.Port = port
		return h, nil
	}
}

func (ep *EmailParam) SetMailType(types EmailType) EM {
	return func(ep *EmailParam) (interface{}, error) {
		ty := ep.MailType
		err := types.CheckIsNull()
		if err != nil {
			return nil, err
		}
		ep.MailType = ty
		return ty, nil
	}
}

func (ep *EmailParam) MailInit(options ...EM) (*EmailParam, error) {
	q := &EmailParam{
		MailType: conf.MAIlTYPE,
	}
	for _, option := range options {
		_, err := option(q)
		if err != nil {
			return nil, err
		}
	}
	mailParam = q
	return q, nil
}

func (ep *EmailParam) SetSubject(s EmailType) *EmailParam {
	ep.Subject = s
	return ep
}

func (ep *EmailParam) SetDescription(de EmailType) *EmailParam {
	ep.Description = de
	return ep
}

func (ep *EmailParam) SetAttaches(a map[string]string) *EmailParam {
	ep.Attaches = a
	return ep
}

func (ep *EmailParam) SetBody(b EmailType) *EmailParam {
	ep.Body = b
	return ep
}

func (ep *EmailParam) SetTo(to EmailType) *EmailParam {
	ep.To = to
	return ep
}

func (ep *EmailParam) SendMail() error {
	subject := string(ep.Subject)
	user := string(ep.User)
	password := string(ep.Password)
	host := string(ep.Host)
	port, _ := strconv.Atoi(string(ep.Port)) // 前面有检测，此处无需处理error
	body := string(ep.Body)
	to := string(ep.To)

	m := gomail.NewMessage()
	// 发送人
	m.SetHeader("From", user)
	// 接收人
	m.SetHeader("To", to)
	// 抄送人
	// m.SetAddressHeader("Cc", "xxx@qq.com", "yyy")
	// 主题
	m.SetHeader("Subject", subject)
	// 内容
	m.SetBody("text/html", body)
	// 附件
	for _, attaFile := range ep.Attaches {
		m.Attach(attaFile)
	}

	// 拿到token，并进行连接,第4个参数是填授权码
	d := gomail.NewDialer(host, port, user, password)

	// 发送邮件
	err := d.DialAndSend(m)
	log.Info("mail, err:", err)
	return err
}

func SendMail(to string, subject string, body string) error {
	user := string(mailParam.User)
	password := string(mailParam.Password)
	host := string(mailParam.Host)
	port, _ := strconv.Atoi(string(mailParam.Port))

	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", "<html><body>"+body+"</body></html>")
	d := gomail.NewDialer(host, port, user, password)
	err := d.DialAndSend(m)
	log.Info("send mail, err:", err)
	return err
}
