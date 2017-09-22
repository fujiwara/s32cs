package s32cs

import (
	"errors"
	"fmt"
	"regexp"
)

// http://docs.aws.amazon.com/cloudsearch/latest/developerguide/preparing-data.html
var InvalidChars = regexp.MustCompile("[^\u0009\u000a\u000d\u0020-\uD7FF\uE000-\uFFFD]")

type SDFRecord struct {
	ID     string                 `json:"id"`
	Type   string                 `json:"type"`
	Fields map[string]interface{} `json:"fields,omitempty"`
}

func (r *SDFRecord) Validate() error {
	if r.ID == "" {
		return errors.New("id not defined")
	}
	if r.Type != "add" && r.Type != "delete" {
		return fmt.Errorf("type %s is not allowed", r.Type)
	}
	for key, _value := range r.Fields {
		switch value := _value.(type) {
		case string:
			r.Fields[key] = InvalidChars.ReplaceAllString(value, "")
		case []interface{}:
			for i, _v := range value {
				switch v := _v.(type) {
				case string:
					value[i] = InvalidChars.ReplaceAllString(v, "")
				}
			}
			r.Fields[key] = value
		}
	}
	return nil
}
