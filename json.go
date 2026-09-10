package tools

import (
	"bytes"
	"encoding/json"
	"os"
	"regexp"

	"github.com/basebytes/types"
)

// DecodeFile deserialize file content.
func DecodeFile(file string, v any) error {
	filePtr, err := os.Open(file)
	if err != nil {
		return err
	}
	defer filePtr.Close()
	decoder := json.NewDecoder(filePtr)
	return decoder.Decode(&v)
}

// EncodeObj serialize object to file.
func EncodeObj(v any, file string) error {
	filePtr, err := os.Create(file)
	if err != nil {
		return err
	}
	defer filePtr.Close()
	encoder := json.NewEncoder(filePtr)
	return encoder.Encode(&v)
}

// Decode deserialize
func Decode(jsonStr string, v any) {
	_ = DecodeBytes([]byte(jsonStr), v)
}

// Encode serialize
func Encode(v any) string {
	return string(EncodeBytes(v))
}

// DecodeBytes deserialize
func DecodeBytes(jsonBytes []byte, v any) error {
	return json.Unmarshal(jsonBytes, v)
}

// EncodeBytes serialize
func EncodeBytes(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func DecodeJson5(filePath string, output any) error {
	content, err := os.ReadFile(filePath)
	if err == nil {
		err = json.Unmarshal([]byte(commentReg.ReplaceAllString(string(content), "")), output)
	}
	return err
}

func Copy(src any, dest any) error {
	return DecodeBytes(EncodeBytes(&src), &dest)
}

func TransJson(json *types.Json, result any) {
	if json != nil {
		_ = DecodeBytes(*json, result)
	}
}

func IndentJson(value any, prefix, indent string) string {
	var (
		out  bytes.Buffer
		b, _ = json.Marshal(value)
	)
	_ = json.Indent(&out, b, prefix, indent)
	return out.String()
}

func ValidNumberList(v *types.Json) bool {
	return v == nil || idListReg.Match(*v)
}

func ValidStringList(v *types.Json) bool {
	return v == nil || stringListReg.Match(*v)
}

func ValidJson(v *types.Json) bool {
	return v == nil || json.Valid(*v)
}

var (
	commentReg    = regexp.MustCompile("(//[^\\r\\n]*)")
	idListReg     = regexp.MustCompile(`^\[\]|\[\d+([ \n]*,[ \n]*\d+)*\]$`)
	stringListReg = regexp.MustCompile(`^\[\]|\["(.+)"([ \n]*,[ \n]*".+")*\]$`)
)
