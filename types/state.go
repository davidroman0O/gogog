package types

import "net/http"

const (
	Hostname = "https://www.gog.com"
)

type Gogog struct{}

type GogAuth struct {
	Cookies []*Cookie
	User    *UserData
}

type GogStates interface {
	GogAuth | Gogog
}

func NewTransporter(user UserAgent) *transporter {
	return &transporter{
		user: user,
	}
}

type transporter struct {
	user UserAgent
}

func (t *transporter) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Origin", Hostname)
	req.Header.Add("Referer", Hostname+"/")
	req.Header.Add("User-Agent", t.user.String())
	return http.DefaultTransport.RoundTrip(req)
}

type UserAgent string

func (u *UserAgent) String() string {
	return string(*u)
}
