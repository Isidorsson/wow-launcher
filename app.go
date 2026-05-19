package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"wow-launcher/internal/config"
	"wow-launcher/internal/downloader"
	"wow-launcher/internal/events"
	"wow-launcher/internal/install"
	"wow-launcher/internal/launch"
	"wow-launcher/internal/manifest"
	"wow-launcher/internal/news"
	"wow-launcher/internal/profile"

	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App is the facade exposed to Svelte via Wails binding.
// Methods named with capital letter become callable JS functions.
type App struct {
	ctx context.Context
	cfg *config.Config
	dl  *downloader.Downloader
	pm  *profile.Manager
	ms  *manifest.StateStore
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	cfg, err := config.Load()
	if err != nil {
		// Without config we can't render anything useful. Crash loudly.
		panic(err)
	}
	a.cfg = cfg

	a.dl = downloader.New()

	profilesRoot, err := profilesDir(cfg.Paths.ProfilesSubdir)
	if err != nil {
		panic(err)
	}
	pm, err := profile.NewManager(profilesRoot)
	if err != nil {
		panic(err)
	}
	a.pm = pm

	sdir, err := stateDir(cfg.Paths.ProfilesSubdir)
	if err == nil {
		if ms, err := manifest.OpenStateStore(sdir); err == nil {
			a.ms = ms
		}
	}

	// Register drag-and-drop handler. Fires whenever the user drops files
	// onto the window — we look for Wow.exe (or a folder containing it).
	wruntime.OnFileDrop(ctx, func(x, y int, paths []string) {
		a.HandleDroppedPaths(paths)
	})

	// Background: check each server's manifest for content drift since last
	// run. Doesn't block startup; emits update:available per server on change.
	go a.checkForUpdates()
}

// checkForUpdates runs once at startup. For each server, it does a conditional
// fetch (If-None-Match) and emits update:available if the content hash differs
// from what we saw last time. Failure (network, bad sig) is silent — the user
// can still hit Sync manually.
func (a *App) checkForUpdates() {
	if a.ms == nil {
		return
	}
	for _, srv := range a.cfg.Servers {
		prev := a.ms.Get(srv.ID)
		res, err := manifest.FetchWithETag(a.ctx, srv.ManifestURL, a.cfg.Security.ManifestPubkeyHex, prev.ETag, prev.ManifestHash)
		if err != nil {
			continue
		}
		next := manifest.ServerState{
			ETag:         res.ETag,
			ManifestHash: prev.ManifestHash,
			LastChecked:  time.Now().Unix(),
		}
		if res.Manifest != nil {
			next.ManifestHash = res.ContentHash
		}
		_ = a.ms.Set(srv.ID, next)

		if res.Changed && res.Manifest != nil {
			events.Emit(a.ctx, events.UpdateAvailable, events.UpdatePayload{
				ServerID:   srv.ID,
				ServerName: srv.Name,
				FileCount:  len(res.Manifest.Files),
			})
		}
	}
}

// --- Frontend-callable methods ---

type BrandingDTO struct {
	LauncherName string `json:"launcherName"`
	WindowTitle  string `json:"windowTitle"`
	PrimaryColor string `json:"primaryColor"`
}

func (a *App) GetBranding() BrandingDTO {
	return BrandingDTO{
		LauncherName: a.cfg.Branding.LauncherName,
		WindowTitle:  a.cfg.Branding.WindowTitle,
		PrimaryColor: a.cfg.Branding.PrimaryColor,
	}
}

func (a *App) GetServers() []config.Server { return a.cfg.Servers }

type DetectedInstall struct {
	Root   string `json:"root"`
	Locale string `json:"locale"`
}

func (a *App) DetectInstalls() []DetectedInstall {
	hits := install.AutoDetect()
	out := make([]DetectedInstall, 0, len(hits))
	for _, h := range hits {
		out = append(out, DetectedInstall{Root: h.Root, Locale: h.Locale})
	}
	return out
}

func (a *App) ValidateInstall(path string) (*DetectedInstall, error) {
	inst, err := install.Validate(path)
	if err != nil {
		return nil, err
	}
	return &DetectedInstall{Root: inst.Root, Locale: inst.Locale}, nil
}

// BrowseForInstall opens a native folder picker and validates the chosen dir.
// Returns nil (no error) if the user cancelled, an error if the chosen dir is
// not a valid WoW install, or the validated install on success.
func (a *App) BrowseForInstall() (*DetectedInstall, error) {
	path, err := wruntime.OpenDirectoryDialog(a.ctx, wruntime.OpenDialogOptions{
		Title: "Select your World of Warcraft install folder",
	})
	if err != nil {
		return nil, err
	}
	if path == "" {
		return nil, nil
	}
	inst, err := install.Validate(path)
	if err != nil {
		return nil, err
	}
	return &DetectedInstall{Root: inst.Root, Locale: inst.Locale}, nil
}

type ProfileDTO struct {
	ServerID string `json:"serverId"`
	Root     string `json:"root"`
	Locale   string `json:"locale"`
	Exists   bool   `json:"exists"`
}

