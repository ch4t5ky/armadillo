package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
)

func GetFilesInDirectory() {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fmt.Println(f.Name())
	}
}