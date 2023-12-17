package controller

import (
	"context"

	tele "gopkg.in/telebot.v3"

	"github.com/vladimish/mdb/internal/python_downloader"
)

type Controller struct {
	bot        *tele.Bot
	downloader *python_downloader.Downloader
}

func NewController(bot *tele.Bot, downloader *python_downloader.Downloader) *Controller {
	c := Controller{
		bot:        bot,
		downloader: downloader,
	}
	c.registerHandlers()

	return &c
}

func (c *Controller) registerHandlers() {
	c.bot.Handle(tele.OnText, c.HandleText)
	c.bot.Handle("/download", c.HandleDownload)
}

func (c *Controller) Run(ctx context.Context) error {
	go c.bot.Start()

	for {
		select {
		case <-ctx.Done():
			c.bot.Stop()
			return ctx.Err()
		}
	}
}
