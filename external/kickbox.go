package external

const KickBoxApiUrl = "https://open.kickbox.com/v1/disposable/"

type KickBox struct {
	IsDisposable bool `json:"disposable"`
}
