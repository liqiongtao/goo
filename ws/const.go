package ws

import "time"

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	filePeriod     = 3 * time.Second
	maxMessageSize = 512
)
