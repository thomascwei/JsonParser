// 用於動態產生所有取json response結果的函數
// 產生解析及取值函數, 該函數response為slice,前端同時提供取值邏輯

package internal

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
)

func GenerateAllParseFuncString() (PackageBase string, err error) {
	PackageBase = `	
// Code generated by EvalFunction. Just for eyeball check
package thomas
import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var (
	Trace     = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	// StructMap = make(map[string]func([]byte) ([]string, error))
)
	`
	FuncBaseReturnSingle := `

func func%d(ByteData []byte) (result []string, err error) {
		var temp Template%d
		// err = json.NewDecoder(bytes.NewReader(ByteData)).Decode(&temp)
		err = json.Unmarshal(ByteData, &temp)
		if err != nil {
			return nil, err
		}
		result = append(result, fmt.Sprint(temp.%s))
		return result, nil
	}
`
	FuncBaseReturnArray := `

func func%d(ByteData []byte) (result []string, err error) {
		var temp Template%d
		// err = json.NewDecoder(bytes.NewReader(ByteData)).Decode(&temp)
		err = json.Unmarshal(ByteData, &temp)
		if err != nil {
			return nil, err
		}
		for _,v:= range temp.%s{
			result = append(result, fmt.Sprint(v))
		}
		return result, nil
	}
`
	FuncBaseComplexReturnArray := `

func func%d(ByteData []byte) (result []string, err error) {
		var temp Template%d
		// err = json.NewDecoder(bytes.NewReader(ByteData)).Decode(&temp)
		err = json.Unmarshal(ByteData, &temp)
		if err != nil {
			return nil, err
		}
		for _,v:= range temp.%s{
			result = append(result, fmt.Sprint(v%s))
		}
		return result, nil
	}
`
	results, err := queries.ListTemplatesForParse(ctx)
	if err != nil {
		sugarLogger.Errorf(err.Error())
		return "", err
	}
	for _, v := range results {
		UniqueIds = append(UniqueIds, fmt.Sprint(v.ID))
		PackageBase += " "
		PackageBase += strings.Replace(v.GoStruct, "AutoGenerated", fmt.Sprintf("Template%d", v.ID), -1)
		switch v.ParseType {
		case 0:
			PackageBase += fmt.Sprintf(FuncBaseReturnSingle, v.ID, v.ID, v.ValueExtract)
		case 1:
			PackageBase += fmt.Sprintf(FuncBaseReturnArray, v.ID, v.ID, v.ValueExtract)
		case 2:
			// ex:Person.Avatars[].Ff
			parts := strings.Split(v.ValueExtract, "[]")
			if len(parts) != 2 {
				return "", errors.New("ValueExtract format error")
			}
			PackageBase += fmt.Sprintf(FuncBaseComplexReturnArray, v.ID, v.ID, parts[0], parts[1])
		}

	}

	os.WriteFile(path.Join(rootPath, "temp", "temp.go"), []byte(PackageBase), 0644)
	return PackageBase, nil
}

func GenerateNewAddParseFuncString(id int64, template RequestTemplateOne) (PackageBase string, err error) {
	PackageBase = `	
// Code generated by EvalFunction. Just for eyeball check
package thomas
import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var (
	Trace     = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	// StructMap = make(map[string]func([]byte) ([]string, error))
)
	`
	FuncBaseReturnSingle := `

func func%d(ByteData []byte) (result []string, err error) {
		var temp Template%d
		// err = json.NewDecoder(bytes.NewReader(ByteData)).Decode(&temp)
		err = json.Unmarshal(ByteData, &temp)
		if err != nil {
			return nil, err
		}
		result = append(result, fmt.Sprint(temp.%s))
		return result, nil
	}
`
	FuncBaseReturnArray := `

func func%d(ByteData []byte) (result []string, err error) {
		var temp Template%d
		// err = json.NewDecoder(bytes.NewReader(ByteData)).Decode(&temp)
		err = json.Unmarshal(ByteData, &temp)
		if err != nil {
			return nil, err
		}
		for _,v:= range temp.%s{
			result = append(result, fmt.Sprint(v))
		}
		return result, nil
	}
`
	FuncBaseComplexReturnArray := `

func func%d(ByteData []byte) (result []string, err error) {
		var temp Template%d
		// err = json.NewDecoder(bytes.NewReader(ByteData)).Decode(&temp)
		err = json.Unmarshal(ByteData, &temp)
		if err != nil {
			return nil, err
		}
		for _,v:= range temp.%s{
			result = append(result, fmt.Sprint(v%s))
		}
		return result, nil
	}
`

	PackageBase += " "
	PackageBase += strings.Replace(template.GoStruct, "AutoGenerated", fmt.Sprintf("Template%d", id), -1)
	switch template.ParseType {
	case 0:
		PackageBase += fmt.Sprintf(FuncBaseReturnSingle, id, id, template.ValueExtract)
	case 1:
		PackageBase += fmt.Sprintf(FuncBaseReturnArray, id, id, template.ValueExtract)
	case 2:
		// ex:Person.Avatars[].Ff
		parts := strings.Split(template.ValueExtract, "[]")
		if len(parts) != 2 {
			return "", errors.New("ValueExtract format error")
		}
		PackageBase += fmt.Sprintf(FuncBaseComplexReturnArray, id, id, parts[0], parts[1])
	}

	os.WriteFile(path.Join(rootPath, "temp", "new.go"), []byte(PackageBase), 0644)

	return PackageBase, nil
}