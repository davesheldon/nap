/*
Copyright © 2021 Dave Sheldon <dave@boldcitysoftware.com>
*/
package internal

import (
	"net/http"
	"time"
)

type NapResult struct {
	Name              string
	HttpResponse      *http.Response
	PreRequestResult  string
	PostRequestResult string
	StartTime         time.Time
	EndTime           time.Time
	Error             error
}

func (r *NapResult) GetElapsedMs() int64 {
	return r.EndTime.Sub(r.StartTime).Milliseconds()
}