func (a *App) GetProfile(serverID string) ProfileDTO {
	p, ok := a.pm.Get(serverID)
	if !ok {
		return ProfileDTO{ServerID: serverID, Exists: false}
	}
	return ProfileDTO{ServerID: serverID, Root: p.Root, Locale: p.Locale, Exists: true}
}

// CreateProfile materializes base client files from baseInstall into a new
// profile dir, ready for patch download.
func (a *App) CreateProfile(serverID, baseInstall string) (*ProfileDTO, error) {
	inst, err := install.Validate(baseInstall)
	if err != nil {
		return nil, fmt.Errorf("validate base install: %w", err)
	}
	p, err := a.pm.Create(a.ctx, serverID, inst.Root, inst.Locale)
	if err != nil {
		return nil, err
	}
	return &ProfileDTO{ServerID: p.ServerID, Root: p.Root, Locale: p.Locale, Exists: true}, nil
}

// SyncServer fetches the server's manifest and downloads any missing/mismatched files.
func (a *App) SyncServer(serverID string, includeOptional bool) error {
	srv, err := a.findServer(serverID)
	if err != nil {
		return err
	}
	p, ok := a.pm.Get(serverID)
	if !ok {
		return errors.New("create profile before syncing")
	}
	var prev manifest.ServerState
	if a.ms != nil {
		prev = a.ms.Get(serverID)
	}
	res, err := manifest.FetchWithETag(a.ctx, srv.ManifestURL, a.cfg.Security.ManifestPubkeyHex, "", prev.ManifestHash)
	if err != nil {
		return err
	}
	m := res.Manifest
	if m == nil {
		// 304 with empty If-None-Match shouldn't happen — but treat as no-op.
		events.Emit(a.ctx, events.StatusMessage, "Already up to date")
		return nil
	}
	events.Emit(a.ctx, events.StatusMessage, fmt.Sprintf("Syncing %s — %d files", srv.Name, len(m.Files)))
	if err := a.dl.SyncManifest(a.ctx, p.Root, m, includeOptional); err != nil {
		return err
	}
	if err := install.WriteRealmlist(&install.Install{Root: p.Root, Locale: p.Locale}, m.Realmlist); err != nil {
		return fmt.Errorf("write realmlist: %w", err)
	}
	if a.ms != nil {
		_ = a.ms.Set(serverID, manifest.ServerState{
			ETag:         res.ETag,
			ManifestHash: res.ContentHash,
			LastChecked:  time.Now().Unix(),
		})
	}
	events.Emit(a.ctx, events.StatusMessage, "Sync complete")
	return nil
}

// FetchNews returns the news feed for a server. Empty slice if URL not set
// or feed empty; error only on hard failure (network/parse).
func (a *App) FetchNews(serverID string) ([]news.Item, error) {
	srv, err := a.findServer(serverID)
	if err != nil {
		return nil, err
	}
	return news.Fetch(a.ctx, srv.NewsFeedURL)
}

// HandleDroppedPaths is invoked by the file-drop handler in main.go.
// It finds a Wow.exe (or its parent dir) in paths and emits a "drop:install"
// event with the validated install info, or "drop:error" on failure.
func (a *App) HandleDroppedPaths(paths []string) {
	var lastValidationErr error
	sawExe := false
	for _, p := range paths {
		candidate := p
		info, err := os.Stat(p)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			if !strings.EqualFold(filepath.Base(p), "Wow.exe") {
				continue
			}
			sawExe = true
			candidate = filepath.Dir(p)
		}
		inst, err := install.Validate(candidate)
		if err != nil {
			lastValidationErr = err
			continue
		}
		events.Emit(a.ctx, "drop:install", DetectedInstall{Root: inst.Root, Locale: inst.Locale})
		return
	}
	switch {
	case lastValidationErr != nil && sawExe:
		events.Emit(a.ctx, "drop:error",
			"Wow.exe found, but its folder is not a full install (missing Data/<locale>/...). "+
				"Use Browse to pick the full WoW folder, or drop the folder itself.")
	case lastValidationErr != nil:
		events.Emit(a.ctx, "drop:error", lastValidationErr.Error())
	default:
		events.Emit(a.ctx, "drop:error",
			"No Wow.exe in dropped items. Drop Wow.exe, the install folder, or use Browse.")
	}
}

func (a *App) Play(serverID string) error {
	p, ok := a.pm.Get(serverID)
	if !ok {
		return errors.New("profile not created")
	}
	_, err := launch.Run(p.Root)
	return err
}

// --- helpers ---

func (a *App) findServer(id string) (*config.Server, error) {
	for i := range a.cfg.Servers {
		if a.cfg.Servers[i].ID == id {
			return &a.cfg.Servers[i], nil
		}
	}
	return nil, fmt.Errorf("unknown server %q", id)
}

func profilesDir(sub string) (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, filepath.FromSlash(sub)), nil
}

// stateDir picks the parent of profilesSubdir for storing non-profile state
// (manifest cache, etc.). E.g. profiles_subdir = "WowLauncher/profiles" →
// state lives in "WowLauncher/state".
func stateDir(profilesSub string) (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	parent := filepath.Dir(filepath.FromSlash(profilesSub))
	if parent == "." || parent == string(filepath.Separator) {
		parent = filepath.FromSlash(profilesSub)
	}
	return filepath.Join(base, parent, "state"), nil
}
