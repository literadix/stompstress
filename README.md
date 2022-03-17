# stompstress
simple stomp stress tool

$ ./stomp_stress --help

Usage of ./stomp_stress
  -count int
        Number of messages to send (default 10)
  -help
        Print help text
  -queue string
        Destination queue (default "/queue/client_test")
  -rate int
        Messages per second (default 100)
  -server string
        STOMP server endpoint (default "localhost:61613")
  -size int
        Size of each message (default 1000)

$ ./stomp_stress -count 10000 -rate 200 -size 200000

Sending messages ...
Sent messages: 1000
Sent messages: 2000
Sent messages: 3000
Sent messages: 4000
Sent messages: 5000
Sent messages: 6000
Sent messages: 7000
Sent messages: 8000
Sent messages: 9000
Sent messages: 10000
disconnecting ...
disconnected
