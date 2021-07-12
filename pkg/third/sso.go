package third

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Sso sso service
type Sso struct {
	Domain string
}

// SendWorkWechatMsg send message to work wechat
func (s *Sso) SendWorkWechatMsg(users []string, msg string) error {
	if len(users) == 0 || len(msg) == 0 {
		return errors.New("Invalid parameters")
	}

	http.DefaultClient.Timeout = time.Second * 5
	resp, err := http.PostForm(s.Domain+"/api/alarm/send", url.Values{
		"users":   {strings.Join(users, ",")},
		"content": {msg},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
