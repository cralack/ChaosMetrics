// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/ARAM": {
            "get": {
                "description": "query @version,loc",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Champion Rank"
                ],
                "summary": "请求一个ARAM英雄榜",
                "parameters": [
                    {
                        "type": "string",
                        "default": "na1",
                        "description": "Region",
                        "name": "loc",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "14.1.1",
                        "description": "Version",
                        "name": "version",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/anres.ChampionBrief"
                            }
                        }
                    }
                }
            }
        },
        "/CLASSIC": {
            "get": {
                "description": "query @version,loc",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Champion Rank"
                ],
                "summary": "请求一个CLASSIC英雄榜",
                "parameters": [
                    {
                        "type": "string",
                        "default": "na1",
                        "description": "Region",
                        "name": "loc",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "14.1.1",
                        "description": "Version",
                        "name": "version",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/anres.ChampionBrief"
                            }
                        }
                    }
                }
            }
        },
        "/champion": {
            "get": {
                "description": "请求一个英雄详情 @name,version,loc,mode",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Champion Detail"
                ],
                "summary": "请求一个英雄详情",
                "parameters": [
                    {
                        "type": "string",
                        "default": "na1",
                        "description": "Region",
                        "name": "loc",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "CLASSIC",
                        "description": "Game mode",
                        "name": "mode",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "Ahri",
                        "description": "Champion name",
                        "name": "name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "14.1.1",
                        "description": "Version",
                        "name": "version",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/anres.ChampionDetail"
                        }
                    }
                }
            }
        },
        "/item": {
            "get": {
                "description": "请求一个物品详情 @version,lang",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Item"
                ],
                "summary": "请求一个物品详情",
                "parameters": [
                    {
                        "type": "string",
                        "default": "2010",
                        "description": "The ID of the item",
                        "name": "itemid",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "zh_CN",
                        "description": "Language",
                        "name": "lang",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "13.8.1",
                        "description": "Version",
                        "name": "version",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/riotmodel.ItemDTO"
                        }
                    }
                }
            }
        },
        "/summoner": {
            "get": {
                "description": "请求一个召唤师详情 @name,loc",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Champion Detail"
                ],
                "summary": "请求一个召唤师详情",
                "parameters": [
                    {
                        "type": "string",
                        "default": "na1",
                        "description": "Region",
                        "name": "loc",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "Solarbacca",
                        "description": "Summoner name",
                        "name": "name",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.SummonerDTO"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "anres.ChampionBrief": {
            "type": "object",
            "properties": {
                "avg_damage_dealt": {
                    "description": "场均输出占比 12%",
                    "type": "number"
                },
                "avg_dead_time": {
                    "description": "场均死亡时长",
                    "type": "number"
                },
                "ban_rate": {
                    "description": "Ban率",
                    "type": "number"
                },
                "id": {
                    "description": "英雄ID:Aatrox",
                    "type": "string"
                },
                "image": {
                    "description": "basic info",
                    "allOf": [
                        {
                            "$ref": "#/definitions/model.Image"
                        }
                    ]
                },
                "pick_rate": {
                    "description": "登场率 15%",
                    "type": "number"
                },
                "win_rate": {
                    "description": "statistical data",
                    "type": "number"
                }
            }
        },
        "anres.ChampionDetail": {
            "type": "object",
            "properties": {
                "avg_damage_dealt": {
                    "description": "场均输出占比 12%",
                    "type": "number"
                },
                "avg_damage_taken": {
                    "description": "场均承伤占比 10%",
                    "type": "number"
                },
                "avg_dead_time": {
                    "description": "场均死亡时长",
                    "type": "number"
                },
                "avg_kda": {
                    "description": "场均KDA 15%",
                    "type": "number"
                },
                "avg_kp": {
                    "description": "场均参团率 10%",
                    "type": "number"
                },
                "avg_time_ccing": {
                    "description": "场均控制时长 5%",
                    "type": "number"
                },
                "avg_vision_score": {
                    "description": "场均视野得分 3%",
                    "type": "number"
                },
                "ban_rate": {
                    "description": "Ban率",
                    "type": "number"
                },
                "game_mode": {
                    "description": "游戏模式",
                    "type": "string"
                },
                "id": {
                    "description": "英雄ID:Aatrox",
                    "type": "string"
                },
                "idx": {
                    "description": "basic info",
                    "type": "string"
                },
                "image": {
                    "description": "图像",
                    "allOf": [
                        {
                            "$ref": "#/definitions/model.Image"
                        }
                    ]
                },
                "item": {
                    "description": "build winrate",
                    "type": "object",
                    "additionalProperties": {
                        "type": "object",
                        "additionalProperties": {
                            "type": "integer"
                        }
                    }
                },
                "key": {
                    "description": "英雄Key:266",
                    "type": "string"
                },
                "loc": {
                    "description": "对局服务器",
                    "type": "string"
                },
                "name": {
                    "description": "英雄名称:暗裔剑魔",
                    "type": "string"
                },
                "perk": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "integer"
                    }
                },
                "pick_rate": {
                    "description": "登场率 15%",
                    "type": "number"
                },
                "rank_score": {
                    "description": "statistical data",
                    "type": "number"
                },
                "skill": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "integer"
                    }
                },
                "spell": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "integer"
                    }
                },
                "title": {
                    "description": "英雄称号:亚托克斯",
                    "type": "string"
                },
                "total_played": {
                    "description": "英雄出场总数",
                    "type": "number"
                },
                "total_win": {
                    "description": "英雄胜利总数",
                    "type": "number"
                },
                "version": {
                    "description": "对局版本",
                    "type": "string"
                },
                "win_rate": {
                    "description": "胜率 30%",
                    "type": "number"
                }
            }
        },
        "model.Image": {
            "type": "object",
            "properties": {
                "full": {
                    "description": "大图文件名",
                    "type": "string"
                },
                "group": {
                    "description": "图像所属组",
                    "type": "string"
                },
                "h": {
                    "description": "图像高度",
                    "type": "integer"
                },
                "sprite": {
                    "description": "小图文件名",
                    "type": "string"
                },
                "w": {
                    "description": "图像宽度",
                    "type": "integer"
                },
                "x": {
                    "description": "图像 X 坐标",
                    "type": "integer"
                },
                "y": {
                    "description": "图像 Y 坐标",
                    "type": "integer"
                }
            }
        },
        "response.SummonerDTO": {
            "type": "object",
            "properties": {
                "loc": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "profileIconID": {
                    "type": "integer"
                },
                "summonerLevel": {
                    "type": "integer"
                }
            }
        },
        "riotmodel.GoldInfo": {
            "type": "object",
            "properties": {
                "base": {
                    "type": "integer"
                },
                "purchasable": {
                    "type": "boolean"
                },
                "sell": {
                    "type": "integer"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "riotmodel.ItemDTO": {
            "type": "object",
            "properties": {
                "colloq": {
                    "type": "string"
                },
                "depth": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "from": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "gold": {
                    "$ref": "#/definitions/riotmodel.GoldInfo"
                },
                "id": {
                    "type": "string"
                },
                "image": {
                    "$ref": "#/definitions/model.Image"
                },
                "into": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "maps": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "boolean"
                    }
                },
                "name": {
                    "type": "string"
                },
                "plaintext": {
                    "type": "string"
                },
                "stats": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "number"
                    }
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "BasicAuth": {
            "type": "basic"
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.7",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "ChaosMetrics API接口文档",
	Description:      "使用Riot官方API获取数据进行分析、统计项目",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
