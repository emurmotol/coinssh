package external

import (
	"time"
	"github.com/gobuffalo/envy"
	"net"
	"net/url"
	"net/http"
)

const (
	ReCaptchaApiUrl = "https://www.google.com/recaptcha/api/siteverify"
	ReCaptchaPostParam = "g-recaptcha-response"
)

type ReCaptcha struct {
	Success     bool      `json:"success"`
	ChallengeTs time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func IsHuman(r *http.Request) (bool, error) {
	response := r.FormValue(ReCaptchaPostParam)
	remoteIp, _, _ := net.SplitHostPort(r.RemoteAddr)
	secret := envy.Get("RECAPTCHA_SECRET_KEY", "")

	data := url.Values{
		"secret": {secret},
		"response": {response},
		"remoteip": {remoteIp},
	}
	rc := &ReCaptcha{}
	err := PostJson(ReCaptchaApiUrl, data, rc)

	if err != nil {
		return false, err
	}
	return rc.Success, nil
}