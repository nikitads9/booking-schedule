// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Nikita Denisenok",
            "url": "https://vk.com/ndenisenok"
        },
        "license": {
            "name": "GNU 3.0",
            "url": "https://www.gnu.org/licenses/gpl-3.0.ru.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/add": {
            "post": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Adds an  associated with user with given parameters. NotificationPeriod is optional and must look like {number}s,{number}m or {number}h. Implemented with the use of transaction: first rooms availibility is checked. In case one's new booking request intersects with and old one(even if belongs to him), the request is considered erratic. startDate is to be before endDate and both should not be expired.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bookings"
                ],
                "summary": "Adds event",
                "operationId": "addByEventJSON",
                "parameters": [
                    {
                        "description": "AddEventRequest",
                        "name": "event",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/AddEventRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/AddEventResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/AddEventResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/AddEventResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/AddEventResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/AddEventResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/AddEventResponse"
                        }
                    }
                }
            }
        },
        "/get-events": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Responds with series of event info objects within given time period. The query parameters are start date and end date (start is to be before end and both should not be expired).",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bookings"
                ],
                "summary": "Get several events info",
                "operationId": "getMultipleEventsByTag",
                "parameters": [
                    {
                        "type": "string",
                        "format": "time.Time",
                        "default": "2024-03-28T17:43:00Z",
                        "description": "start",
                        "name": "start",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "format": "time.Time",
                        "default": "2024-03-29T17:43:00Z",
                        "description": "end",
                        "name": "end",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/GetEventsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/GetEventsResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/GetEventsResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/GetEventsResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/GetEventsResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/GetEventsResponse"
                        }
                    }
                }
            }
        },
        "/get-vacant-rooms": {
            "get": {
                "description": "Receives two dates as query parameters. start is to be before end and both should not be expired. Responds with list of vacant rooms and their parameters for given interval.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bookings"
                ],
                "summary": "Get list of vacant rooms",
                "operationId": "getRoomsByDates",
                "parameters": [
                    {
                        "type": "string",
                        "format": "time.Time",
                        "default": "2024-03-28T17:43:00Z",
                        "description": "start",
                        "name": "start",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "format": "time.Time",
                        "default": "2024-03-29T17:43:00Z",
                        "description": "end",
                        "name": "end",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/GetVacantRoomsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/GetVacantRoomsResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/GetVacantRoomsResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/GetVacantRoomsResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/GetVacantRoomsResponse"
                        }
                    }
                }
            }
        },
        "/user/sign-in": {
            "get": {
                "security": [
                    {
                        "BasicAuth": []
                    }
                ],
                "description": "Get auth token to access user restricted api methods. Requires nickname and password passed via basic auth.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Sign in",
                "operationId": "getOauthToken",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SignInResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SignInResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/SignInResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/SignInResponse"
                        }
                    }
                }
            }
        },
        "/user/sign-up": {
            "post": {
                "description": "Creates user with given tg id, nickname, name and password hashed by bcrypto. Every parameter is required. Returns jwt token.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Sign up",
                "operationId": "signUpUserJson",
                "parameters": [
                    {
                        "description": "User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/SignInResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/SignInResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/SignInResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/SignInResponse"
                        }
                    }
                }
            }
        },
        "/{event_id}/delete": {
            "delete": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Deletes an event with given UUID.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bookings"
                ],
                "summary": "Deletes an event",
                "operationId": "removeByEventID",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "default": "550e8400-e29b-41d4-a716-446655440000",
                        "description": "event_id",
                        "name": "event_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/DeleteEventResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/DeleteEventResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/DeleteEventResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/DeleteEventResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/DeleteEventResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/DeleteEventResponse"
                        }
                    }
                }
            }
        },
        "/{event_id}/get": {
            "get": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Responds with booking info for booking with given EventID.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bookings"
                ],
                "summary": "Get event info",
                "operationId": "getEventbyTag",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "default": "550e8400-e29b-41d4-a716-446655440000",
                        "description": "event_id",
                        "name": "event_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/GetEventResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/GetEventResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/GetEventResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/GetEventResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/GetEventResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/GetEventResponse"
                        }
                    }
                }
            }
        },
        "/{event_id}/update": {
            "patch": {
                "security": [
                    {
                        "Bearer": []
                    }
                ],
                "description": "Updates an existing event with given EventID, suiteID, startDate, endDate values (notificationPeriod being optional). Implemented with the use of transaction: first room availibility is checked. In case one attempts to alter his previous booking (i.e. widen or tighten its' limits) the booking is updated.  If it overlaps with smb else's booking or with clients' another booking the request is considered unsuccessful. startDate parameter  is to be before endDate and both should not be expired.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bookings"
                ],
                "summary": "Updates event info",
                "operationId": "modifyEventByJSON",
                "parameters": [
                    {
                        "type": "string",
                        "format": "uuid",
                        "default": "550e8400-e29b-41d4-a716-446655440000",
                        "description": "event_id",
                        "name": "event_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "UpdateEventRequest",
                        "name": "event",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/UpdateEventRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/UpdateEventResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/UpdateEventResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/UpdateEventResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/UpdateEventResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/UpdateEventResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/UpdateEventResponse"
                        }
                    }
                }
            }
        },
        "/{suite_id}/get-vacant-dates": {
            "get": {
                "description": "Responds with list of vacant intervals within month for selected suite.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "bookings"
                ],
                "summary": "Get vacant intervals",
                "operationId": "getDatesBySuiteID",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "int64",
                        "default": 1,
                        "description": "suite_id",
                        "name": "suite_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/GetVacantDateResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/GetVacantDateResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/GetVacantDateResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/GetVacantDateResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/GetVacantDateResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "AddEventRequest": {
            "type": "object",
            "required": [
                "endDate",
                "startDate",
                "suiteID"
            ],
            "properties": {
                "endDate": {
                    "description": "Дата и время окончания бронировании",
                    "type": "string",
                    "example": "2024-03-29T17:43:00Z"
                },
                "notifyAt": {
                    "description": "Интервал времени для предварительного уведомления о бронировании",
                    "type": "string",
                    "example": "24h"
                },
                "startDate": {
                    "description": "Дата и время начала бронировании",
                    "type": "string",
                    "example": "2024-03-28T17:43:00Z"
                },
                "suiteID": {
                    "description": "Номер апаратаментов",
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "AddEventResponse": {
            "type": "object",
            "properties": {
                "eventID": {
                    "type": "string",
                    "format": "uuid",
                    "example": "550e8400-e29b-41d4-a716-446655440000"
                },
                "response": {
                    "$ref": "#/definitions/Response"
                }
            }
        },
        "DeleteEventResponse": {
            "type": "object",
            "properties": {
                "response": {
                    "$ref": "#/definitions/Response"
                }
            }
        },
        "EventInfo": {
            "type": "object",
            "properties": {
                "EventID": {
                    "description": "Уникальный идентификатор бронирования",
                    "type": "string",
                    "format": "uuid",
                    "example": "550e8400-e29b-41d4-a716-446655440000"
                },
                "createdAt": {
                    "description": "Дата и время создания",
                    "type": "string",
                    "example": "2024-03-27T17:43:00Z"
                },
                "endDate": {
                    "description": "Дата и время окончания бронировании",
                    "type": "string",
                    "example": "2024-03-29T17:43:00Z"
                },
                "notifyAt": {
                    "description": "Интервал времени для уведомления о бронировании",
                    "type": "string",
                    "example": "24h00m00s"
                },
                "startDate": {
                    "description": "Дата и время начала бронировании",
                    "type": "string",
                    "example": "2024-03-28T17:43:00Z"
                },
                "suiteID": {
                    "description": "Номер апартаментов",
                    "type": "integer",
                    "example": 1
                },
                "updatedAt": {
                    "description": "Дата и время обновления",
                    "type": "string",
                    "example": "2024-03-27T18:43:00Z"
                },
                "userID": {
                    "description": "Идентификатор владельца бронирования",
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "GetEventResponse": {
            "type": "object",
            "properties": {
                "event": {
                    "$ref": "#/definitions/EventInfo"
                },
                "response": {
                    "$ref": "#/definitions/Response"
                }
            }
        },
        "GetEventsResponse": {
            "type": "object",
            "properties": {
                "events": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/EventInfo"
                    }
                },
                "response": {
                    "$ref": "#/definitions/Response"
                }
            }
        },
        "GetVacantDateResponse": {
            "type": "object",
            "properties": {
                "intervals": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Interval"
                    }
                },
                "response": {
                    "$ref": "#/definitions/Response"
                }
            }
        },
        "GetVacantRoomsResponse": {
            "type": "object",
            "properties": {
                "response": {
                    "$ref": "#/definitions/Response"
                },
                "rooms": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/Suite"
                    }
                }
            }
        },
        "Interval": {
            "type": "object",
            "properties": {
                "end": {
                    "description": "Номер свободен по",
                    "type": "string",
                    "example": "2024-04-10T15:04:05Z"
                },
                "start": {
                    "description": "Номер свободен с",
                    "type": "string",
                    "example": "2024-03-10T15:04:05Z"
                }
            }
        },
        "Response": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "application-specific error code",
                    "type": "integer"
                },
                "error": {
                    "description": "application-level error message, for debugging",
                    "type": "string"
                },
                "status": {
                    "description": "user-level status message",
                    "type": "string"
                }
            }
        },
        "SignInResponse": {
            "type": "object",
            "properties": {
                "response": {
                    "$ref": "#/definitions/Response"
                },
                "token": {
                    "description": "JWT токен для доступа",
                    "type": "string"
                }
            }
        },
        "Suite": {
            "type": "object",
            "properties": {
                "capacity": {
                    "description": "Вместимость в персонах",
                    "type": "integer",
                    "example": 4
                },
                "name": {
                    "description": "Название апартаментов",
                    "type": "string",
                    "example": "Winston Churchill"
                },
                "suiteID": {
                    "description": "Номер апартаментов",
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "UpdateEventRequest": {
            "type": "object",
            "required": [
                "endDate",
                "startDate",
                "suiteID"
            ],
            "properties": {
                "endDate": {
                    "description": "Дата и время окончания бронирования",
                    "type": "string",
                    "example": "2024-03-29T17:43:00-03:00"
                },
                "notifyAt": {
                    "description": "Интервал времени для уведомления о бронировании",
                    "type": "string",
                    "example": "24h"
                },
                "startDate": {
                    "description": "Дата и время начала бронировании",
                    "type": "string",
                    "example": "2024-03-28T17:43:00-03:00"
                },
                "suiteID": {
                    "description": "Номер апартаментов",
                    "type": "integer",
                    "format": "int64",
                    "example": 1
                }
            }
        },
        "UpdateEventResponse": {
            "type": "object",
            "properties": {
                "response": {
                    "$ref": "#/definitions/Response"
                }
            }
        },
        "User": {
            "type": "object",
            "required": [
                "name",
                "password",
                "telegramID",
                "telegramNickname"
            ],
            "properties": {
                "name": {
                    "description": "Имя пользователя",
                    "type": "string",
                    "example": "Pavel Durov"
                },
                "password": {
                    "description": "Пароль",
                    "type": "string",
                    "example": "12345"
                },
                "telegramID": {
                    "description": "Телеграм ID пользователя",
                    "type": "integer",
                    "example": 1235678
                },
                "telegramNickname": {
                    "description": "Никнейм пользователя в телеграме",
                    "type": "string",
                    "example": "pavel_durov"
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        },
        "Bearer": {
            "description": "Type \"Bearer\" followed by a space and JWT token.",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "tags": [
        {
            "description": "operations with bookings, suites and intervals",
            "name": "bookings"
        },
        {
            "description": "operations with user profile such as sign in, sign up, getting profile editing it and deleting",
            "name": "users"
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "127.0.0.1:3000",
	BasePath:         "/events",
	Schemes:          []string{"http", "https"},
	Title:            "event-schedule API",
	Description:      "This is a service for writing and reading booking entries.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
