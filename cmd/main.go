package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tele "gopkg.in/telebot.v3"

	"github.com/vladimish/mdb/controller"
	"github.com/vladimish/mdb/internal/python_downloader"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	ctx := context.Background()

	pref := tele.Settings{
		Token:  os.Getenv("TG_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return fmt.Errorf("can't init tg: %w", err)
	}

	d, err := python_downloader.NewDownloader(os.Getenv("YANDEX_TOKEN"))
	if err != nil {
		return fmt.Errorf("can't init yandex: %w", err)
	}

	c := controller.NewController(b, d)
	err = c.Run(ctx)
	if err != nil {
		return err
	}

	return nil
}
