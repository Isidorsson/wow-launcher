<script lang="ts">
  import { onMount } from 'svelte';
  import { EventsOn } from '../wailsjs/runtime/runtime';
  import {
    GetServers, GetBranding, DetectInstalls, GetProfile
  } from '../wailsjs/go/main/App';
  import {
    servers, selectedServerId, detectedInstalls,
    phase, statusMsg, errorMsg, currentFile, overallPct, bytesPerSec
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
    { id: 'changelog', label: 'Changelog' },
    { id: 'shop',      label: 'Shop' },
    { id: 'more',      label: 'More' },
  ];

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
  });

  function onInstalled() {
    profileExists = true;
  }
</script>

<div class="app">
  <header class="topbar">
    <div class="brand-wrap">
      <span class="crest" aria-hidden="true">
        <svg viewBox="0 0 32 32" width="22" height="22">
          <defs>
            <linearGradient id="arc" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#d8ecff"/>
              <stop offset="50%" stop-color="#4ea4ff"/>
              <stop offset="100%" stop-color="#1e60b8"/>
            </linearGradient>
            <radialGradient id="arc-core" cx="50%" cy="50%" r="50%">
              <stop offset="0%" stop-color="#ffffff"/>
              <stop offset="100%" stop-color="#67e8f9"/>
            </radialGradient>
          </defs>
          <path d="M16 2 L20 10 L29 11 L22 17 L24 26 L16 22 L8 26 L10 17 L3 11 L12 10 Z"
                fill="url(#arc)" stroke="#0a1124" stroke-width="0.8"/>
          <circle cx="16" cy="15" r="2.2" fill="url(#arc-core)"/>
        </svg>
      </span>
      <h1 class="brand">{launcherName}</h1>
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
        <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="3"/>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09a1.65 1.65 0 0 0-1-1.51 1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09a1.65 1.65 0 0 0 1.51-1 1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33h0a1.65 1.65 0 0 0 1-1.51V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82v0a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
        </svg>
        {#if !profileExists && $selectedServerId}<span class="pip"></span>{/if}
      </button>
    </div>
  </header>

  <nav class="tabnav">
    <TabNav {tabs} bind:active={activeTab} />
  </nav>

  <main class="main">
    {#if !$selectedServerId}
      <div class="empty"><p>Select a realm to begin.</p></div>
    {:else}
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

  .topbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-4);
    padding: var(--space-3) var(--space-5);
    background: linear-gradient(180deg, rgba(19,26,48,0.95), rgba(6,9,18,0.95));
    border-bottom: 1px solid var(--border-subtle);
  }
  .brand-wrap { display: flex; align-items: center; gap: var(--space-3); min-width: 0; }
  .crest {
    filter: drop-shadow(0 0 10px var(--accent-glow));
    display: flex;
  }
  .brand {
    margin: 0;
    font-family: var(--font-display);
    font-weight: 700;
    font-size: var(--fs-md);
    letter-spacing: 0.22em;
    text-transform: uppercase;
    color: var(--fg-bright);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .topbar-right {
    display: flex;
    align-items: flex-end;
    gap: var(--space-3);
  }

  .cogwheel {
    background: transparent;
    border: 1px solid transparent;
    color: var(--fg-soft);
    width: 36px; height: 36px;
    border-radius: 50%;
    display: flex; align-items: center; justify-content: center;
    position: relative;
    transition: all var(--dur-base) var(--ease-out);
    align-self: center;
  }
  .cogwheel:hover {
    border-color: var(--border-default);
    color: var(--accent-soft);
    background: var(--bg-hover);
    transform: rotate(35deg);
  }
  .cogwheel.alert { color: var(--accent-soft); }
  .pip {
    position: absolute; top: 4px; right: 4px;
    width: 8px; height: 8px; border-radius: 50%;
    background: var(--status-error);
    box-shadow: 0 0 8px var(--status-error);
    animation: pulse 1.4s ease-in-out infinite;
  }
  @keyframes pulse {
    0%, 100% { opacity: 0.7; transform: scale(1); }
    50%      { opacity: 1;   transform: scale(1.25); }
  }

  .tabnav {
    padding: 0 var(--space-5);
    background: rgba(6, 9, 18, 0.55);
  }

  .main {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    padding: var(--space-6) var(--space-5);
  }
  .panel {
    max-width: 980px;
    margin: 0 auto;
  }
  .panel[hidden] { display: none; }
  .empty {
    padding: var(--space-8);
    color: var(--fg-mute);
    font-style: italic;
    text-align: center;
  }
</style>
