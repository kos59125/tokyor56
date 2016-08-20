package main

import (
	"github.com/satori/go.uuid"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

var footprint func(*AppLog)

type User struct {
	ID          uuid.UUID
	IP          net.IP
	CurrentPath string
}

func NewUser(landingPage string) *User {
	id := uuid.NewV4()
	ip := newRandomIP()
	return &User{
		ID:          id,
		IP:          ip,
		CurrentPath: landingPage,
	}
}

func newRandomIP() net.IP {
	ip := make([]byte, 4)
	rand.Read(ip)
	return net.IPv4(ip[0], ip[1], ip[2], ip[3])
}

func (user *User) MoveTo(path string) {
	user.CurrentPath = path
}

func (user *User) LeaveFootprint() {
	applog := &AppLog{
		UserID:     user.ID,
		Timestamp:  time.Now().UTC(),
		Path:       user.CurrentPath,
		StatusCode: getStatusCode(user.CurrentPath),
		ClientIP:   user.IP,
	}
	footprint(applog)
}

func SetFootprint(f func(*AppLog)) {
	footprint = f
}

// パスからステータスコードを生成
func getStatusCode(path string) int {
	// パスが /404 のように / + 数字の場合はそれをステータスコードとし、
	// それ以外の場合は 200 とする。
	name := strings.TrimPrefix(path, "/")
	if code, err := strconv.Atoi(name); err == nil {
		return code
	} else {
		return 200
	}
}
