package main

import (
	"fmt"
	"github.com/ch4t5ky/armadillo/internal/helpers"
	"github.com/ch4t5ky/armadillo/internal/service"
	"github.com/ch4t5ky/armadillo/internal/windows"
	"golang.org/x/sys/windows/svc"
	"log"
	"os"
	"strings"
)

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <command>\n"+
			"       where <command> is one of\n"+
			"       install, remove, start, stop \n",
		errmsg, os.Args[0])
	os.Exit(1)
}

func main() {
	const svcName = "armadillo"

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	inService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running in service: %v", err)
	}

	if inService {
		service.RunService(svcName)
		return
	}

	if len(os.Args) < 2 {
		usage("no command specified")
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "install":
		err = service.InstallService(svcName, "Service for protection directory")
	case "remove":
		err = service.RemoveService(svcName)
	case "start":
		var password string
		fmt.Println("Enter password: ")
		fmt.Scanf("%s\n", password)
		helpers.UpdatePasswordInFile(path, password)
		err = service.StartService(svcName, path)
	case "stop":
		path, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		windows.Chmod(path+"\\template.tbl", 0777)
		password, _ := helpers.GetValuesFromTemplate(path + "\\template.tbl")

		var enteredPassword string
		fmt.Println("Enter password for file: ")
		fmt.Scanf("%s\n", enteredPassword)

		if helpers.CreateMD5Hash(enteredPassword) != password {
			fmt.Println("Incorrect Password: Protection continue")
			windows.Chmod(path+"\\template.tbl", 0000)
		} else {
			fmt.Println("Correct Password: Protection stopped")
			err = service.ControlService(svcName, svc.Stop, svc.Stopped)
		}

	default:
		usage(fmt.Sprintf("invalid command %s", cmd))
	}
	if err != nil {
		log.Fatalf("failed to %s %s: %v", cmd, svcName, err)
	}
	return
}
