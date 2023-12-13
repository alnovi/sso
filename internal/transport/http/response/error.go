package response

type Error struct {
	Message string `json:"message"`
}

type ErrorValidate struct {
	Message  string            `json:"message"`
	Validate map[string]string `json:"validate"`
}
