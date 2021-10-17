package helpers

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"golang.org/x/sys/windows/svc/debug"
	"log"
	"os"
)

var elog debug.Log

func GetValuesFromTemplate(path string) (string, []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var fileTextLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileTextLines = append(fileTextLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var password = fileTextLines[0]
	var maskTemplates = fileTextLines[1:len(fileTextLines)]

	return password, maskTemplates
}

func UpdatePasswordInFile(path string, password string) {
	fileName := path + "\\template.tbl"
	hashPassword := CreateMD5Hash(password)
	_, patterns := GetValuesFromTemplate(fileName)
	file, err := os.OpenFile(fileName, os.O_RDWR, 644)
	if err != nil {
		fmt.Println(1, err.Error())
	}
	defer file.Close()

	_, err = file.WriteString(hashPassword + "\n")
	for _, value := range patterns {
		_, err = file.WriteString(value + "\n")
	}
	err = file.Sync()
}

func CreateMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
