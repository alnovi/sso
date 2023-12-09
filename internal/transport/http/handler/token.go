package handler

import "github.com/labstack/echo/v4"

type TokenHandler struct {
}

func NewTokenHandler() *TokenHandler {
	return &TokenHandler{}
}

func (h *TokenHandler) GenerateToken(c echo.Context) error {
	//TODO: обработать query параметры
	//TODO: добавить use-case

	return nil
}
