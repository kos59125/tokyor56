package main

import (
	"github.com/satori/go.uuid"
	"net"
	"time"
)

type AppLog struct {
	UserID     uuid.UUID `json:"user_id"`
	Timestamp  time.Time `json:"timestamp"`
	Path       string    `json:"path"`
	StatusCode int       `json:"status"`
	ClientIP   net.IP    `json:"client_ip"`
}
