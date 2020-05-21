package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"shadow-racer/autopilot/v1/pkg/telemetry"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/majordomusio/commons/pkg/util"
	log "github.com/majordomusio/log15"
)

var (
	broker     string
	queue      string
	cl         mqtt.Client
	logger     log.Logger
	currentDir string
	dumpFile   *os.File
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

	// close the dumpfile
	if dumpFile != nil {
		dumpFile.Sync()
		dumpFile.Close()
	}
}

func workerHandler() {
	// FIXME do periodic work
}

//
// Structure of the CSV file:
//
// TS, Batch, N, DeviceID, RecordingTS, TH, ST, HEAD
//

func dataFrameToCSVString(df *telemetry.DataFrame) string {
	return fmt.Sprintf("%d,%d,%d,%s,%s,%s,%s,%s\n", df.TS, df.Batch, df.N, df.DeviceID, df.Data["recording_ts"], df.Data["th"], df.Data["st"], df.Data["head"])
}

func receiveDataFrame(client mqtt.Client, msg mqtt.Message) {
	var df telemetry.DataFrame
	err := json.Unmarshal(msg.Payload(), &df)
	if err == nil {
		if df.Type == telemetry.KV {
			_, err := dumpFile.WriteString(dataFrameToCSVString(&df))
			if err != nil {
				logger.Error("Error dumping data to file", "err", err.Error())
			}
		} else {
			if len(df.Blob) != 0 {
				blob, err := base64.StdEncoding.DecodeString(df.Blob)
				if err != nil {
					logger.Error("Error unmarshalling a blob", "err", err.Error())
				} else {
					fn := fmt.Sprintf("%s/%d_%d.jpg", currentDir, df.Batch, df.N)
					err := ioutil.WriteFile(fn, blob, 0644)
					if err != nil {
						logger.Error("Error dumping blob to file", "file", fn, "err", err.Error())
					}
				}
			}
		}
	} else {
		logger.Error("Error unmarshalling a dataframe", "err", err.Error())
	}
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

	// prepare the datadump location
	currentDir = fmt.Sprintf("dump/%d", util.Timestamp())
	err := os.MkdirAll(currentDir, 0755)
	if err != nil {
		logger.Error("Error creating the data directory", "dir", currentDir, "err", err.Error())
		os.Exit(1)
	}

	// datadump file
	dumpFile, err = os.Create(currentDir + "/data.csv")
	if err != nil {
		logger.Error("Error creating the data file", "file", currentDir+"/data.csv", "err", err.Error())
		os.Exit(1)
	}

	// MQTT setup and configuration
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID("mqtt-local")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	// create a MQTT client and subscribe to the queue
	cl = mqtt.NewClient(opts)
	if token := cl.Connect(); token.Wait() && token.Error() != nil {
		logger.Error("Error connecting to broker", "err", token.Error())
		os.Exit(1)
	}

	if token := cl.Subscribe(queue, 0, receiveDataFrame); token.Wait() && token.Error() != nil {
		logger.Error("Error subscribing to queue", "err", token.Error())
		os.Exit(1)
	}

	logger.Info("Ready to receive data", "queue", queue, "broker", broker)

	// periodic background processes
	backgroundChannel := time.NewTicker(time.Second * time.Duration(10)).C
	for {
		<-backgroundChannel
		workerHandler()
	}
}
