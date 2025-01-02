// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://example.com/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.example.com/support",
            "email": "support@example.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/matcher/api/v1/search": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "Find the nearest driver around a GeoJSON point",
                "tags": [
                    "Matcher"
                ],
                "summary": "Search for a Driver",
                "parameters": [
                    {
                        "type": "number",
                        "description": "Latitude",
                        "name": "latitude",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Longitude",
                        "name": "longitude",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Search radius in meters",
                        "name": "radius",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/bitaksi-go-matcher_internal_models.DriverWithDistance"
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Driver not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "bitaksi-go-matcher_internal_models.DriverWithDistance": {
            "type": "object",
            "properties": {
                "distance": {
                    "description": "Distance from the \"near\" point",
                    "type": "number"
                },
                "id": {
                    "type": "string"
                },
                "location": {
                    "$ref": "#/definitions/bitaksi-go-matcher_internal_models.Location"
                }
            }
        },
        "bitaksi-go-matcher_internal_models.Location": {
            "type": "object",
            "properties": {
                "coordinates": {
                    "type": "array",
                    "items": {
                        "type": "number"
                    }
                },
                "type": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Matcher Service API",
	Description:      "This is the API documentation for the Matcher Service.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}