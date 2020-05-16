package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/abiosoft/ishell"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	frequency = 50
)

type board struct {
	adaptor        *raspi.Adaptor
	driver         *i2c.PCA9685Driver
	channels       []pulse
	defaultChannel int
}

type pulse struct {
	min int
	max int
}

var thing *board

// set the PWM for a channel
func setFunc(c *ishell.Context) {
	if len(c.Args) >= 2 {
		ch := thing.defaultChannel
		min := 0
		max := 0
		if len(c.Args) == 2 {
			min, _ = strconv.Atoi(c.Args[0])
			max, _ = strconv.Atoi(c.Args[1])
		} else {
			ch, _ = strconv.Atoi(c.Args[0])
			min, _ = strconv.Atoi(c.Args[1])
			max, _ = strconv.Atoi(c.Args[2])
		}
		if ch < 0 || ch > 15 {
			c.Println("set: invalid arguments")
			return
		}
		if min < 0 || min > 4096 {
			c.Println("set: invalid arguments")
			return
		}
		if max < 0 || max > 4096 {
			c.Println("set: invalid arguments")
			return
		}

		err := setPWM(ch, min, max)
		if err != nil {
			c.Println("set: error setting parameters")
		}
	} else {
		c.Println("set: not enough arguments")
	}
}

// turn the PWM for a channel OFF i.e. pulse = 0,0
func offFunc(c *ishell.Context) {
	s := ""
	if len(c.Args) == 1 {
		ch := getChannel(c, 1)
		setPWM(ch, 0, 0)
		s = fmt.Sprintf("Set channel %d to OFF", ch)
	} else {
		for i := 0; i < 16; i++ {
			setPWM(i, 0, 0)
		}
		s = "Set all channels to OFF"
	}
	c.Println(s)
}

// shows the current values of a channel
func showFunc(c *ishell.Context) {
	ch := getChannel(c, 1)
	s := fmt.Sprintf("channel %d: %d,%d", ch, thing.channels[ch].min, thing.channels[ch].max)
	c.Println(s)
}

// sets the default channel
func channelFunc(c *ishell.Context) {
	thing.defaultChannel = getChannel(c, 1)
	s := fmt.Sprintf("Default channel: %d", thing.defaultChannel)
	c.Println(s)
}

func catchC(c *ishell.Context, count int, input string) {
	c.Println("Ctrl-C, exiting now ...")
	cleanup()
	os.Exit(0)
}

// returns the channel number from pos. pos is 1 based !!!!
func getChannel(c *ishell.Context, pos int) int {
	if len(c.Args) == 0 {
		return thing.defaultChannel
	}
	if pos > len(c.Args) {
		return thing.defaultChannel
	}
	ch, err := strconv.Atoi(c.Args[pos-1])
	if err != nil {
		return thing.defaultChannel
	}
	if ch < 0 {
		return thing.defaultChannel
	}
	if ch > 15 {
		return thing.defaultChannel
	}
	return ch
}

func initialize() {
	a := raspi.NewAdaptor()
	d := i2c.NewPCA9685Driver(a)

	t := &board{
		adaptor:        a,
		driver:         d,
		channels:       make([]pulse, 16),
		defaultChannel: 0,
	}

	t.driver.Start()
	time.Sleep(1 * time.Second)
	t.driver.SetPWMFreq(frequency)

	thing = t
}

func cleanup() {
	thing.driver.Halt()
	thing.driver.Connection().Finalize()
	time.Sleep(1 * time.Second)
}

func setPWM(ch, min, max int) error {
	thing.channels[ch].min = min
	thing.channels[ch].max = max

	return thing.driver.SetPWM(ch, uint16(min), uint16(max))
}

func main() {
	shell := ishell.New()
	shell.Println("PWM calibration")

	// initialize the hardware
	initialize()

	// register custom exit handlers
	shell.Interrupt(catchC)
	shell.DeleteCmd("exit")

	// register the commands
	shell.AddCmd(&ishell.Cmd{
		Name: "set",
		Help: "Set the PWM values on the selected channel",
		Func: setFunc,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "channel",
		Help: "Select a channel (0..15)",
		Func: channelFunc,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "show",
		Help: "Shows the PWM values on the selected channel",
		Func: showFunc,
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "off",
		Help: "Turn one or all channels OFF",
		Func: offFunc,
	})

	// run the CLI
	shell.Run()
}
