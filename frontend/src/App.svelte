<script lang="ts">
  import { onMount } from 'svelte';
  import { EventsOn } from '../wailsjs/runtime/runtime';
  import {
    GetServers, GetBranding, DetectInstalls, GetProfile
  } from '../wailsjs/go/main/App';
  import {
    servers, selectedServerId, detectedInstalls,
    phase, statusMsg, errorMsg, currentFile, overallPct, bytesPerSec,
    updateAvailable
  } from './stores';
  import ServerSelect from './lib/ServerSelect.svelte';
  import TabNav from './lib/TabNav.svelte';
  import NewsFeed from './lib/NewsFeed.svelte';
  import Changelog from './lib/Changelog.svelte';
  import Shop from './lib/Shop.svelte';
  import More from './lib/More.svelte';
  import PlayStrip from './lib/PlayStrip.svelte';
  import Settings from './lib/Settings.svelte';

  let launcherName = $state('Azeroth Launcher');
  let settingsOpen = $state(false);
  let profileExists = $state(false);
  let activeTab = $state('news');

  const tabs = [
    { id: 'news',      label: 'News' },
    { id: 'changelog', label: 'Patch Notes' },
    { id: 'shop',      label: 'Shop' },
    { id: 'more',      label: 'Connect' },
  ];

  const currentServer = $derived(
    $servers.find(s => s.id === $selectedServerId)
  );

  $effect(() => {
    if ($selectedServerId) refreshProfile();
  });

  async function refreshProfile() {
    const p = await GetProfile($selectedServerId);
    profileExists = !!p.exists;
  }

  let initialTab = $state<'installation' | 'downloads' | 'profile'>('installation');

  function openSettings(t: 'installation' | 'downloads' | 'profile' = 'installation') {
    initialTab = t;
    settingsOpen = true;
  }

  onMount(async () => {
    // Block WebView/browser default file-drop (otherwise dropping Wow.exe
    // triggers a navigation/download instead of firing OnFileDrop in Go).
    const preventDefault = (e: DragEvent) => e.preventDefault();
    window.addEventListener('dragenter', preventDefault);
    window.addEventListener('dragover', preventDefault);
    window.addEventListener('drop', preventDefault);

    const b = await GetBranding();
    launcherName = b.launcherName;
    const list = await GetServers();
    servers.set(list);
    if (list.length > 0) selectedServerId.set(list[0].id);

    const hits = await DetectInstalls();
    detectedInstalls.set(hits);

    EventsOn('download:start', (p: any) => {
      phase.set('syncing');
      currentFile.set(p.file);
      statusMsg.set(`Downloading ${p.index}/${p.ofTotal}`);
    });
    EventsOn('download:progress', (p: any) => {
      overallPct.set(p.overallPct);
      bytesPerSec.set(p.bytesPerSec);
    });
    EventsOn('download:error', (p: any) => {
      phase.set('error');
      errorMsg.set(`${p.file}: ${p.message}`);
    });
    EventsOn('status:message', (msg: string) => statusMsg.set(msg));
    EventsOn('drop:install', (inst: { root: string; locale: string }) => {
      detectedInstalls.update(list => {
        if (list.some(i => i.root === inst.root)) return list;
        return [...list, inst];
      });
      statusMsg.set(`Detected: ${inst.root}`);
      errorMsg.set('');
      if (!profileExists) openSettings('installation');
    });
    EventsOn('drop:error', (msg: string) => {
      errorMsg.set(String(msg));
    });
    EventsOn('update:available', (p: { serverId: string; serverName: string; fileCount: number }) => {
      updateAvailable.set(p);
    });
  });

  function dismissUpdate() {
    updateAvailable.set(null);
  }

  function onInstalled() {
    profileExists = true;
  }
</script>

