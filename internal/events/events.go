package events

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Event names emitted to the Svelte frontend.
const (
	DownloadStart    = "download:start"    // payload: StartPayload
	DownloadProgress = "download:progress" // payload: ProgressPayload
	DownloadDone     = "download:done"     // payload: DonePayload
	DownloadError    = "download:error"    // payload: ErrorPayload
	StatusMessage    = "status:message"    // payload: string
)

type StartPayload struct {
	File      string `json:"file"`
	TotalSize int64  `json:"totalSize"`
	Index     int    `json:"index"`
	OfTotal   int    `json:"ofTotal"`
}

type ProgressPayload struct {
	File         string  `json:"file"`
	Downloaded   int64   `json:"downloaded"`
	Total        int64   `json:"total"`
	BytesPerSec  int64   `json:"bytesPerSec"`
	OverallPct   float64 `json:"overallPct"`
}

type DonePayload struct {
	File string `json:"file"`
}

type ErrorPayload struct {
	File    string `json:"file"`
	Message string `json:"message"`
}

func Emit(ctx context.Context, name string, payload any) {
	if ctx == nil {
		return
	}
	runtime.EventsEmit(ctx, name, payload)
}
