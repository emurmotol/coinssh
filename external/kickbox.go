package external

import "strings"

const KickBoxApiUrl = "https://open.kickbox.com/v1/disposable/"

type KickBox struct {
	IsDisposable bool `json:"disposable"`
}

func IsEmailDisposable(email string) (bool, error) {
	kb := &KickBox{}
	url := strings.Join([]string{KickBoxApiUrl, email}, "")
	err := GetJson(url, kb)
	return kb.IsDisposable, err
}
