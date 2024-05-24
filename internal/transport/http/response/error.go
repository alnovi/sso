package response

type Error struct {
	Code     int               `json:"-"`
	Message  string            `json:"message"`
	Validate map[string]string `json:"validate,omitempty"`
}
