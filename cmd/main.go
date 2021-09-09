package main

import (
	"fmt"
	"github.com/ch4t5ky/armadillo/internal/service"
	svc2 "github.com/ch4t5ky/armadillo/internal/svc"
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
			"       install, remove, debug, start, stop, pause or continue.\n",
		errmsg, os.Args[0])
	os.Exit(2)
}

func main() {
	const svcName = "armadillo"

	inService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running in service: %v", err)
	}

	if inService {
		service.RunService(svcName, false)
		return
	}

	if len(os.Args) < 2 {
		usage("no command specified")
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "debug":
		service.RunService(svcName, true)
		return
	case "install":
		err = svc2.InstallService(svcName, "my service")
	case "remove":
		err = svc2.RemoveService(svcName)
	case "start":
		err = svc2.StartService(svcName)
	case "stop":
		if len(os.Args) < 3 {
			fmt.Println("No password entered")
		}
		err = svc2.ControlService(svcName, svc.Stop, svc.Stopped)
	default:
		usage(fmt.Sprintf("invalid command %s", cmd))
	}
	if err != nil {
		log.Fatalf("failed to %s %s: %v", cmd, svcName, err)
	}
	return
}
