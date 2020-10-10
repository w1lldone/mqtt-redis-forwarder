package main

import (
	"context"
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

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

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		choke <- [2]string{msg.Topic(), string(msg.Payload())}
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe(topic, byte(qos), nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for {
		incoming := <-choke
		fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
		err := rdb.Publish(ctx, redisChannel, incoming[1]).Err()
		if err != nil {
			panic(err)
		}
	}
}
