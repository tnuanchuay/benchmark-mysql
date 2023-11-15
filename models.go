package main

import "time"

type Read struct {
	Id         int
	Ip         string
	Url        string
	RedirectTo string
	CreateAt   time.Time
}
