package lmstfy

import (
	"falcon-mail-transmit/interval/config"
	"falcon-mail-transmit/lib/log"
	"os"
	"sync"
)
import "github.com/bitleak/lmstfy/client"

var once sync.Once
var lmstfy *client.LmstfyClient

type QueueLmstfyClient struct {
	Client *client.LmstfyClient
}

func GetInstance() *QueueLmstfyClient {
	once.Do(func() {
		lmstfy = create()
	})
	return &QueueLmstfyClient{Client: lmstfy}
}

func create() *client.LmstfyClient {
	cfg, err := config.Load()
	if err != nil {
		log.Logger.Fatal("failed to load application configuration: ", err.Error())
		os.Exit(-1)
	}
	c := client.NewLmstfyClient(cfg.Lmstfy.Host, cfg.Lmstfy.Port, cfg.Lmstfy.Namespace, cfg.Lmstfy.Token)

	return c
}
