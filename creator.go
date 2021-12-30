package main

import (
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"sync"
)

type Creator struct {
	lock    sync.Mutex
	pool    chan *QRCode
	creator int
}

type QRCode struct {
	uuid  string
	image string
}

func NewCreator(creator int, pool int) *Creator {
	c := new(Creator)
	c.pool = make(chan *QRCode, pool)
	c.creator = creator
	c.init()
	return c
}

func (c *Creator) init() {
	for i := 0; i < c.creator; i++ {
		go func() {
			//无限生成
			for {
				u := uuid.New().String()
				i, e := qrcode.Encode("scan://"+u, qrcode.Medium, 256)
				if nil != e {
					continue
				}
				q := &QRCode{
					uuid:  u,
					image: base64.StdEncoding.EncodeToString(i),
				}
				c.pool <- q
			}
		}()
	}
}

func (c *Creator) Get() *QRCode {
	return <-c.pool
}
