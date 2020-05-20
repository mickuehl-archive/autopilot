package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"shadow-racer/autopilot/v1/pkg/telemetry"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/majordomusio/log15"
)

var (
	broker string
	queue  string
	cl     mqtt.Client
	logger log.Logger
)

func init() {
	logger = log.New("module", "mqtt-local")
}

func shutdownHandler() {
	logger.Info("Shutting down")

	if token := cl.Unsubscribe(queue); token.Wait() && token.Error() != nil {
		logger.Error("Error unsubscribing from queue", "err", token.Error())
		os.Exit(1)
	}
	cl.Disconnect(250)
}

func workerHandler() {
	// FIXME do periodic work
}

func main() {
	flag.StringVar(&broker, "b", "tcp://localhost:1883", "MQTT Broker endpoint")
	flag.StringVar(&queue, "q", "shadow-racer/telemetry", "Default queue for telemetry data")
	flag.Parse()

	// setup shutdown handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		shutdownHandler()
		os.Exit(1)
	}()

	// setup and configuration
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID("mqtt-local")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	// create a client
	cl = mqtt.NewClient(opts)
	if token := cl.Connect(); token.Wait() && token.Error() != nil {
		logger.Error("Error connecting to broker", "err", token.Error())
		os.Exit(1)
	}

	if token := cl.Subscribe(queue, 0, receiveDataFrame); token.Wait() && token.Error() != nil {
		logger.Error("Error subscribing to queue", "err", token.Error())
		os.Exit(1)
	}

	logger.Info("Starting")

	// periodic background processes
	backgroundChannel := time.NewTicker(time.Second * time.Duration(10)).C
	for {
		<-backgroundChannel
		workerHandler()
	}
}

func receiveDataFrame(client mqtt.Client, msg mqtt.Message) {
	var df telemetry.DataFrame

	err := json.Unmarshal(msg.Payload(), &df)
	if err == nil {
		logger.Debug("dataframe", "data", df)
	} else {
		logger.Error("Error unmarshalling a dataframe", "err", err.Error())
	}
}
