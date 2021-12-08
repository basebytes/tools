package tools

import (
	"encoding/json"
	"os"
)

//deserialize file content.
func DecodeFile(file string, v interface{}) error {
	filePtr, err := os.Open(file)
	if err != nil {
		return err
	}
	defer filePtr.Close()
	decoder := json.NewDecoder(filePtr)
	return decoder.Decode(&v)
}

//serialize object to file.
func EncodeObj(v interface{}, file string) error {
	filePtr, err := os.Create(file)
	if err != nil {
		return err
	}
	defer filePtr.Close()
	encoder := json.NewEncoder(filePtr)
	return encoder.Encode(&v)
}

//deserialize
func Decode(jsonStr string, v interface{}) {
	_ = DecodeBytes([]byte(jsonStr), v)
}

//serialize
func Encode(v interface{}) string {
	return string(EncodeBytes(v))
}

//deserialize
func DecodeBytes(jsonBytes []byte, v interface{}) error {
	return json.Unmarshal(jsonBytes, v)
}

//serialize
func EncodeBytes(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
