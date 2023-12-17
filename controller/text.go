package controller

import (
	tele "gopkg.in/telebot.v3"
)

func (c *Controller) HandleText(ctx tele.Context) error {
	return ctx.Reply(
		"To download from yandex, send link in the following format: "+
			"`/download https://music.yandex.ru/album/1924729/track/31412855`",
		&tele.SendOptions{
			ParseMode: tele.ModeMarkdownV2,
		})
}
