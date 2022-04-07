package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"json/internal"
	"log"
	"os"
	"path/filepath"
)

var (
	//Trace *log.Logger
	Info *log.Logger
	//Error *log.Logger
)

func init() {
	newPath := filepath.Join(".", "log")
	err := os.MkdirAll(newPath, os.ModePerm)
	if err != nil {
		log.Fatal("can not create log folder")
	}
	file, err := os.OpenFile("./log/main.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("can not open log file")
	}
	Info = log.New(io.MultiWriter(file, os.Stdout), "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {

	Info.Println("Starting...")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/json_parser/V1/get_all_templates", internal.GetAllTemplatesRoute)
	r.POST("/json_parser/V1/get_single_template", internal.GetSingleTemplateRoute)
	r.POST("/json_parser/V1/create_template", internal.CreateTemplateRoute)
	r.POST("/json_parser/V1/delete_template", internal.DeleteTemplateRoute)
	r.POST("/json_parser/V1/update_template", internal.UpdateTemplateRoute)

	r.POST("/json_parser/V1/get_eval_result", internal.GetParseResultRoute)
	r.POST("/json_parser/V1/insert_sensor_data", internal.InsertSensorDataRoute)

	r.Run(":9109")
}
