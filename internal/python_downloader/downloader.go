package python_downloader

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/bogem/id3v2"
	"github.com/hajimehoshi/go-mp3"

	"github.com/vladimish/mdb/internal/domain"
)

type Downloader struct {
	token string
}

func NewDownloader(token string) (*Downloader, error) {
	cmd := exec.Command("yandex-music-downloader", "-h")
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("can't run downloader: %w", err)
	}

	return &Downloader{
		token: token,
	}, nil
}

const downloadsLocation = "./downloads"

func (d *Downloader) Download(ctx context.Context, link string) (song *domain.Song, err error) {
	album, track, err := extractSongAndAlbum(link)
	if err != nil {
		return nil, fmt.Errorf("can't extract data from link: %w", err)
	}
	if track == "" {
		return nil, domain.ErrAlbumDownloadsUnsupported
	}

	filename := fmt.Sprintf("%s/%s/%s.mp3", downloadsLocation, album, track)

	cmd := exec.CommandContext(
		ctx,
		"yandex-music-downloader",
		"--session-id",
		fmt.Sprintf(`%s`, d.token),
		"--skip-existing",
		"--embed-cover",
		// "--hq",
		"--path-pattern",
		filename[:len(filename)-4], // remove .mp3 from pattern
		"--url",
		fmt.Sprintf(`%s`, link),
	)
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}

	// TODO: handle concurrency if two users downloading same file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("can't open file: %w", err)
	}
	defer func() {
		err = file.Close()
	}()

	duration, err := getDuration(file)
	if err != nil {
		return nil, fmt.Errorf("can't read duration: %w", err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("can't reset file pointer: %w", err)
	}

	opts, err := id3v2.ParseReader(file, id3v2.Options{Parse: true})
	if err != nil {
		return nil, fmt.Errorf("can't parse id3v2: %w", err)
	}

	pictures := opts.GetFrames(opts.CommonID("Attached picture"))
	if len(pictures) == 0 {
		return nil, fmt.Errorf("no album cover found")
	}

	songData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("can't read file: %w", err)
	}

	song = &domain.Song{
		Title:     opts.Title(),
		Artist:    opts.Artist(),
		Album:     opts.Album(),
		Duration:  duration,
		SongData:  songData,
		CoverData: pictures[0].(id3v2.PictureFrame).Picture,
	}

	return song, err
}

func getDuration(r io.ReadSeeker) (time.Duration, error) {
	_, err := r.Seek(0, 0)
	if err != nil {
		return 0, fmt.Errorf("can't reset file pointer: %w", err)
	}

	decoder, err := mp3.NewDecoder(r)
	if err != nil {
		return 0, err
	}

	length := decoder.Length()
	duration := time.Duration(length) * time.Second / time.Duration(decoder.SampleRate()*2*2)

	return duration, nil
}

func extractSongAndAlbum(uri string) (album string, track string, err error) {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return "", "", fmt.Errorf("error parsing URL: %w", err)
	}

	pathComponents := strings.Split(parsedURL.Path, "/")

	for i, component := range pathComponents {
		if component == "album" {
			if i+1 >= len(pathComponents) {
				return "", "", domain.ErrInvalidLinkFormat
			}
			album = pathComponents[i+1]
		}
		if component == "track" {
			if i+1 >= len(pathComponents) {
				return "", "", domain.ErrInvalidLinkFormat
			}
			track = pathComponents[i+1]
		}
	}

	if album == "" {
		return "", "", domain.ErrInvalidLinkFormat
	}

	return album, track, nil
}
