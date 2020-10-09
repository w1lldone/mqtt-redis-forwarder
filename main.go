package main

import (
	"flag"
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	topic := flag.String("topic", "", "The topic name to/from which to publish/subscribe")
	broker := flag.String("broker", "tcp://iot.eclipse.org:1883", "The broker URI. ex: tcp://10.10.1.1:1883")
	password := flag.String("password", "", "The password (optional)")
	user := flag.String("user", "", "The User (optional)")
	id := flag.String("id", "testgoid", "The ClientID (optional)")
	cleansess := flag.Bool("clean", true, "Set Clean Session (default true)")
	qos := flag.Int("qos", 0, "The Quality of Service 0,1,2 (default 0)")
	flag.Parse()

	if *topic == "" {
		fmt.Println("Invalid setting for -topic, must not be empty")
		return
	}
	
	opts := mqtt.NewClientOptions()
	opts.AddBroker(*broker)
	opts.SetClientID(*id)
	opts.SetUsername(*user)
	opts.SetPassword(*password)
	opts.SetCleanSession(*cleansess)

	choke := make(chan [2]string)

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		choke <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe(*topic, byte(*qos), nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for {
		incoming := <-choke
		fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
	}
}
