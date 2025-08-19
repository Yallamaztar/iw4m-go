package iw4m

import "net/http"

type IW4MWrapper struct {
	BaseURL  string
	ServerID string
	Cookie   string
	Client   *http.Client
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
