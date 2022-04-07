package internal

import (
	"json/pkg/db"
)

type Config struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	DB       string `mapstructure:"DB_DB"`
	SeverIP  string `mapstructure:"SERVER_IP"`
	Token    string `mapstructure:"TOKEN"`
}

// 刪除使用
type RequestId struct {
	Id int32 `json:"id"`
}

// 查詢單一template與新增與修改使用
type RequestTemplateOne struct {
	db.CreateTemplateParams
	ObjectIDList []string `json:"object_id_list"`
	TemplateID   int32    `json:"template_id,omitempty"`
	Name string `json:"name"`
}

// template清單
type ResponseTemplates struct {
	TemplateID   int32  `json:"template_id"`
	TemplateName string `json:"template_name"`
	Description  string `json:"description"`
}

type EvalPayload struct {
	Body []byte `json:"body"`
	Id   int    `json:"id"`
}

type InsideData struct {
	ObjectID  string `json:"ObjectId"`
	Value     string `json:"Value"`
	TimeStamp int    `json:"TimeStamp"`
}
type InsertSenorDataSingle struct {
	System    string       `json:"System"`
	TimeStamp int          `json:"TimeStamp"`
	TimeZone  string       `json:"TimeZone"`
	Data      []InsideData `json:"Data"`
}

type InsertSenorDataArray []InsertSenorDataSingle
