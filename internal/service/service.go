package service

import (
	"bufio"
	"fmt"
	notify "github.com/fsnotify/fsnotify"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var elog debug.Log

type Service struct {
}

func (m *Service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}

	path := args[1]
	elog.Info(1, "Get path for work: "+path)

	_, patterns := m.GetPatternsFromPath(path)
	elog.Info(1, "Get patterns for work.")

	watcher, err := notify.NewWatcher()
	if err != nil {
		elog.Error(1, err.Error())
	}
	go m.StartEventsChecker(watcher.Events, patterns)

	watcher.Add(path)
	elog.Info(1, "Watcher started")
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Stop, svc.Shutdown:
				elog.Info(1, "Watcher stopped")
				watcher.Close()
				break loop
			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func (m *Service) GetPatternsFromPath(path string) (string, []string) {
	file, err := os.Open(path + "\\template.tbl")
	if err != nil {
		elog.Error(1, err.Error())
	}
	defer file.Close()

	var fileTextLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fileTextLines = append(fileTextLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		elog.Error(1, err.Error())
	}

	var password = fileTextLines[0]
	var maskTemplates = fileTextLines[1:len(fileTextLines)]
	return password, maskTemplates
}

func (m *Service) StartEventsChecker(events chan notify.Event, patterns []string) {
	var isRenamed = false
	var isByWatcher = false
	var lastName = ""
	var checkPatterns = patterns
	for {
		select {
		case event, ok := <-events:
			if !ok {
				return
			}

			filename := filepath.Base(event.Name)
			if !m.PatternCheck(checkPatterns, filename) {
				return
			}
			elog.Info(1, "Event: "+event.Op.String()+".\n File: "+event.Name)
			switch event.Op {
			case notify.Create:
				if strings.Contains(filename, "копия") {
					for {
						err := os.Remove(event.Name)
						if err == nil {
							break
						}
					}
				} else if isRenamed {
					err := os.Rename(event.Name, lastName)
					if err != nil {
						log.Println("Sorry: ", err.Error())
					}
					isRenamed, isByWatcher = false, true
				} else {
					err := os.Remove(event.Name)
					if err == nil {
						elog.Error(1, err.Error())
					}
				}
			case notify.Rename:
				if !isByWatcher {
					isRenamed = true
					lastName = event.Name
				}
				isByWatcher = false
			}
		}
	}
}

func (m *Service) PatternCheck(patterns []string, filename string) bool {
	for _, pattern := range patterns {
		result, _ := filepath.Match(pattern, filename)
		if result {
			elog.Info(1, "File match pattern: "+pattern)
			return true
		}
	}
	return false
}
