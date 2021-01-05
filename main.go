package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// Message struc
type Message struct {
	Reconnecting bool `json:"reconnecting"`
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	redisServer := viper.GetString("redis.server")
	redisPassword := viper.GetString("redis.password")
	redisDb := viper.GetInt("redis.db")
	redisChannel := viper.GetString("redis.channel")

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisServer,
		Password: redisPassword,
		DB:       redisDb,
	})

	pubsub := rdb.Subscribe(ctx, redisChannel)
	ch := pubsub.Channel()

	topic := viper.GetString("mqtt.topic")
	broker := viper.GetString("mqtt.broker")
	password := viper.GetString("mqtt.password")
	user := viper.GetString("mqtt.user")
	id := viper.GetString("mqtt.id")
	cleansess := viper.GetBool("mqtt.cleaness")
	qos := viper.GetInt("mqtt.qos")

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(id)
	opts.SetUsername(user)
	opts.SetPassword(password)
	opts.SetCleanSession(cleansess)

	choke := make(chan [2]string)

	handler := func(client mqtt.Client, msg mqtt.Message) {
		choke <- [2]string{msg.Topic(), string(msg.Payload())}
	}

	opts.SetOnConnectHandler(func(c mqtt.Client) {
		fmt.Println("connection established")
		if token := c.Subscribe(topic, byte(qos), handler); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		select {
		case incoming := <-choke:
			fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])

			err := rdb.Publish(ctx, redisChannel, incoming[1]).Err()
			if err != nil {
				panic(err)
			}
		case msg := <-ch:
			if subscriberIsDisconnected(msg.Payload) {
				os.Exit(0)
			}
		}
	}
}

func subscriberIsDisconnected(jsonMessage string) bool {
	var message Message
	err := json.Unmarshal([]byte(jsonMessage), &message)
	if err != nil {
		log.Println(err)
	}

	if message.Reconnecting == true {
		fmt.Printf("Subscriber reconnected \n")
		return true
	}

	return false
}
