package tools

import (
	"encoding/json"
	"os"
	"regexp"
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

var commentReg = regexp.MustCompile("(//[^\\r\\n]*)")
