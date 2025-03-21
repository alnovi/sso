package response

type Error struct {
	Code     int               `json:"-"`
	Error    string            `json:"error"`
	Validate map[string]string `json:"validate,omitempty"`
}
