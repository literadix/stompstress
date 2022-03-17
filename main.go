package main

import (
	"flag"
	"fmt"
	"github.com/go-stomp/stomp/v3"
	"go.uber.org/ratelimit"
	"math/rand"
	"os"
	"strings"
)

const defaultPort = ":61613"

var letters = []rune("ABCDEF1234567890")

var serverAddr = flag.String("server", "localhost:61613", "STOMP server endpoint")
var messageCount = flag.Int("count", 10, "Number of messages to send")
var messageSize = flag.Int("size", 1000, "Size of each message")
var rateLimit = flag.Int("rate", 100, "Messages per second")

var queueName = flag.String("queue", "/queue/client_test", "Destination queue")
var helpFlag = flag.Bool("help", false, "Print help text")

var stop = make(chan bool)
var sentAll = make(chan bool)

// these are the default options that work with RabbitMQ
var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
	stomp.ConnOpt.Login("guest", "guest"),
	stomp.ConnOpt.Host("/"),
}

func main() {

	flag.Parse()
	if *helpFlag {
		fmt.Fprintf(os.Stderr, "Usage of %s\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	go sendMessages()
	<-sentAll
	<-stop

}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		if i%100 == 0 {
			b[i] = '\n'
		} else {
			b[i] = letters[rand.Intn(len(letters))]
		}

	}
	return strings.TrimPrefix(string(b), "\n")
}

func sendMessages() {
	defer func() {
		stop <- true
	}()

	conn, err := stomp.Dial("tcp", *serverAddr, options...)
	if err != nil {
		println("cannot connect to server", err.Error())
		return
	}

	rl := ratelimit.New(*rateLimit) // messages per second

	println("Sending messages ...")
	for i := 1; i <= *messageCount; i++ {

		text := randSeq(*messageSize)

		if i%1000 == 0 {
			fmt.Printf("Sent messages: %d\n", i)
		}

		rl.Take()

		tx := conn.Begin()

		err = tx.Send(*queueName, "text/plain",
			[]byte(text), stomp.SendOpt.Receipt)
		if err != nil {
			println("failed to send to server, will disconnect", err)
			err = conn.Disconnect()
			sentAll <- true
			return
		}

		err = tx.Commit()
		if err != nil {
			return
		}

	}

	println("disconnecting ...")
	err = conn.Disconnect()

	if err != nil {
		println("failed to disconnect", err)
		return
	} else {
		println("disconnected")
	}

	sentAll <- true
}
