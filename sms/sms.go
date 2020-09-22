package sms

import (
	"errors"
	"regexp"
)

type isms interface {
	Send(mobile, action string) (string, error)
	Verify(mobile, action, code string) error
}

type gooSms struct {
	sms isms
}

func New(sms isms) *gooSms {
	return &gooSms{
		sms: sms,
	}
}

func (s *gooSms) Send(mobile, action string) (string, error) {
	if mobile == "" {
		return "", errors.New("mobile is null")
	}
	if regexp.MustCompile(`^1[3,4,5,7,8]\d{9}$`).MatchString(mobile) == false {
		return "", errors.New("invalid mobile")
	}
	if action == "" {
		return "", errors.New("action is null")
	}
	return s.sms.Send(mobile, action)
}

func (s *gooSms) Verify(mobile, action, code string) error {
	if mobile == "" {
		return errors.New("mobile is null")
	}
	if regexp.MustCompile(`^1[3,4,5,7,8]\d{9}$`).MatchString(mobile) == false {
		return errors.New("invalid mobile")
	}
	if action == "" {
		return errors.New("action is null")
	}
	if code == "" {
		return errors.New("code is null")
	}
	return s.sms.Verify(mobile, action, code)
}
