package request

type CreateClient struct {
	Id       string  `json:"id" validate:"required,min=3,max=30,client_id,lowercase"`
	Name     string  `json:"name" validate:"required,min=5,max=50"`
	Icon     *string `json:"icon" validate:"omitnil,uri,max=250"`
	Callback string  `json:"callback" validate:"required,url,max=250"`
	Secret   *string `json:"secret" validate:"omitnil,min=5,max=100"`
}

type UpdateClient struct {
	Name     string  `json:"name" validate:"required,min=5,max=50"`
	Icon     *string `json:"icon" validate:"omitnil,uri,max=250"`
	Callback string  `json:"callback" validate:"required,uri,max=250"`
	Secret   string  `json:"secret" validate:"required,min=5,max=100"`
}
