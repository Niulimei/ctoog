// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "title": "translator",
    "termsOfService": "http://swagger.io/terms/",
    "contact": {
      "email": "apiteam@swagger.io"
    },
    "license": {
      "name": "Apache 2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    },
    "version": "1.0.0"
  },
  "basePath": "/api",
  "paths": {
    "/login": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "登录",
        "operationId": "Login",
        "parameters": [
          {
            "description": "用户信息",
            "name": "userInfo",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/UserModel"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "成功",
            "schema": {
              "$ref": "#/definitions/Authorization"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/pvobs": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "pvob列表",
        "operationId": "ListPvob",
        "responses": {
          "200": {
            "description": "PVOB列表",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/pvobs/{id}/components": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "组件列表",
        "operationId": "ListPvobComponent",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "组件列表",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/pvobs/{pvob_id}/components/{component_id}": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "流列表",
        "operationId": "ListPvobComponentStream",
        "parameters": [
          {
            "type": "string",
            "name": "pvob_id",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "name": "component_id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "流列表",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/tasks": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "任务列表",
        "operationId": "ListTask",
        "parameters": [
          {
            "type": "integer",
            "default": 0,
            "name": "limit",
            "in": "query",
            "required": true
          },
          {
            "type": "integer",
            "default": 0,
            "name": "offset",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "Authorization",
            "in": "header",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "任务列表",
            "schema": {
              "$ref": "#/definitions/TaskPageInfoModel"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      },
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "新建任务",
        "operationId": "CreateTask",
        "parameters": [
          {
            "description": "任务信息",
            "name": "taskInfo",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/TaskModel"
            }
          },
          {
            "type": "string",
            "name": "Authorization",
            "in": "header",
            "required": true
          }
        ],
        "responses": {
          "201": {
            "description": "成功",
            "schema": {
              "$ref": "#/definitions/OK"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/tasks/restart": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "任务重启",
        "operationId": "RestartTask",
        "parameters": [
          {
            "type": "string",
            "name": "Authorization",
            "in": "header",
            "required": true
          },
          {
            "description": "任务重启",
            "name": "restartTrigger",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/TaskRestart"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "任务列表",
            "schema": {
              "$ref": "#/definitions/TaskPageInfoModel"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/tasks/{id}": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "任务详情",
        "operationId": "GetTask",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "name": "Authorization",
            "in": "header",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "任务信息",
            "schema": {
              "$ref": "#/definitions/TaskDetail"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      },
      "put": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "更新任务",
        "operationId": "UpdateTask",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "name": "Authorization",
            "in": "header",
            "required": true
          },
          {
            "description": "任务信息",
            "name": "taskLog",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/TaskLogInfo"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "成功",
            "schema": {
              "$ref": "#/definitions/OK"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/users": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "账户列表",
        "operationId": "ListUser",
        "parameters": [
          {
            "type": "integer",
            "default": 0,
            "name": "limit",
            "in": "query",
            "required": true
          },
          {
            "type": "integer",
            "default": 0,
            "name": "offset",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "Authorization",
            "in": "header",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "账户列表",
            "schema": {
              "$ref": "#/definitions/UserPageInfoModel"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      },
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "开通账户",
        "operationId": "CreateUser",
        "parameters": [
          {
            "description": "用户信息",
            "name": "userInfo",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/UserModel"
            }
          },
          {
            "type": "string",
            "name": "Authorization",
            "in": "header",
            "required": true
          }
        ],
        "responses": {
          "201": {
            "description": "成功",
            "schema": {
              "$ref": "#/definitions/OK"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/users/self": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "账户",
        "operationId": "GetUser",
        "parameters": [
          {
            "type": "string",
            "name": "Authorization",
            "in": "header",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "账户",
            "schema": {
              "$ref": "#/definitions/UserInfoModel"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/workers": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "worker注册和心跳",
        "operationId": "PingWorker",
        "parameters": [
          {
            "description": "worker信息",
            "name": "workerInfo",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/WorkerModel"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "成功",
            "schema": {
              "$ref": "#/definitions/OK"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Authorization": {
      "type": "object",
      "properties": {
        "token": {
          "type": "string"
        }
      }
    },
    "ErrorModel": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer",
          "example": 400
        },
        "message": {
          "type": "string",
          "example": "error message"
        }
      }
    },
    "OK": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string"
        }
      }
    },
    "TaskDetail": {
      "type": "object",
      "properties": {
        "logList": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/TaskLogInfo"
          }
        },
        "taskModel": {
          "$ref": "#/definitions/TaskModel"
        }
      }
    },
    "TaskInfoModel": {
      "type": "object",
      "properties": {
        "component": {
          "type": "string"
        },
        "gitRepo": {
          "type": "string"
        },
        "id": {
          "type": "integer"
        },
        "lastCompleteDateTime": {
          "type": "string"
        },
        "pvob": {
          "type": "string"
        },
        "status": {
          "type": "string"
        }
      }
    },
    "TaskLogInfo": {
      "type": "object",
      "properties": {
        "duration": {
          "type": "string"
        },
        "endTime": {
          "type": "string"
        },
        "logID": {
          "type": "string"
        },
        "startTime": {
          "type": "string"
        },
        "status": {
          "type": "string"
        }
      }
    },
    "TaskMatchInfo": {
      "type": "object",
      "properties": {
        "gitBranch": {
          "type": "string"
        },
        "stream": {
          "type": "string"
        }
      }
    },
    "TaskModel": {
      "type": "object",
      "properties": {
        "ccPassword": {
          "type": "string"
        },
        "ccUser": {
          "type": "string"
        },
        "component": {
          "type": "string"
        },
        "gitPassword": {
          "type": "string"
        },
        "gitURL": {
          "type": "string"
        },
        "gitUser": {
          "type": "string"
        },
        "matchInfo": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/TaskMatchInfo"
          }
        },
        "pvob": {
          "type": "string"
        }
      }
    },
    "TaskPageInfoModel": {
      "type": "object",
      "properties": {
        "count": {
          "type": "integer"
        },
        "limit": {
          "type": "integer"
        },
        "offset": {
          "type": "integer"
        },
        "taskInfo": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/TaskInfoModel"
          }
        }
      }
    },
    "TaskRestart": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string"
        }
      }
    },
    "UserInfoModel": {
      "type": "object",
      "properties": {
        "id": {
          "type": "integer"
        },
        "password": {
          "type": "string"
        },
        "role_id": {
          "type": "integer"
        },
        "username": {
          "type": "string"
        }
      }
    },
    "UserModel": {
      "type": "object",
      "properties": {
        "password": {
          "type": "string"
        },
        "role_id": {
          "type": "integer"
        },
        "username": {
          "type": "string"
        }
      }
    },
    "UserPageInfoModel": {
      "type": "object",
      "properties": {
        "count": {
          "type": "integer"
        },
        "limit": {
          "type": "integer"
        },
        "offset": {
          "type": "integer"
        },
        "userInfo": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/UserInfoModel"
          }
        }
      }
    },
    "WorkerModel": {
      "type": "object",
      "properties": {
        "host": {
          "type": "string",
          "example": "192.168.1.1"
        },
        "port": {
          "type": "integer",
          "example": 80
        }
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "title": "translator",
    "termsOfService": "http://swagger.io/terms/",
    "contact": {
      "email": "apiteam@swagger.io"
    },
    "license": {
      "name": "Apache 2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    },
    "version": "1.0.0"
  },
  "basePath": "/api",
  "paths": {
    "/login": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "登录",
        "operationId": "Login",
        "parameters": [
          {
            "description": "用户信息",
            "name": "userInfo",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/UserModel"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "成功",
            "schema": {
              "$ref": "#/definitions/Authorization"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/pvobs": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "pvob列表",
        "operationId": "ListPvob",
        "responses": {
          "200": {
            "description": "PVOB列表",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/pvobs/{id}/components": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "组件列表",
        "operationId": "ListPvobComponent",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "组件列表",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/pvobs/{pvob_id}/components/{component_id}": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "流列表",
        "operationId": "ListPvobComponentStream",
        "parameters": [
          {
            "type": "string",
            "name": "pvob_id",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "name": "component_id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "流列表",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/tasks": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "任务列表",
        "operationId": "ListTask",
        "parameters": [
          {
            "type": "integer",
            "default": 0,
            "name": "limit",
            "in": "query",
            "required": true
          },
          {
            "type": "integer",
            "default": 0,
            "name": "offset",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "Authorization",
            "in": "header",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "任务列表",
            "schema": {
              "$ref": "#/definitions/TaskPageInfoModel"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      },
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "新建任务",
        "operationId": "CreateTask",
        "parameters": [
          {
            "description": "任务信息",
            "name": "taskInfo",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/TaskModel"
            }
          },
          {
            "type": "string",
            "name": "Authorization",
            "in": "header",
            "required": true
          }
        ],
        "responses": {
          "201": {
            "description": "成功",
            "schema": {
              "$ref": "#/definitions/OK"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/tasks/restart": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "任务重启",
        "operationId": "RestartTask",
        "parameters": [
          {
            "type": "string",
            "name": "Authorization",
            "in": "header",
            "required": true
          },
          {
            "description": "任务重启",
            "name": "restartTrigger",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/TaskRestart"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "任务列表",
            "schema": {
              "$ref": "#/definitions/TaskPageInfoModel"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/tasks/{id}": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "任务详情",
        "operationId": "GetTask",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "name": "Authorization",
            "in": "header",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "任务信息",
            "schema": {
              "$ref": "#/definitions/TaskDetail"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      },
      "put": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "更新任务",
        "operationId": "UpdateTask",
        "parameters": [
          {
            "type": "string",
            "name": "id",
            "in": "path",
            "required": true
          },
          {
            "type": "string",
            "name": "Authorization",
            "in": "header",
            "required": true
          },
          {
            "description": "任务信息",
            "name": "taskLog",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/TaskLogInfo"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "成功",
            "schema": {
              "$ref": "#/definitions/OK"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/users": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "账户列表",
        "operationId": "ListUser",
        "parameters": [
          {
            "type": "integer",
            "default": 0,
            "name": "limit",
            "in": "query",
            "required": true
          },
          {
            "type": "integer",
            "default": 0,
            "name": "offset",
            "in": "query",
            "required": true
          },
          {
            "type": "string",
            "name": "Authorization",
            "in": "header",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "账户列表",
            "schema": {
              "$ref": "#/definitions/UserPageInfoModel"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      },
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "开通账户",
        "operationId": "CreateUser",
        "parameters": [
          {
            "description": "用户信息",
            "name": "userInfo",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/UserModel"
            }
          },
          {
            "type": "string",
            "name": "Authorization",
            "in": "header",
            "required": true
          }
        ],
        "responses": {
          "201": {
            "description": "成功",
            "schema": {
              "$ref": "#/definitions/OK"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/users/self": {
      "get": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "账户",
        "operationId": "GetUser",
        "parameters": [
          {
            "type": "string",
            "name": "Authorization",
            "in": "header",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "账户",
            "schema": {
              "$ref": "#/definitions/UserInfoModel"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    },
    "/workers": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "worker注册和心跳",
        "operationId": "PingWorker",
        "parameters": [
          {
            "description": "worker信息",
            "name": "workerInfo",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/WorkerModel"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "成功",
            "schema": {
              "$ref": "#/definitions/OK"
            }
          },
          "500": {
            "description": "内部错误",
            "schema": {
              "$ref": "#/definitions/ErrorModel"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Authorization": {
      "type": "object",
      "properties": {
        "token": {
          "type": "string"
        }
      }
    },
    "ErrorModel": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer",
          "example": 400
        },
        "message": {
          "type": "string",
          "example": "error message"
        }
      }
    },
    "OK": {
      "type": "object",
      "properties": {
        "message": {
          "type": "string"
        }
      }
    },
    "TaskDetail": {
      "type": "object",
      "properties": {
        "logList": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/TaskLogInfo"
          }
        },
        "taskModel": {
          "$ref": "#/definitions/TaskModel"
        }
      }
    },
    "TaskInfoModel": {
      "type": "object",
      "properties": {
        "component": {
          "type": "string"
        },
        "gitRepo": {
          "type": "string"
        },
        "id": {
          "type": "integer"
        },
        "lastCompleteDateTime": {
          "type": "string"
        },
        "pvob": {
          "type": "string"
        },
        "status": {
          "type": "string"
        }
      }
    },
    "TaskLogInfo": {
      "type": "object",
      "properties": {
        "duration": {
          "type": "string"
        },
        "endTime": {
          "type": "string"
        },
        "logID": {
          "type": "string"
        },
        "startTime": {
          "type": "string"
        },
        "status": {
          "type": "string"
        }
      }
    },
    "TaskMatchInfo": {
      "type": "object",
      "properties": {
        "gitBranch": {
          "type": "string"
        },
        "stream": {
          "type": "string"
        }
      }
    },
    "TaskModel": {
      "type": "object",
      "properties": {
        "ccPassword": {
          "type": "string"
        },
        "ccUser": {
          "type": "string"
        },
        "component": {
          "type": "string"
        },
        "gitPassword": {
          "type": "string"
        },
        "gitURL": {
          "type": "string"
        },
        "gitUser": {
          "type": "string"
        },
        "matchInfo": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/TaskMatchInfo"
          }
        },
        "pvob": {
          "type": "string"
        }
      }
    },
    "TaskPageInfoModel": {
      "type": "object",
      "properties": {
        "count": {
          "type": "integer"
        },
        "limit": {
          "type": "integer"
        },
        "offset": {
          "type": "integer"
        },
        "taskInfo": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/TaskInfoModel"
          }
        }
      }
    },
    "TaskRestart": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string"
        }
      }
    },
    "UserInfoModel": {
      "type": "object",
      "properties": {
        "id": {
          "type": "integer"
        },
        "password": {
          "type": "string"
        },
        "role_id": {
          "type": "integer"
        },
        "username": {
          "type": "string"
        }
      }
    },
    "UserModel": {
      "type": "object",
      "properties": {
        "password": {
          "type": "string"
        },
        "role_id": {
          "type": "integer"
        },
        "username": {
          "type": "string"
        }
      }
    },
    "UserPageInfoModel": {
      "type": "object",
      "properties": {
        "count": {
          "type": "integer"
        },
        "limit": {
          "type": "integer"
        },
        "offset": {
          "type": "integer"
        },
        "userInfo": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/UserInfoModel"
          }
        }
      }
    },
    "WorkerModel": {
      "type": "object",
      "properties": {
        "host": {
          "type": "string",
          "example": "192.168.1.1"
        },
        "port": {
          "type": "integer",
          "example": 80
        }
      }
    }
  }
}`))
}