<div class="app">
  <header class="topbar">
    <div class="brand-wrap">
      <span class="crest" aria-hidden="true">
        <svg viewBox="0 0 32 32" width="22" height="22">
          <path d="M16 3 L26 9 L26 21 L16 27 L6 21 L6 9 Z"
                fill="none" stroke="#ededee" stroke-width="1.6" stroke-linejoin="round"/>
          <circle cx="16" cy="15" r="2" fill="#4fb3d9"/>
        </svg>
      </span>
      <div class="brand-text">
        <span class="brand-eyebrow">Launcher</span>
        <h1 class="brand">{launcherName}</h1>
      </div>
    </div>

    <div class="topbar-center">
      <TabNav {tabs} bind:active={activeTab} />
    </div>

    <div class="topbar-right">
      <ServerSelect />
      <button
        class="cogwheel"
        class:alert={!profileExists && $selectedServerId}
        onclick={() => openSettings('installation')}
        aria-label="Settings"
        title="Settings"
      >
        <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="3"/>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09a1.65 1.65 0 0 0-1-1.51 1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09a1.65 1.65 0 0 0 1.51-1 1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33h0a1.65 1.65 0 0 0 1-1.51V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82v0a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
        </svg>
        {#if !profileExists && $selectedServerId}<span class="pip"></span>{/if}
      </button>
    </div>
  </header>

  {#if $updateAvailable && $updateAvailable.serverId === $selectedServerId}
    <div class="update-banner" role="status">
      <span class="update-dot" aria-hidden="true"></span>
      <span class="update-text">
        Patches updated for <strong>{$updateAvailable.serverName}</strong>
        — {$updateAvailable.fileCount} file{$updateAvailable.fileCount === 1 ? '' : 's'} in manifest. Click Sync to download.
      </span>
      <button class="update-dismiss" onclick={dismissUpdate} aria-label="Dismiss">×</button>
    </div>
  {/if}

  <main class="main">
    {#if !$selectedServerId}
      <div class="empty"><p>Select a realm to begin.</p></div>
    {:else}
      <div class="realm-hero">
        <div class="realm-meta">
          <span class="realm-eyebrow">Connected realm</span>
          <h2 class="realm-name">{currentServer?.name ?? '—'}</h2>
          <div class="realm-stats">
            <span class="stat">
              <span class="stat-dot" class:online={profileExists}></span>
              <span class="stat-label">{profileExists ? 'Installed' : 'Setup required'}</span>
            </span>
            {#if currentServer?.website}
              <span class="stat divider"></span>
              <a class="stat link" href={currentServer.website} target="_blank" rel="noopener">
                {currentServer.website.replace(/^https?:\/\//, '')}
              </a>
            {/if}
          </div>
        </div>
        {#if !profileExists}
          <button class="setup-cta" onclick={() => openSettings('installation')}>
            <span class="cta-eyebrow">Get started</span>
            <span class="cta-label">Install for this realm</span>
            <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <path d="M5 12h14M13 6l6 6-6 6"/>
            </svg>
          </button>
        {/if}
      </div>

      <div
        class="panel"
        id="panel-news"
        role="tabpanel"
        aria-labelledby="tab-news"
        hidden={activeTab !== 'news'}
      >
        <NewsFeed />
      </div>
      <div
        class="panel"
        id="panel-changelog"
        role="tabpanel"
        aria-labelledby="tab-changelog"
        hidden={activeTab !== 'changelog'}
      >
        <Changelog />
      </div>
      <div
        class="panel"
        id="panel-shop"
        role="tabpanel"
        aria-labelledby="tab-shop"
        hidden={activeTab !== 'shop'}
      >
        <Shop />
      </div>
      <div
        class="panel"
        id="panel-more"
        role="tabpanel"
        aria-labelledby="tab-more"
        hidden={activeTab !== 'more'}
      >
        <More />
      </div>
    {/if}
  </main>

  <PlayStrip {profileExists} />

  <Settings
    bind:open={settingsOpen}
    {initialTab}
    onclose={() => { settingsOpen = false; refreshProfile(); }}
    oninstalled={onInstalled}
  />
</div>

<style>
  .app {
    display: flex;
    flex-direction: column;
    height: 100vh;
    position: relative;
    z-index: 3;
    --wails-drop-target: drop;
  }

  /* ---------- Top bar ---------- */
  /* ---------- Topbar ---------- */
  .topbar {
    display: grid;
    grid-template-columns: minmax(0, 1fr) auto minmax(0, 1fr);
    align-items: center;
    gap: var(--space-5);
    padding: var(--space-3) var(--space-6);
    background: var(--bg-surface);
    border-bottom: 1px solid var(--border-subtle);
    position: relative;
    z-index: 4;
  }

  .brand-wrap {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    min-width: 0;
  }
  .crest { display: flex; }

  .brand-text {
    display: flex;
    flex-direction: column;
    min-width: 0;
  }
  .brand-eyebrow {
    font-family: var(--font-heading);
    font-size: var(--fs-2xs);
    font-weight: 500;
    color: var(--fg-faint);
    letter-spacing: var(--tracking-wider);
    text-transform: uppercase;
    line-height: 1;
  }
  .brand {
    margin: 3px 0 0;
    font-family: var(--font-heading);
    font-weight: 700;
    font-size: var(--fs-base);
    letter-spacing: var(--tracking-tight);
    color: var(--fg-bright);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    line-height: 1.2;
    text-transform: none;
  }

  .topbar-center {
    display: flex;
    justify-content: center;
  }

  .topbar-right {
    display: flex;
    align-items: end;
    justify-content: flex-end;
    gap: var(--space-3);
  }

  .cogwheel {
    background: transparent;
    border: 1px solid var(--border-subtle);
    color: var(--fg-soft);
    width: 32px; height: 32px;
    border-radius: var(--radius-sm);
    display: flex; align-items: center; justify-content: center;
    position: relative;
    transition:
      color var(--dur-fast) var(--ease-out),
      border-color var(--dur-fast) var(--ease-out),
      background var(--dur-fast) var(--ease-out);
  }
  .cogwheel:hover {
    border-color: var(--border-strong);
    color: var(--fg-bright);
    background: var(--bg-hover);
  }
  .cogwheel.alert {
    color: var(--accent);
    border-color: var(--accent);
  }
  .pip {
    position: absolute; top: 3px; right: 3px;
    width: 6px; height: 6px; border-radius: 50%;
    background: var(--status-error);
  }

  /* ---------- Update banner ---------- */
  .update-banner {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-2) var(--space-5);
    background: var(--bg-surface);
    border-bottom: 1px solid var(--accent);
    color: var(--fg-bright);
    font-family: var(--font-heading);
    font-size: var(--fs-sm);
  }
  .update-dot {
    width: 8px; height: 8px;
    border-radius: 50%;
    background: var(--accent);
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--accent) 25%, transparent);
    flex-shrink: 0;
  }
  .update-text { flex: 1; color: var(--fg-soft); }
  .update-text strong { color: var(--fg-bright); }
  .update-dismiss {
    background: transparent;
    border: none;
    color: var(--fg-mute);
    font-size: 18px;
    line-height: 1;
    padding: 0 var(--space-2);
    cursor: pointer;
    transition: color var(--dur-fast) var(--ease-out);
  }
  .update-dismiss:hover { color: var(--fg-bright); }

  /* ---------- Main ---------- */
  .main {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    padding: var(--space-7) var(--space-7) var(--space-6);
    scroll-behavior: smooth;
  }
  .panel { max-width: 1180px; margin: 0 auto; }
  .panel[hidden] { display: none; }
  .empty {
    padding: var(--space-12);
    color: var(--fg-mute);
    font-family: var(--font-heading);
    text-align: center;
  }

  /* ---------- Realm hero ---------- */
  .realm-hero {
    max-width: 1180px;
    margin: 0 auto var(--space-7);
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    gap: var(--space-6);
    padding: var(--space-5) var(--space-6);
    background: var(--bg-surface);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-md);
    position: relative;
  }

  .realm-meta {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
    min-width: 0;
    flex: 1;
  }
  .realm-eyebrow {
    font-family: var(--font-heading);
    font-size: var(--fs-2xs);
    font-weight: 600;
    color: var(--fg-faint);
    letter-spacing: var(--tracking-wider);
    text-transform: uppercase;
  }
  .realm-name {
    margin: 0;
    font-family: var(--font-heading);
    font-weight: 700;
    font-size: var(--fs-2xl);
    letter-spacing: var(--tracking-tight);
    line-height: 1.1;
    color: var(--fg-bright);
    text-wrap: balance;
  }
  .realm-stats {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    margin-top: var(--space-2);
    font-family: var(--font-mono);
    font-size: var(--fs-2xs);
    letter-spacing: 0;
  }
  .stat {
    display: inline-flex;
    align-items: center;
    gap: var(--space-2);
    color: var(--fg-mute);
    text-transform: none;
  }
  .stat-dot {
    width: 6px; height: 6px;
    border-radius: 50%;
    background: var(--fg-faint);
  }
  .stat-dot.online {
    background: var(--status-success);
  }
  .stat.divider {
    width: 1px; height: 10px;
    background: var(--border-default);
    padding: 0;
  }
  .stat.link {
    color: var(--accent);
    text-decoration: none;
    transition: color var(--dur-fast) var(--ease-out);
  }
  .stat.link:hover { color: var(--accent-bright); text-decoration: underline; }

  .setup-cta {
    display: inline-flex;
    align-items: center;
    gap: var(--space-2);
    padding: 0.55rem 0.9rem;
    background: var(--fg-bright);
    border: 1px solid var(--fg-bright);
    border-radius: var(--radius-sm);
    color: var(--fg-on-accent);
    font-family: var(--font-heading);
    font-size: var(--fs-sm);
    font-weight: 600;
    letter-spacing: 0;
    cursor: pointer;
    transition:
      background var(--dur-fast) var(--ease-out),
      border-color var(--dur-fast) var(--ease-out);
    flex-shrink: 0;
  }
  .setup-cta:hover {
    background: var(--c-bone-200);
    border-color: var(--c-bone-200);
  }
  .setup-cta .cta-eyebrow {
    font-size: var(--fs-2xs);
    font-weight: 600;
    letter-spacing: var(--tracking-wider);
    text-transform: uppercase;
    color: var(--c-carbon-500);
    padding-right: var(--space-2);
    border-right: 1px solid rgba(10, 10, 12, 0.18);
  }
  .setup-cta .cta-label { font-weight: 600; }
  .setup-cta svg { transition: transform var(--dur-fast) var(--ease-out); }
  .setup-cta:hover svg { transform: translateX(2px); }
</style>
