package main

import (
	"errors"
	"fmt"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/hanfa/bili-go"
	"log"
	"os"
	"strings"
)

//extractBvidFromUrl uses a very dumb method to extract Bvid from a bilibili video link
func extractBvidFromUrl(url string) (string, error) {
	if !strings.Contains(url, "/video/") {
		return "", errors.New("only url like https://www.bilibili.com/video/<Bvid> is currently supported")
	}
	videoIndex := strings.LastIndex(url, "/video/")
	startIndex := videoIndex + len("/video/")

	questionIndex := strings.Index(url, "?")
	if questionIndex == -1 {
		return url[startIndex:], nil
	} else {
		if url[questionIndex-1] != '/' {
			return url[startIndex:questionIndex], nil
		}
		return url[startIndex : questionIndex-1], nil
	}
}

type DownloaderProgressWriter struct {
	curLength   int
	totalLength int
	window      *astilectron.Window
}

// update the UI as long as the download is making progress
func (m *DownloaderProgressWriter) Write(p []byte) (n int, err error) {
	m.curLength += len(p)
	err = m.window.SendMessage(Message{
		ID: MessageDownloadUpdate,
		Payload: MessageDownloadUpdatePayload{
			CurLength:   m.curLength,
			TotalLength: m.totalLength,
		},
	})
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (m *DownloaderProgressWriter) SetContentLength(contentLength int) (err error) {
	m.totalLength = contentLength
	return nil
}

func main() {
	appLogger := log.New(log.Writer(), log.Prefix(), log.Flags())
	appOption := astilectron.Options{
		AppName:            "bDownloader",
		AppIconDefaultPath: "resources/icon.png",
		BaseDirectoryPath:  ".",
		DataDirectoryPath:  ".",
	}

	app, err := astilectron.New(appLogger, appOption)

	if err != nil {
		appLogger.Fatal(fmt.Errorf("main: creating astilectron failed: %w", err))
	}
	defer app.Close()
	app.HandleSignals()

	dataDirectory := app.Paths().DataDirectory()

	if err = app.Start(); err != nil {
		appLogger.Fatal(fmt.Errorf("main: starting astilectron failed: %w", err))
	}

	// create a downloader window
	window, err := app.NewWindow(dataDirectory+"/resources/app/index.html", &astilectron.WindowOptions{
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(700),
		Width:  astikit.IntPtr(700),
	})
	if err != nil {
		appLogger.Fatal(fmt.Errorf("main: new window failed: %w", err))
	}

	// create bili client
	bClient, err := bili.New()

	// create window listeners
	window.OnMessage(func(m *astilectron.EventMessage) interface{} {
		var message Message
		if err := m.Unmarshal(&message); err != nil {
			return Message{
				ID:      MessageError,
				Payload: "error unmarshalling data from view",
			}
		}

		// extract bvid from the url
		url := message.Payload.(string)
		bvid, err := extractBvidFromUrl(url)
		if err != nil {
			appLogger.Print(fmt.Errorf("main: extract bvid failed: %w", err))
			return Message{
				ID:      MessageError,
				Payload: err.Error(),
			}
		}

		// fetch video information if the url is valid
		appLogger.Print(fmt.Printf("main: extracted bivid %s", bvid))
		videoInfo, err := bClient.GetVideoInfoByBvid(bvid)
		if err != nil {
			return Message{ID: MessageError, Payload: err.Error()}
		}

		pw := DownloaderProgressWriter{window: window}

		switch message.ID {
		case MessageUrlUpdate:
			// client changes the url of the download link
			return Message{
				ID:      MessageUrlUpdate,
				Payload: videoInfo.Data,
			}
		case MessageDownloadRequest:
			// client wants to download the video of given link
			outPath, err := os.UserHomeDir()
			outPath += "/" + bvid + ".mp4"

			if err != nil {
				appLogger.Print(fmt.Errorf("main: cannot determine user home directory"))
				return Message{
					ID:      MessageError,
					Payload: err.Error(),
				}
			}
			err = bClient.DownloadByBvid(bili.DownloadOptionBvid{
				Bvid: bvid,
				DownloadOptionCommon: bili.DownloadOptionCommon{
					Resolution: bili.Stream1080P,
					OutPath:    outPath,
				},
			}, true, &pw)
			if err != nil {
				return Message{
					ID:      MessageError,
					Payload: err.Error(),
				}
			}
			return Message{
				ID:      MessageDownloadRequest,
				Payload: nil,
			}
		}
		return nil
	})

	err = window.Create()
	if err != nil {
		appLogger.Fatal(fmt.Errorf("main: window create failed: %w", err))
	}

	app.Wait()
}
