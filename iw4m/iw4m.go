package iw4m

import (
	"fmt"
	"net/http"
)

type IW4MWrapper struct {
	BaseURL  string
	ServerID string
	Cookie   string
	Client   *http.Client
}

func (iw4m *IW4MWrapper) DoRequest(endpoint string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", iw4m.BaseURL, endpoint)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", iw4m.Cookie)
	res, err := iw4m.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Create a new instance of the iw4m wrapper
func NewWrapper(baseUrl string, serverID string, cookie string) *IW4MWrapper {
	return &IW4MWrapper{
		BaseURL:  baseUrl,
		ServerID: serverID,
		Cookie:   cookie,
		Client:   &http.Client{},
	}
}
