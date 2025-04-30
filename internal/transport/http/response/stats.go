package response

type Stats struct {
	Users    int `json:"users"`
	Clients  int `json:"clients"`
	Sessions int `json:"sessions"`
}
