package main

type MessageIDType string

type Message struct {
	ID      MessageIDType `json:"id"`
	Payload interface{}   `json:"payload"`
}

type MessageDownloadUpdatePayload struct {
	CurLength   int `json:"cur_length"`
	TotalLength int `json:"total_length"`
}

const (
	MessageError           MessageIDType = "downloader.error"
	MessageUrlUpdate       MessageIDType = "downloader.url_update"
	MessageDownloadRequest MessageIDType = "downloader.download_request"
	MessageDownloadUpdate  MessageIDType = "downloader.download_update"
)
