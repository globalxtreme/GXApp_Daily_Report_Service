package config

import (
	"fmt"
	xtremedb "github.com/globalxtreme/go-core/v2/database"
	xtremerabbitmq "github.com/globalxtreme/go-core/v2/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"os"
	"time"
)

func InitRabbitMQ() func() {
	var rabbitMQClose func()

	xtremerabbitmq.RabbitMQSQL, rabbitMQClose = xtremedb.Connect(xtremedb.DBConf{
		Driver:    xtremedb.MYSQL_DRIVER,
		Host:      os.Getenv("DB_RABBITMQ_HOST"),
		Port:      os.Getenv("DB_RABBITMQ_PORT"),
		Username:  os.Getenv("DB_RABBITMQ_USERNAME"),
		Password:  os.Getenv("DB_RABBITMQ_PASSWORD"),
		Database:  os.Getenv("DB_RABBITMQ_DATABASE"),
		ParseTime: true,
	})

	xtremerabbitmq.RabbitMQConf.Connection = make(map[string]xtremerabbitmq.RabbitMQConnectionConf, 2)

	globalHost := os.Getenv("RABBITMQ_GLOBAL_HOST")
	if globalHost != "" {
		xtremerabbitmq.RabbitMQConf.Connection[xtremerabbitmq.RABBITMQ_CONNECTION_GLOBAL] = xtremerabbitmq.RabbitMQConnectionConf{
			Host:     globalHost,
			Port:     os.Getenv("RABBITMQ_GLOBAL_PORT"),
			Username: os.Getenv("RABBITMQ_GLOBAL_USER"),
			Password: os.Getenv("RABBITMQ_GLOBAL_PASSWORD"),
		}
	}

	localHost := os.Getenv("RABBITMQ_LOCAL_HOST")
	if localHost != "" {
		xtremerabbitmq.RabbitMQConf.Connection[xtremerabbitmq.RABBITMQ_CONNECTION_LOCAL] = xtremerabbitmq.RabbitMQConnectionConf{
			Host:     localHost,
			Port:     os.Getenv("RABBITMQ_LOCAL_PORT"),
			Username: os.Getenv("RABBITMQ_LOCAL_USER"),
			Password: os.Getenv("RABBITMQ_LOCAL_PASSWORD"),
		}
	}

	xtremerabbitmq.RabbitMQConf.Timeout = 5 * time.Second

	return rabbitMQClose
}

func InitRabbitMQConnection() func() {
	xtremerabbitmq.RabbitMQConnectionDial = make(map[string]*amqp091.Connection, 2)
	for connectionName, connectionConf := range xtremerabbitmq.RabbitMQConf.Connection {
		conn, err := amqp091.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
			connectionConf.Username, connectionConf.Password, connectionConf.Host, connectionConf.Port))
		if err != nil {
			log.Panicf("Failed to connect to RabbitMQ: %s", err)
		}

		xtremerabbitmq.RabbitMQConnectionDial[connectionName] = conn
	}

	rabbitMQConnClose := func() {
		for _, connection := range xtremerabbitmq.RabbitMQConnectionDial {
			connection.Close()
		}
	}

	return rabbitMQConnClose
}
