package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	log "github.com/majordomusio/log15"
)

var (
	logger log.Logger
	proc   *exec.Cmd
)

func init() {
	logger = log.New()
}

func shutdown() {
	logger.Debug("shutdown")
}

func work() {
	logger.Debug("work")
}

func main() {
	// setup shutdown handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		shutdown()
		os.Exit(0)
	}()

	// start the port
	proc = exec.Command("./ports.py")
	stdout, err := proc.StdoutPipe()
	if err != nil {
		logger.Error("Error getting STDOUT", "err", err.Error())
		os.Exit(1)
	}

	if err = proc.Start(); err != nil {
		logger.Error("Error starting the process", "err", err.Error())
		os.Exit(1)
	}

	// read from STDOUT
	go func() {
		var buffer []byte
		buffer = make([]byte, 100)

		for {
			logger.Debug("Waiting for STDOUT")

			n, err := stdout.Read(buffer)
			if err != nil {
				logger.Error("Error reading from STDOUT", "err", err.Error())
			} else {
				logger.Debug("stdout", "n", n, "msg", buffer[:n])
			}
		}
	}()

	// loop until SIGTERM, SIGINT
	backgroundChannel := time.NewTicker(time.Second * time.Duration(10)).C
	for {
		<-backgroundChannel
		work()
	}

	/*if err := proc.Wait(); err != nil {
		logger.Error("Error", "err", err.Error())
		os.Exit(1)
	}*/
}
