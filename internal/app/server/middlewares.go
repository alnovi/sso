package server

type Middlewares struct {
}

func newMiddlewares(_ *App, services *Services) (*Middlewares, error) {
	return &Middlewares{}, nil
}
