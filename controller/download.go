package controller

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v3"

	"github.com/vladimish/mdb/internal/domain"
)

func (c *Controller) HandleDownload(ctx tele.Context) (err error) {
	err = ctx.Notify(tele.UploadingAudio)
	if err != nil {
		return fmt.Errorf("can't send chat action: %w", err)
	}

	commands := strings.Split(ctx.Text(), " ")
	if len(commands) != 2 {
		return domain.ErrInvalidLinkFormat
	}

	song, err := c.downloader.Download(context.TODO(), commands[1])
	if err != nil {
		return fmt.Errorf("can't download track: %w", err)
	}

	songBuff := bytes.NewBuffer(song.SongData)
	coverBuff := bytes.NewBuffer(song.CoverData)

	err = ctx.Send(&tele.Audio{
		File: tele.File{
			FileReader: songBuff,
		},
		Duration: int(song.Duration.Seconds()),
		Caption:  "",
		Thumbnail: &tele.Photo{
			File: tele.File{
				FileReader: coverBuff,
			},
			Width:  400,
			Height: 400,
		},
		Title:     song.Title,
		Performer: song.Artist,
		MIME:      "",
		FileName:  fmt.Sprintf("%s – %s", song.Artist, song.Title),
	})
	if err != nil {
		return fmt.Errorf("can't send response: %w", err)
	}

	return err
}
