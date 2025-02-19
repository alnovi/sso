// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Alnovi",
            "url": "https://github.com/alnovi/sso"
        },
        "license": {
            "name": "MIT",
            "url": "https://github.com/alnovi/sso/LICENSE.md"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/oauth/v1/authorize": {
            "post": {
                "description": "Аутентификация пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuth"
                ],
                "summary": "Аутентификация пользователя",
                "operationId": "Authorize",
                "parameters": [
                    {
                        "type": "string",
                        "example": "code",
                        "description": "Response type",
                        "name": "response_type",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "app_id",
                        "description": "Client ID",
                        "name": "client_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Redirect URI",
                        "name": "redirect_uri",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "State",
                        "name": "state",
                        "in": "query"
                    },
                    {
                        "description": "Логин и пароль пользователя",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.Authorize"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ссылка для перехода",
                        "schema": {
                            "$ref": "#/definitions/response.URL"
                        }
                    },
                    "302": {
                        "description": "Found"
                    }
                }
            }
        },
        "/oauth/v1/client": {
            "get": {
                "description": "Проверка параметров клиента и информация о нем",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuth"
                ],
                "summary": "Информация о клиенте",
                "operationId": "CheckClient",
                "parameters": [
                    {
                        "type": "string",
                        "example": "code",
                        "description": "Response type",
                        "name": "response_type",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "app_id",
                        "description": "Client ID",
                        "name": "client_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Redirect URI",
                        "name": "redirect_uri",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/oauth/v1/profile": {
            "get": {
                "description": "Профиль пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuth profile"
                ],
                "summary": "Профиль пользователя",
                "operationId": "Profile",
                "responses": {
                    "200": {
                        "description": "Профиль пользователя",
                        "schema": {
                            "$ref": "#/definitions/response.Profile"
                        }
                    }
                }
            }
        },
        "/oauth/v1/profile/logout": {
            "post": {
                "description": "Удаление сессии пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuth profile"
                ],
                "summary": "Выход",
                "operationId": "Logout",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/v1/oauth/forgot-password": {
            "post": {
                "description": "Отправка ссылки для смены пароля",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuth-password"
                ],
                "summary": "Забыли пароль",
                "operationId": "ForgotPassword",
                "parameters": [
                    {
                        "type": "string",
                        "example": "code",
                        "description": "Response type",
                        "name": "response_type",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "app_id",
                        "description": "Client ID",
                        "name": "client_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Redirect URI",
                        "name": "redirect_uri",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "State",
                        "name": "state",
                        "in": "query"
                    },
                    {
                        "description": "Логин пользователя",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.ForgotPassword"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Сообщение для пользователя",
                        "schema": {
                            "$ref": "#/definitions/response.Message"
                        }
                    }
                }
            }
        },
        "/v1/oauth/reset-password": {
            "post": {
                "description": "Изменение пароля пользователя",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuth-password"
                ],
                "summary": "Смена пароля",
                "operationId": "ResetPassword",
                "parameters": [
                    {
                        "type": "string",
                        "example": "secret",
                        "description": "Разовый токен",
                        "name": "hash",
                        "in": "query",
                        "required": true
                    },
                    {
                        "description": "Новый пароль пользователя",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.ResetPassword"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ссылка для перехода",
                        "schema": {
                            "$ref": "#/definitions/response.URL"
                        }
                    }
                }
            }
        },
        "/v1/oauth/token": {
            "post": {
                "description": "Получение токена доступа по code или refresh токену",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "OAuth"
                ],
                "summary": "Токен доступа",
                "operationId": "Token",
                "parameters": [
                    {
                        "enum": [
                            "authorization_code",
                            "refresh_token"
                        ],
                        "type": "string",
                        "description": "Grant type",
                        "name": "grant_type",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "app_id",
                        "description": "Client ID",
                        "name": "client_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "secret",
                        "description": "Client secret",
                        "name": "client_secret",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "secret",
                        "description": "Code token",
                        "name": "code",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "secret",
                        "description": "Refresh token",
                        "name": "refresh_token",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Токен доступа",
                        "schema": {
                            "$ref": "#/definitions/response.AccessToken"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "request.Authorize": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string",
                    "minLength": 5,
                    "example": "name@example.com"
                },
                "password": {
                    "type": "string",
                    "maxLength": 24,
                    "minLength": 5,
                    "example": "qwerty"
                },
                "remember": {
                    "type": "boolean"
                }
            }
        },
        "request.ForgotPassword": {
            "type": "object",
            "required": [
                "login"
            ],
            "properties": {
                "login": {
                    "type": "string",
                    "minLength": 5,
                    "example": "name@example.com"
                }
            }
        },
        "request.ResetPassword": {
            "type": "object",
            "required": [
                "password"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "maxLength": 24,
                    "minLength": 5,
                    "example": "qwerty"
                }
            }
        },
        "response.AccessToken": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "expires_in": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "response.Message": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "response.Profile": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "response.URL": {
            "type": "object",
            "properties": {
                "url": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "SSO",
	Description:      "Single sign-on (сервис единого входа)",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
