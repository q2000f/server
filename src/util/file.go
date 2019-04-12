package util

import (
	"encoding/json"
	"github.com/sipt/GoJsoner"
	"go/format"
	"io/ioutil"
	"os"
)

func LoadJsonFile(fileName string, out interface{}) error {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	txt, err := GoJsoner.Discard(string(bytes))
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(txt), out)
	if err != nil {
		return err
	}
	return nil
}

func SaveFile(dir, fileName string, bytes []byte) error {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dir+"/"+fileName, bytes, os.ModePerm)
}

func SaveGoFile(dir, fileName string, bytes []byte) error {
	fmtBytes, err := format.Source(bytes)
	if err != nil {
		return err
	}
	return SaveFile(dir, fileName, fmtBytes)
}
