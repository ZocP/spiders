package queryWithOption

const (
	REGEX = "regex"
	FUZZY = "fuzzy"
)

type RequestBody struct {
	Keyword string `json:"keyword"`
	Option  string `json:"option"`
}
