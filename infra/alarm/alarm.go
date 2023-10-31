package alarm

import (
	"regexp"
	"strings"

	"github.com/ChangSZ/blog/infra/log"
	"github.com/ChangSZ/blog/infra/mail"
	"github.com/go-errors/errors"
)

// Define AlarmType to string
// for to check the params is right
type AlarmType string

type AlarmMailReceive string

// this are some const params what i defined
// only this can be to input
const (
	AlarmTypeMail    AlarmType = "mail"
	AlarmTypeWechat  AlarmType = "wechat"
	AlarmTypeMessage AlarmType = "message"
)

type AlarmParam struct {
	Types  AlarmType
	MailTo AlarmMailReceive
}

var alarmParam *AlarmParam

// Define a closure type to next
type ap func(*AlarmParam) (interface{}, error)

// can use this function to set a new value
// but to check it is a right type
func (alarm *AlarmParam) SetType(t AlarmType) ap {
	return func(alarm *AlarmParam) (interface{}, error) {
		str := strings.Split(string(t), ",")
		if len(str) == 0 {
			log.Error("you must input a value")
			return nil, errors.New("you must input a value")
		}
		for _, types := range str {
			s := AlarmType(types)
			_, err := s.IsCurrentType()
			if err != nil {
				log.Error("the value type is error")
				return nil, err
			}
		}
		ty := alarm.Types
		alarm.Types = t
		return ty, nil
	}
}

func (alarm *AlarmParam) SetMailTo(t AlarmMailReceive) ap {
	return func(alarm *AlarmParam) (interface{}, error) {
		to := alarm.MailTo
		_, err := t.CheckIsNull()
		if err != nil {
			return nil, err
		}
		_, err = t.MustMailFormat()
		if err != nil {
			return nil, err
		}
		alarm.MailTo = t
		return to, nil
	}
}

// alarm receive account can not null
func (t AlarmMailReceive) CheckIsNull() (AlarmMailReceive, error) {
	if len(t) == 0 {
		log.Error("value can not be null")
		return "", errors.New("value can not be null")
	}
	return t, nil
}

// alarm receive account must be mail format
func (t AlarmMailReceive) MustMailFormat() (AlarmMailReceive, error) {
	if m, _ := regexp.MatchString("^([a-zA-Z0-9_-])+@([a-zA-Z0-9_-])+(.[a-zA-Z0-9_-])+", string(t)); !m {
		log.Error("value format is not right")
		return "", errors.New("value format is not right")
	}
	return t, nil
}

// judge it is a right type what i need
// if is it a wrong type, i must return a panic to above
func (at AlarmType) IsCurrentType() (AlarmType, error) {
	switch at {
	case AlarmTypeMail:
		return at, nil
	case AlarmTypeWechat:
		return at, nil
	case AlarmTypeMessage:
		return at, nil
	default:
		log.Error("the alarm type is error")
		return at, errors.New("the alarm type is error")
	}
}

// implementation value
func (alarm *AlarmParam) AlarmInit(options ...ap) error {
	q := &AlarmParam{}
	for _, option := range options {
		_, err := option(q)
		if err != nil {
			return err
		}
	}
	alarmParam = q
	return nil
}

func Alarm(content string) {
	types := strings.Split(string(alarmParam.Types), ",")
	var err error
	for _, a := range types {
		switch AlarmType(a) {
		case AlarmTypeMail:
			if alarmParam.MailTo == "" {
				log.Error("邮件接收者不能为空")
				break
			}
			err = mail.SendMail(string(alarmParam.MailTo), "报警", content)
			break
		case AlarmTypeWechat:
			break
		case AlarmTypeMessage:
			break
		}
		if err != nil {
			log.Errorf("alarm is error: %v", err)
		}
	}
}
