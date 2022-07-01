package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/bluebreezecf/opentsdb-goclient/client"
	"github.com/bluebreezecf/opentsdb-goclient/config"
	consumer "github.com/fogcloud-io/fog-amqp-consumer"
)

var (
	AMQP_HOST         = os.Getenv("AMQP_HOST")
	AMQP_PORT         = os.Getenv("AMQP_PORT")
	FOG_ACCESS_KEY    = os.Getenv("FOG_ACCESS_KEY")
	FOG_ACCESS_SECRET = os.Getenv("FOG_ACCESS_SECRET")
)

var (
	OPENTSDB_HOST = os.Getenv("OPENTSDB_HOST")
	tsdbCli       client.Client
)

func init() {
	amqpCli := consumer.InitAMQPConsumer(context.Background(), AMQP_HOST, AMQP_PORT, FOG_ACCESS_KEY, FOG_ACCESS_SECRET)
	tsdbCli = initTSDBClient()
	go amqpCli.ConsumeWithHanlder(context.Background(), handleMsg)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func initTSDBClient() client.Client {
	cli, err := client.NewClient(config.OpenTSDBConfig{
		OpentsdbHost: OPENTSDB_HOST,
	})
	if err != nil {
		panic(err)
	}
	return cli
}

func handleMsg(msg []byte) {
	data := []client.DataPoint{
		{
			Metric:    "fogcloud",
			Timestamp: time.Now().Unix(),
			Value:     string(msg),
		},
	}
	tsdbCli.Put(data, "details")
}
