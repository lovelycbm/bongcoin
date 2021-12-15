package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

func HandleError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	HandleError(encoder.Encode(i))
	return aBuffer.Bytes()
}

func FromBytes( i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	HandleError(decoder.Decode(i))
}

func Hash(i interface{}) string {
	s:= fmt.Sprint(i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func Splitter(s string, sep string, i int) string {
	r := strings.Split(s, sep)
	
	if len(r)-1 < i {
		return ""
	}
	return r[i]
}

func ToJSON(i interface{}) []byte {
	b, err := json.Marshal(i)
	HandleError(err)
	return b
}