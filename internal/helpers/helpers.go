package helpers

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

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

func updateValueInFile(password string, templates []string) {
	hashPassword := createMD5Hash(password)
	fmt.Println(hashPassword)

	file, err := os.OpenFile("template.tbl", os.O_RDWR, 644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.WriteString(hashPassword + "\n")
	for _, value := range templates {
		_, err = file.WriteString(value + "\n")
	}
	err = file.Sync()
}

func LockFiles() {
	password, templates := GetValuesFromTemplate("test")
	for _, value := range templates {
		fmt.Println("Template: " + value)
	}

	fmt.Println("Enter new password: ")
	fmt.Scanf("%s\n", &password)

	updateValueInFile(password, templates)
}

func createMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
