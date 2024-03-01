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
        "/{user_id}/add": {
            "post": {
                "description": "Adds an  associated with user with given parameters. NotificationPeriod is optional and must look like {number}s,{number}m or {number}h. Implemented with the use of transaction: first rooms availibility is checked. In case one's new booking request intersects with and old one(even if belongs to him), the request is considered erratic. startDate is to be before endDate and both should not be expired.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "events"
                ],
                "summary": "Adds event",
                "operationId": "addByEventJSON",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "int64",
                        "default": 1234,
                        "description": "user_id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "AddEventRequest",
                        "name": "event",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.AddEventResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.AddEventResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.AddEventResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.AddEventResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.AddEventResponse"
                        }
                    }
                }
            }
        },
        "/{user_id}/get-events": {
            "get": {
                "description": "Responds with series of event info objects within given time period. The query parameters are start date and end date (start is to be before end and both should not be expired).",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "events"
                ],
                "summary": "Get several events info",
                "operationId": "getMultipleEventsByTag",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "int64",
                        "default": 1234,
                        "description": "user_id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "format": "time.Time",
                        "default": "2024-03-28T17:43:00-03:00",
                        "description": "start",
                        "name": "start",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "format": "time.Time",
                        "default": "2024-03-29T17:43:00-03:00",
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
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetEventsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetEventsResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetEventsResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetEventsResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetEventsResponse"
                        }
                    }
                }
            }
        },
        "/{user_id}/get-vacant-rooms": {
            "get": {
                "description": "Receives two dates as query parameters. start is to be before end and both should not be expired. Responds with list of vacant rooms and their parameters for given interval.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "events"
                ],
                "summary": "Get list of vacant rooms",
                "operationId": "getRoomsByDates",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "int64",
                        "default": 1234,
                        "description": "user_id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "format": "time.Time",
                        "default": "2024-03-28T17:43:00-03:00",
                        "description": "start",
                        "name": "start",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "format": "time.Time",
                        "default": "2024-03-29T17:43:00-03:00",
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
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetVacantRoomsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetVacantRoomsResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetVacantRoomsResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetVacantRoomsResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetVacantRoomsResponse"
                        }
                    }
                }
            }
        },
        "/{user_id}/{event_id}/delete": {
            "delete": {
                "description": "Deletes an event with given UUID.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "events"
                ],
                "summary": "Deletes an event",
                "operationId": "removeByEventID",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "int64",
                        "default": 1234,
                        "description": "user_id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    },
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
                            "$ref": "#/definitions/event-schedule_internal_app_api.DeleteEventResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.DeleteEventResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.DeleteEventResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.DeleteEventResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.DeleteEventResponse"
                        }
                    }
                }
            }
        },
        "/{user_id}/{event_id}/get": {
            "get": {
                "description": "Responds with booking info for booking with given EventID.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "events"
                ],
                "summary": "Get event info",
                "operationId": "getEventbyTag",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "int64",
                        "default": 1234,
                        "description": "user_id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    },
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
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetEventResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetEventResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetEventResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetEventResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetEventResponse"
                        }
                    }
                }
            }
        },
        "/{user_id}/{event_id}/update": {
            "patch": {
                "description": "Updates an existing event with given EventID, suiteID, startDate, endDate values (notificationPeriod being optional). Implemented with the use of transaction: first room availibility is checked. In case one attempts to alter his previous booking (i.e. widen or tighten its' limits) the booking is updated.  If it overlaps with smb else's booking or with clients' another booking the request is considered unsuccessful. startDate parameter  is to be before endDate and both should not be expired.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "events"
                ],
                "summary": "Updates event info",
                "operationId": "modifyEventByJSON",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "int64",
                        "default": 1234,
                        "description": "user_id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    },
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
                            "$ref": "#/definitions/event-schedule_internal_app_api.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.UpdateEventResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.UpdateEventResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.UpdateEventResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.UpdateEventResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.UpdateEventResponse"
                        }
                    }
                }
            }
        },
        "/{user_id}/{suite_id}/get-vacant-dates": {
            "get": {
                "description": "Responds with list of vacant intervals within month for selected suite.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "events"
                ],
                "summary": "Get vacant intervals",
                "operationId": "getDatesBySuiteID",
                "parameters": [
                    {
                        "type": "integer",
                        "format": "int64",
                        "default": 1234,
                        "description": "user_id",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "format": "int64",
                        "default": 1234,
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
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetVacantDatesResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetVacantDatesResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetVacantDatesResponse"
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetVacantDatesResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/event-schedule_internal_app_api.GetVacantDatesResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "event-schedule_internal_app_api.AddEventResponse": {
            "type": "object",
            "properties": {
                "eventID": {
                    "type": "string",
                    "format": "uuid",
                    "example": "550e8400-e29b-41d4-a716-446655440000"
                },
                "response": {
                    "$ref": "#/definitions/event-schedule_internal_app_api.Response"
                }
            }
        },
        "event-schedule_internal_app_api.DeleteEventResponse": {
            "type": "object",
            "properties": {
                "response": {
                    "$ref": "#/definitions/event-schedule_internal_app_api.Response"
                }
            }
        },
        "event-schedule_internal_app_api.GetEventResponse": {
            "type": "object",
            "properties": {
                "event": {
                    "$ref": "#/definitions/event-schedule_internal_app_model.EventInfo"
                },
                "response": {
                    "$ref": "#/definitions/event-schedule_internal_app_api.Response"
                }
            }
        },
        "event-schedule_internal_app_api.GetEventsResponse": {
            "type": "object",
            "properties": {
                "events": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/event-schedule_internal_app_model.EventInfo"
                    }
                },
                "response": {
                    "$ref": "#/definitions/event-schedule_internal_app_api.Response"
                }
            }
        },
        "event-schedule_internal_app_api.GetVacantDatesResponse": {
            "type": "object",
            "properties": {
                "intervals": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/event-schedule_internal_app_model.Interval"
                    }
                },
                "response": {
                    "$ref": "#/definitions/event-schedule_internal_app_api.Response"
                }
            }
        },
        "event-schedule_internal_app_api.GetVacantRoomsResponse": {
            "type": "object",
            "properties": {
                "response": {
                    "$ref": "#/definitions/event-schedule_internal_app_api.Response"
                },
                "rooms": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/event-schedule_internal_app_model.Suite"
                    }
                }
            }
        },
        "event-schedule_internal_app_api.Request": {
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
                    "example": "2024-03-29T17:43:00-03:00"
                },
                "notificationPeriod": {
                    "description": "Интервал времени для предварительного уведомления о бронировании",
                    "type": "string",
                    "example": "24h"
                },
                "startDate": {
                    "description": "Дата и время начала бронировании",
                    "type": "string",
                    "example": "2024-03-28T17:43:00-03:00"
                },
                "suiteID": {
                    "description": "номер апаратаментов",
                    "type": "integer",
                    "example": 123
                }
            }
        },
        "event-schedule_internal_app_api.Response": {
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
        "event-schedule_internal_app_api.UpdateEventResponse": {
            "type": "object",
            "properties": {
                "response": {
                    "$ref": "#/definitions/event-schedule_internal_app_api.Response"
                }
            }
        },
        "event-schedule_internal_app_model.EventInfo": {
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
                    "type": "string"
                },
                "endDate": {
                    "description": "Дата и время окончания бронировании",
                    "type": "string"
                },
                "notifyAt": {
                    "description": "Интервал времени для уведомления о бронировании",
                    "type": "string"
                },
                "ownerID": {
                    "description": "Идентификатор владельца бронирования",
                    "type": "integer"
                },
                "startDate": {
                    "description": "Дата и время начала бронировании",
                    "type": "string"
                },
                "suiteID": {
                    "description": "Номер апартаментов",
                    "type": "integer"
                },
                "updatedAt": {
                    "description": "Дата и время обновления",
                    "type": "string"
                }
            }
        },
        "event-schedule_internal_app_model.Interval": {
            "type": "object",
            "properties": {
                "end": {
                    "type": "string"
                },
                "start": {
                    "type": "string",
                    "example": "2024-03-02T15:04:05-07:00"
                }
            }
        },
        "event-schedule_internal_app_model.Suite": {
            "type": "object",
            "properties": {
                "capacity": {
                    "type": "integer",
                    "example": 4
                },
                "name": {
                    "type": "string",
                    "example": "Winston Churchill"
                },
                "suiteID": {
                    "type": "integer",
                    "example": 123
                }
            }
        }
    }
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
