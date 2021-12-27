/*
Package internal
Copyright © 2021 Bold City Software <dave@boldcitysoftware.com>
*/
package internal

import (
	"net/http"
	"time"
)

type NapRequestResult struct {
	Name              string
	HttpResponse      *http.Response
	PreRequestResult  string
	PostRequestResult string
	StartTime         time.Time
	EndTime           time.Time
	Error             error
}

func (r *NapRequestResult) GetElapsedMs() int64 {
	return r.EndTime.Sub(r.StartTime).Milliseconds()
}

func NapRequestResultError(err error) *NapRequestResult {
	result := new(NapRequestResult)
	result.Error = err
	return result
}
