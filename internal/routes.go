package internal

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// 檢查若有錯誤將錯誤信息返回給client
func ResponseErrorInfo(c *gin.Context, err error) {
	c.JSON(200, gin.H{
		"result": "fail",
		"error":  err.Error(),
	})
}

// request完成返回結果
func FinishReturnResult(c *gin.Context, result interface{}) {
	c.JSON(200, gin.H{
		"result": "ok",
		"data":   result,
	})
	return
}

func GetAllTemplatesRoute(c *gin.Context) {
	result, err := ListTemplates()
	if err != nil {
		sugarLogger.Errorf(err.Error())
		ResponseErrorInfo(c, err)
		return
	}
	FinishReturnResult(c, result)

}

func GetSingleTemplateRoute(c *gin.Context) {
	ID := RequestId{}
	err := c.Bind(&ID)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		ResponseErrorInfo(c, err)
		return
	}
	result, err := store.GetSingleTemplateTX(ID.Id)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		ResponseErrorInfo(c, err)
		return
	}
	FinishReturnResult(c, result)
}

// 新增的模板立即編譯成函數以檢查有沒錯誤, 若有錯誤刪除該模板
func CreateTemplateRoute(c *gin.Context) {
	payload := RequestTemplateOne{}
	err := c.Bind(&payload)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		ResponseErrorInfo(c, err)
		return
	}
	//result, err := store.CreateTemplateTX(payload)
	result, err := CreateTemplateRawTX(payload)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		ResponseErrorInfo(c, err)
		return
	}
	Trace.Println("CreateTemplate")
	// 測試新增template可否成功轉成parse函數
	SingleString, err := GenerateNewAddParseFuncString(result, payload)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		ResponseErrorInfo(c, err)
		return
	}
	err = EvalFuncString(SingleString, result)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		ResponseErrorInfo(c, err)
		return
	}
	FinishReturnResult(c, result)
}

func DeleteTemplateRoute(c *gin.Context) {
	ID := RequestId{}
	err := c.Bind(&ID)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		ResponseErrorInfo(c, err)
		return
	}
	err = store.DeleteTemplateTx(ID.Id)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		ResponseErrorInfo(c, err)
		return
	}
	FinishReturnResult(c, err)
}

func UpdateTemplateRoute(c *gin.Context) {
	payload := RequestTemplateOne{}
	err := c.Bind(&payload)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		ResponseErrorInfo(c, err)
		return
	}
	err = store.UpdateTemplateTx(payload)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		ResponseErrorInfo(c, err)
		return
	}
	FinishReturnResult(c, err)
}

func GetParseResultRoute(c *gin.Context) {
	payload := EvalPayload{}
	err := c.Bind(&payload)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		ResponseErrorInfo(c, err)
		return
	}
	result, err := GetParseResult(payload)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		ResponseErrorInfo(c, err)
		return
	}
	FinishReturnResult(c, result)
}

func InsertSensorDataRoute(c *gin.Context) {
	payload := EvalPayload{}
	err := c.Bind(&payload)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		ResponseErrorInfo(c, err)
		return
	}
	values, err := GetParseResult(payload)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		ResponseErrorInfo(c, err)
		return
	}
	if len(values) == 0 {
		sugarLogger.Errorf("got empty array")
		ResponseErrorInfo(c, errors.New("got empty array"))
		return
	}
	OriginObjectID, ObjectIDs, err := GetOriAndObjectIDs(payload)
	//	若ObjectIDs長度為0表示自動產生, 若不為0則必須與values長度一致
	if len(ObjectIDs) == 0 {
		for i, _ := range values {
			ObjectIDs = append(ObjectIDs, OriginObjectID+"_P"+strconv.Itoa(i))
		}
	} else {
		if len(ObjectIDs) != len(values) {
			sugarLogger.Errorf("values length mismatch")
			ResponseErrorInfo(c, errors.New("values length mismatch"))
			return
		}
	}
	part := InsertSenorDataSingle{
		System:    OriginObjectID,
		TimeStamp: int(time.Now().Unix()),
		TimeZone:  "",
		Data:      []InsideData{},
	}
	for i, v := range ObjectIDs {
		singleValue := InsideData{
			ObjectID:  v,
			Value:     fmt.Sprint(values[i]),
			TimeStamp: int(time.Now().Unix()),
		}
		part.Data = append(part.Data, singleValue)
	}
	result := InsertSenorDataArray{part}
	resp, err := InsertSensorDataPost(result)
	if err != nil {
		ResponseErrorInfo(c, err)
		return
	}
	FinishReturnResult(c, resp)
}
