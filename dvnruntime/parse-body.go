package dvnruntime

import "encoding/json"

func ParseBody(body string, val interface{}) error {
	return json.Unmarshal([]byte(body), &val)
}
