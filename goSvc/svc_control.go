//go:build windows

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

type myservice struct{}


func (m *myservice) Execute(
	args []string, 
	r <-chan svc.ChangeRequest, 
	changes chan<- svc.Status,
) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	serve_it()
	return
}

func exePath() (string, error) {
	var err error
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is directory", p)
	}
	if filepath.Ext(p) == "" {
		var fi os.FileInfo

		p += ".exe"
		fi, err = os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			err = fmt.Errorf("%s is directory", p)
		}
	}
	return "", err
}

func runService() {
	run := svc.Run
	err := run(svcName, &myservice{})
	if err != nil {
		log.Fatalf("Could not run service: %s", svcName)
	}
}

func installService() {
	exepath, err := exePath()
	if err != nil {
		log.Fatal(err)
	}
	m, err := mgr.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer m.Disconnect()
	s, err := m.OpenService(svcName)
	if err == nil {
		s.Close()
		log.Fatalf("service %s already exists", svcName)
	}
	s, err = m.CreateService(
		svcName, 
		exepath, 
		mgr.Config{DisplayName: svcName}, 
		"is", 
		"auto-started",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	if err != nil {
		s.Delete()
		log.Fatal(err)
	}
}

func removeService() {
	m, err := mgr.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer m.Disconnect()
	s, err := m.OpenService(svcName)
	if err != nil {
		log.Fatalf("service %s is not installed", svcName)
	}
	defer s.Close()
	err = s.Delete()
	if err != nil {
		log.Fatal(err)
	}
}

func startService() {
	m, err := mgr.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer m.Disconnect()
	s, err := m.OpenService(svcName)
	if err != nil {
		log.Fatalf("could not access service: %v", err)
	}
	defer s.Close()
	err = s.Start("is", "manual-started")
	if err != nil {
		log.Fatalf("could not start service: %v", err)
	}
}

func stopService() {
	m, err := mgr.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer m.Disconnect()
	s, err := m.OpenService(svcName)
	if err != nil {
		log.Fatalf("could not access service: %v", err)
	}
	defer s.Close()
	status, err := s.Control(svc.Stop)
	if err != nil {
		log.Fatalf("could not send control=%d: %v", svc.Stop, err)
	}
	timeout := time.Now().Add(10 * time.Second)
	for status.State != svc.Stopped {
		if timeout.Before(time.Now()) {
			log.Fatalf("timeout waiting for service to go to state=%d", svc.Stopped)
		}
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			log.Fatalf("could not retrieve service status: %v", err)
		}
	}
}
