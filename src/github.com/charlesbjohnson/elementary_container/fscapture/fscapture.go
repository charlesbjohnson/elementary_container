package fscapture

import "os"

type Capturable interface {
	Capture() error
	RegisterFileCaptureHook(string, FileCaptureHook) Capturable
	FileCaptureEvents() <-chan FileCaptureEvent
	Path() string
	File() *os.File
	Close() error
}

// TODO
// expand the return value to signify whether to
//   do nothing
//   ignore only this file
//   ignore all files in the same directory
type FileCaptureHook func(path string, info os.FileInfo) bool

type FileCaptureEvent struct {
	Path     string
	Message  string
	Info     os.FileInfo
	Captured bool
}
