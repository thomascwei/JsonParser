	
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
	 type Template10 struct {
 Person struct {
  Name struct {
   First    string `json:"first"`
   Last     string `json:"last"`
   FullName string `json:"fullName"`
  } `json:"name"`
  Github struct {
   Handle    string `json:"handle"`
   Followers int    `json:"followers"`
  } `json:"github"`
  Avatars []struct {
   Ff     int64  `json:"ff"`
   URL    string `json:"url"`
   Type   string `json:"type"`
   Inside []struct {
    Aa struct {
     Zz float64 `json:"zz"`
    } `json:"aa"`
    Bb int `json:"bb"`
   } `json:"inside"`
  } `json:"avatars"`
 } `json:"person"`
 Company struct {
  Name string    `json:"name"`
  List []float64 `json:"list"`
  Strs []string  `json:"strs"`
 } `json:"company"`
}

func func10(ByteData []byte) (result []string, err error) {
		var temp Template10
		// err = json.NewDecoder(bytes.NewReader(ByteData)).Decode(&temp)
		err = json.Unmarshal(ByteData, &temp)
		if err != nil {
			return nil, err
		}
		result = append(result, fmt.Sprint(temp.Person.Name.First))
		return result, nil
	}
