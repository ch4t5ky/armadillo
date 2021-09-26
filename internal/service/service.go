package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
	"net/http"
)

var elog debug.Log

type Service struct {
	http *http.Server
}

func New(addr string) *Service {
	var service = &Service{
		http: &http.Server{
			Addr: addr,
		},
	}
	service.http.Handler = service.setupRouter()
	return service
}

func (s *Service) setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/api/lock", s.testProtect)
	r.POST("/api/unlock")
	return r
}

func (s *Service) Start() error {
	elog.Info(1, "Start REST API")
	return s.http.ListenAndServe()
}

func (s *Service) Stop() error {
	elog.Info(1, "Stop REST API")
	return s.http.Close()
}

func (s *Service) testProtect(c *gin.Context) {
	c.JSON(http.StatusOK, "nice")

}
func (m *Service) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}

	errs := make(chan error, 1)
	go func() {
		errs <- m.Start()
	}()

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Continue:
				elog.Info(1, "Running protection.")
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			case svc.Pause:
				elog.Info(1, "Stop protection.")
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Stop, svc.Shutdown:
				_ = m.Stop()
				m.Stop()
				break loop
			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func RunService(name string, api_addr string) {
	var err error
	elog, err = eventlog.Open(name)
	if err != nil {
		return
	}
	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name))
	run := svc.Run
	service := New(api_addr)
	err = run(
		name,
		service,
	)
	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	elog.Info(1, fmt.Sprintf("%s service stopped", name))
}
