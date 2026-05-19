package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"wow-launcher/internal/config"
	"wow-launcher/internal/downloader"
	"wow-launcher/internal/events"
	"wow-launcher/internal/install"
	"wow-launcher/internal/launch"
	"wow-launcher/internal/manifest"
	"wow-launcher/internal/profile"
)

// App is the facade exposed to Svelte via Wails binding.
// Methods named with capital letter become callable JS functions.
type App struct {
	ctx context.Context
	cfg *config.Config
	dl  *downloader.Downloader
	pm  *profile.Manager
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
	p, err := a.pm.Create(serverID, inst.Root, inst.Locale)
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
	m, err := manifest.Fetch(a.ctx, srv.ManifestURL, a.cfg.Security.ManifestPubkeyHex)
	if err != nil {
		return err
	}
	events.Emit(a.ctx, events.StatusMessage, fmt.Sprintf("Syncing %s — %d files", srv.Name, len(m.Files)))
	if err := a.dl.SyncManifest(a.ctx, p.Root, m, includeOptional); err != nil {
		return err
	}
	// Write realmlist after files are in place.
	if err := install.WriteRealmlist(&install.Install{Root: p.Root, Locale: p.Locale}, m.Realmlist); err != nil {
		return fmt.Errorf("write realmlist: %w", err)
	}
	events.Emit(a.ctx, events.StatusMessage, "Sync complete")
	return nil
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
