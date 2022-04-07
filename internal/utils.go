package internal

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	//"io"
	"io/ioutil"
	"json/pkg/db"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var (
	Trace    = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	TxOption = sql.TxOptions{
		Isolation: 6,
		ReadOnly:  false,
	}

	Configs      = LoadConfig("./config")
	DBConnection = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&parseTime=true&loc=Local",
		Configs.User, Configs.Password, Configs.Host, Configs.Port, Configs.DB)
	MyDB, _    = sql.Open("mysql", DBConnection)
	ctx        = context.Background()
	queries    = db.New(MyDB)
	store      = NewStore(MyDB)
	_, b, _, _ = runtime.Caller(0)
	rootPath   = filepath.Dir(filepath.Dir(b))
	// 把取值的函數裝在這裏面, key為TemplateID
	EvalFuncHub map[string]func([]byte) ([]string, error)
	UniqueIds   = make([]string, 0, 0)

	sugarLogger *zap.SugaredLogger
)

// 讀專案中的config檔
func LoadConfig(MyPath string) (config Config) {
	// 若有同名環境變量則使用環境變量
	viper.AutomaticEnv()
	viper.AddConfigPath(MyPath)
	// 為了讓執行test也能讀到config添加索引路徑
	wd, err := os.Getwd()
	parent := filepath.Dir(wd)
	viper.AddConfigPath(path.Join(parent, MyPath))
	viper.SetConfigName("db")
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("can not load config: " + err.Error())
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("can not load config: " + err.Error())
	}
	Trace.Printf("%+v\n", config)
	return
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "log/internal.log",
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
func init() {
	InitLogger()
	defer sugarLogger.Sync()
	CreateParseFunction()
	sugarLogger.Infof("finish create parse function")
}

// 服務初始化與更新時調用
func CreateParseFunction() {
	EvalFuncHub = make(map[string]func([]byte) ([]string, error))
	funcString, err := GenerateAllParseFuncString()
	if err != nil {
		sugarLogger.Errorf(err.Error())
		log.Fatal(err)
		return
	}
	// 初始化eval功能
	i := interp.New(interp.Options{})
	err = i.Use(stdlib.Symbols)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		log.Fatal(err)
	}
	// 套件程式碼, 照抄
	_, err = i.Eval(funcString)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		log.Fatal(err)
	} // 逐一將parse函數塞進map
	for _, objid := range UniqueIds {
		v, err := i.Eval("thomas.func" + objid)
		if err != nil {
			sugarLogger.Errorf(err.Error())
			log.Fatal(err)
		}
		EvalFuncHub[objid] = v.Interface().(func([]byte) ([]string, error))
	}
}

func GetParseResult(payload EvalPayload) ([]string, error) {
	if thisFunc, ok := EvalFuncHub[fmt.Sprint(payload.Id)]; ok {
		result, err := thisFunc(payload.Body)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return nil, errors.New("template id " + fmt.Sprint(payload.Id) + " not found")
}

func GetOriAndObjectIDs(payload EvalPayload) (OriginObjectID string, ObjectIDs []string, err error) {
	template, err := queries.GetTemplate(ctx, int32(payload.Id))
	if err != nil {
		return
	}
	// 若ID自動產生則ObjectIDs會回空, 上一層依照取得的response長度來自動產生
	if template.AutoGenObjectID {
		return template.OriginObjectID, ObjectIDs, err
	}
	objectIDs, err := queries.GetObjectIDs(ctx, int32(payload.Id))
	if err != nil {
		return
	}
	for _, v := range objectIDs {
		ObjectIDs = append(ObjectIDs, v.ObjectID)
	}
	return template.OriginObjectID, ObjectIDs, err
}

// 將hot data post給middleware
func InsertSensorDataPost(data InsertSenorDataArray) (resp string, err error) {
	url := fmt.Sprintf("http://%s:9322/api/CCAU/NADI_3DOCMS/CCAU/insertSensorData", Configs.SeverIP)

	method := "POST"
	var payload bytes.Buffer
	err = json.NewEncoder(&payload).Encode(data)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, &payload)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		return
	}
	req.Header.Add("Authorization", Configs.Token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		return
	}
	return string(body), err
}

func EvalFuncString(SingleString string, id int64) (err error) {
	// eval 新函數
	// 初始化eval功能
	i := interp.New(interp.Options{})
	err = i.Use(stdlib.Symbols)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		return
	}
	// 套件程式碼, 照抄
	_, err = i.Eval(SingleString)
	// 解析成程式碼出錯, 代表前端給的值有問題
	if err != nil {
		sugarLogger.Errorf(err.Error())
		return
	}
	// 將parse函數塞進map
	objid := fmt.Sprint(id)
	v, err := i.Eval("thomas.func" + objid)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		return
	}
	EvalFuncHub[objid] = v.Interface().(func([]byte) ([]string, error))
	return
}
