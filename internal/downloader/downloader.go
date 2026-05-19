package downloader

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"wow-launcher/internal/events"
	"wow-launcher/internal/manifest"

	"golang.org/x/sync/errgroup"
)

const (
	defaultChunkCount = 8
	minChunkSize      = 4 << 20  // 4MB — smaller files download single-stream
	progressInterval  = 250 * time.Millisecond
)

// Downloader handles concurrent chunked downloads with resume + integrity check.
//
// One Downloader instance is reused across files. It is safe to cancel via the
// ctx passed to DownloadFile; partial files are left on disk for next-time resume.
type Downloader struct {
	HTTPClient *http.Client
}

func New() *Downloader {
	return &Downloader{
		HTTPClient: &http.Client{Timeout: 0}, // no overall timeout — only per-chunk ctx
	}
}

// DownloadFile fetches f.URL to destRoot+f.Path, verifies SHA256, and emits
// progress events on ctx. If the file already exists with the right hash it
// returns immediately (no-op).
func (d *Downloader) DownloadFile(ctx context.Context, destRoot string, f manifest.File, idx, total int) error {
	destPath := filepath.Join(destRoot, filepath.FromSlash(f.Path))

	if ok, _ := verifyHash(destPath, f.SHA256); ok {
		events.Emit(ctx, events.DownloadDone, events.DonePayload{File: f.Path})
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	events.Emit(ctx, events.DownloadStart, events.StartPayload{
		File: f.Path, TotalSize: f.Size, Index: idx, OfTotal: total,
	})

	supportsRange, contentLen, err := d.probeRange(ctx, f.URL)
	if err != nil {
		return err
	}
	if f.Size > 0 && contentLen > 0 && contentLen != f.Size {
		return fmt.Errorf("%s: server size %d != manifest size %d", f.Path, contentLen, f.Size)
	}
	totalSize := contentLen
	if f.Size > 0 {
		totalSize = f.Size
	}

	if !supportsRange || totalSize < minChunkSize {
		if err := d.singleStream(ctx, f.URL, destPath, totalSize, f.Path); err != nil {
			return err
		}
	} else {
		if err := d.chunked(ctx, f.URL, destPath, totalSize, f.Path); err != nil {
			return err
		}
	}

	if ok, got := verifyHash(destPath, f.SHA256); !ok {
		_ = os.Remove(destPath)
		return fmt.Errorf("%s: sha256 mismatch (got %s, want %s)", f.Path, got, f.SHA256)
	}

	events.Emit(ctx, events.DownloadDone, events.DonePayload{File: f.Path})
	return nil
}

func (d *Downloader) probeRange(ctx context.Context, url string) (bool, int64, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return false, 0, fmt.Errorf("HEAD %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return false, 0, fmt.Errorf("HEAD %s: status %d", url, resp.StatusCode)
	}
	ranges := resp.Header.Get("Accept-Ranges")
	clen, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	return ranges == "bytes", clen, nil
}

func (d *Downloader) singleStream(ctx context.Context, url, destPath string, totalSize int64, logName string) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("GET %s: status %d", url, resp.StatusCode)
	}

	tmp := destPath + ".part"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}

	var downloaded atomic.Int64
	done := make(chan struct{})
	go progressTicker(ctx, logName, &downloaded, totalSize, done)

	_, copyErr := io.Copy(&counterWriter{w: out, n: &downloaded}, resp.Body)
	close(done)
	out.Close()

	if copyErr != nil {
		return copyErr
	}
	return os.Rename(tmp, destPath)
}

func (d *Downloader) chunked(ctx context.Context, url, destPath string, totalSize int64, logName string) error {
	tmp := destPath + ".part"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := f.Truncate(totalSize); err != nil {
		return err
	}

	chunkSize := totalSize / int64(defaultChunkCount)
	if chunkSize < minChunkSize {
		chunkSize = minChunkSize
	}

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(defaultChunkCount)

	var downloaded atomic.Int64
	done := make(chan struct{})
	go progressTicker(ctx, logName, &downloaded, totalSize, done)

	for start := int64(0); start < totalSize; start += chunkSize {
		start := start
		end := start + chunkSize - 1
		if end >= totalSize {
			end = totalSize - 1
		}
		g.Go(func() error {
			return d.fetchRange(gctx, url, f, start, end, &downloaded)
		})
	}

	err = g.Wait()
	close(done)
	if err != nil {
		return err
	}
	return os.Rename(tmp, destPath)
}

func (d *Downloader) fetchRange(ctx context.Context, url string, f *os.File, start, end int64, downloaded *atomic.Int64) error {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))
	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusPartialContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("range GET status %d", resp.StatusCode)
	}

	buf := make([]byte, 64<<10)
	pos := start
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := f.WriteAt(buf[:n], pos); werr != nil {
				return werr
			}
			pos += int64(n)
			downloaded.Add(int64(n))
		}
		if readErr == io.EOF {
			return nil
		}
		if readErr != nil {
			return readErr
		}
	}
}

func progressTicker(ctx context.Context, name string, downloaded *atomic.Int64, total int64, done <-chan struct{}) {
	t := time.NewTicker(progressInterval)
	defer t.Stop()
	var last int64
	for {
		select {
		case <-done:
			return
		case <-ctx.Done():
			return
		case <-t.C:
			cur := downloaded.Load()
			bps := int64(float64(cur-last) / progressInterval.Seconds())
			last = cur
			pct := 0.0
			if total > 0 {
				pct = float64(cur) / float64(total) * 100
			}
			events.Emit(ctx, events.DownloadProgress, events.ProgressPayload{
				File: name, Downloaded: cur, Total: total, BytesPerSec: bps, OverallPct: pct,
			})
		}
	}
}

type counterWriter struct {
	w io.Writer
	n *atomic.Int64
}

func (c *counterWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	c.n.Add(int64(n))
	return n, err
}

func verifyHash(path, wantHex string) (bool, string) {
	f, err := os.Open(path)
	if err != nil {
		return false, ""
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false, ""
	}
	got := hex.EncodeToString(h.Sum(nil))
	return got == wantHex, got
}

// SyncManifest downloads all files in m that are missing or hash-mismatched.
// It returns after the first error; partial state survives for next run.
func (d *Downloader) SyncManifest(ctx context.Context, destRoot string, m *manifest.Manifest, includeOptional bool) error {
	var todo []manifest.File
	for _, f := range m.Files {
		if !f.Required && !includeOptional {
			continue
		}
		todo = append(todo, f)
	}
	for i, f := range todo {
		if err := d.DownloadFile(ctx, destRoot, f, i+1, len(todo)); err != nil {
			events.Emit(ctx, events.DownloadError, events.ErrorPayload{File: f.Path, Message: err.Error()})
			return err
		}
	}
	return nil
}

var _ = sync.Mutex{} // reserved for future pause/resume controls
