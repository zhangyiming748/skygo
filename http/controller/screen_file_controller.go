package controller

import (
	"io/ioutil"
)

type MyFileContreller struct{}

func (this *MyFileContreller) Open(name string) []byte {
	result, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return result
}
