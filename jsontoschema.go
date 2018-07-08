package jsontoschema

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type JSONType string

const (
	JSONTypeString  = "string"
	JSONTypeInteger = "integer"
	JSONTypeBool    = "boolean"
	JSONTypeNumber  = "number"
	JSONTypeArray   = "array"
	JSONTypeObject  = "object"
)

type JSONArrayValidation struct {
	Type JSONType `json:"type"`
}

type JSONProperty struct {
	Type       JSONType                `json:"type"`
	Items      *JSONArrayValidation    `json:"items,omitempty"`      // Populated if item is an array
	Properties map[string]JSONProperty `json:"properties,omitempty"` // Used if property type is an object
}

type JSONSchema struct {
	Schema      string                  `json:"$schema"`
	ID          string                  `json:"$id"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Type        JSONType                `json:"type"`
	Properties  map[string]JSONProperty `json:"properties"`
}

func JsonToSchema(jsonStr string) (string, error) {
	m, err := jsonToMap(strings.NewReader(jsonStr))
	if err != nil {
		return "", err
	}

	js := JSONSchema{Type: JSONTypeObject}

	js.Properties = iterMap(m)

	var b []byte
	b, err = json.MarshalIndent(js, "", "\t")

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func iterMap(jsonMap map[string]interface{}) map[string]JSONProperty {
	props := make(map[string]JSONProperty)
	for k, v := range jsonMap {
		switch v.(type) {
		case string:
			prop := JSONProperty{Type: JSONTypeString}
			props[k] = prop
		case float64:
			prop := JSONProperty{Type: JSONTypeNumber}
			props[k] = prop
		case bool:
			prop := JSONProperty{Type: JSONTypeBool}
			props[k] = prop
		case map[string]interface{}:
			vmap := v.(map[string]interface{})
			prop := JSONProperty{Type: JSONTypeObject}
			prop.Properties = iterMap(vmap)
			props[k] = prop
		case []interface{}:
			prop := JSONProperty{Type: JSONTypeArray}
			arr := v.([]interface{})
			if len(arr) > 0 {
				prop.Items = &JSONArrayValidation{}
				prop.Items.Type = getJSONType(arr[0])
			}
			props[k] = prop
		}
	}

	return props
}

func getJSONType(v interface{}) JSONType {
	switch v.(type) {
	case string:
		return JSONTypeString
	case float64:
		return JSONTypeNumber
	case bool:
		return JSONTypeBool
	case map[string]interface{}:
		return JSONTypeObject
	case []interface{}:
		return JSONTypeArray
	}

	return JSONTypeString
}

func jsonToMap(r io.Reader) (map[string]interface{}, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	jsonStr := buf.String()

	var i interface{}

	err := json.Unmarshal([]byte(jsonStr), &i)
	if err != nil {
		return nil, err
	}

	switch i.(type) {
	case map[string]interface{}:
		return i.(map[string]interface{}), nil

	case []interface{}:
		l := i.([]interface{})
		if len(l) > 0 {
			v := l[0]
			if _, ok := v.(map[string]interface{}); !ok {
				return nil, fmt.Errorf("unable to cast json string %s to map[string]interface{} og []interface{}", jsonStr)
			}
			return v.(map[string]interface{}), nil
		}
		return map[string]interface{}{}, nil

	default:
		return nil, fmt.Errorf("unable to cast json string %s to map[string]interface{} og []interface{}", jsonStr)
	}
}
