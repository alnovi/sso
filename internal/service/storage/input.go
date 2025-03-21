package storage

type InputClientCreate struct {
	Id       string
	Name     string
	Icon     *string
	Callback string
	Secret   *string
}

type InputClientUpdate struct {
	Id       string
	Name     string
	Icon     *string
	Callback string
	Secret   string
}

type InputUserCreate struct {
	Name     string
	Email    string
	Password string
}

type InputUserUpdate struct {
	Id       string
	Name     string
	Email    string
	Password *string
}
