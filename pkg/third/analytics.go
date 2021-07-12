package third

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/51xpx/pkgx/pkg/jsonx"
	"github.com/51xpx/pkgx/pkg/threading"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

const (
	// AnalyticsEnvDefault default env
	AnalyticsEnvDefault = ""
	// AnalyticsEnvTest test env
	AnalyticsEnvTest = "test"
)

// Analytics bigdata log tracing service
type Analytics struct {
	Logger   *zap.Logger
	Domain   string
	Env      string
	Version  string
	AppID    int
	Platform int
}

// NewAnalytics new Analytics
func NewAnalytics(logger *zap.Logger, domain, env, version string, appID, platform int) *Analytics {
	return &Analytics{
		Logger:   logger,
		Domain:   domain,
		Env:      env,
		Version:  version,
		AppID:    appID,
		Platform: platform,
	}
}

// AnalyticsRequest request for bigdata logging
type AnalyticsRequest struct {
	domain  string
	data    map[string]interface{}
	isDebug bool
}

// SendDefault send to bigdata with default request
// uid: user id
// eid: event id
// seid: source event id
func (s *Analytics) SendDefault(uid string, eid string, seid string, extends map[string]interface{}) error {
	threading.GoSafe(func() {
		err := s.DefaultRequest(uid, eid, seid, extends).Send()
		if err != nil && s.Logger != nil {
			s.Logger.Error("SendDefault failed", zap.Error(err))
		}
	})
	return nil
}

// DefaultRequest default request for common backend event tracing
func (s *Analytics) DefaultRequest(uid string, eid string, seid string, extends map[string]interface{}) *AnalyticsRequest {
	now := time.Now()
	isDebug := false
	if s.Env == AnalyticsEnvTest {
		isDebug = true
	}
	request := AnalyticsRequest{
		domain:  s.Domain,
		data:    make(map[string]interface{}),
		isDebug: isDebug,
	}
	request.data["common"] = map[string]interface{}{
		"p":   s.Platform,
		"ver": s.Version,
	}
	event := map[string]interface{}{
		"et": 1,
		"payload": map[string]interface{}{
			"app":     cast.ToString(s.AppID),
			"env":     s.Env,
			"eid":     eid,
			"src_eid": seid,
			"ts":      cast.ToString(now.UnixNano() / 1e6),
			"uid":     uid,
		},
	}
	payload := event["payload"].(map[string]interface{})
	for k, v := range extends {
		payload[k] = v
	}
	request.data["events"] = []map[string]interface{}{
		event,
	}
	return &request
}

// Send send event to bigdata
func (r *AnalyticsRequest) Send() error {
	body, err := jsonx.Marshal(r.data)
	if err != nil {
		return err
	}

	url := r.domain + "/upload"
	if r.isDebug {
		url = url + "?debug=1"
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if r.isDebug {
		fmt.Println("post analytics request:", jsonx.MarshalToString(r.data))
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		bodyString := string(bodyBytes)
		fmt.Println("get analytics response:", bodyString)
	}
	return nil
}
