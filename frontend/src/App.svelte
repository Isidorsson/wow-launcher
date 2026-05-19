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
  import ServerList from './lib/ServerList.svelte';
  import NewsFeed from './lib/NewsFeed.svelte';
  import SetupCard from './lib/SetupCard.svelte';
  import PlayStrip from './lib/PlayStrip.svelte';
  import Settings from './lib/Settings.svelte';

  let launcherName = $state('WoW Launcher');
  let settingsOpen = $state(false);
  let profileExists = $state(false);

  $effect(() => {
    if ($selectedServerId) refreshProfile();
  });

  async function refreshProfile() {
    const p = await GetProfile($selectedServerId);
    profileExists = !!p.exists;
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
    <span class="brand">{launcherName}</span>
    <button class="cogwheel" onclick={() => settingsOpen = true} aria-label="Settings" title="Settings">
      ⚙
    </button>
  </header>

  <main class="main">
    <ServerList />
    <section class="content">
      {#if $selectedServerId}
        {#if !profileExists}
          <SetupCard oninstalled={onInstalled} />
        {/if}
        <NewsFeed />
      {:else}
        <div class="empty"><p>Select a realm to begin.</p></div>
      {/if}
    </section>
  </main>

  <PlayStrip {profileExists} />

  <Settings bind:open={settingsOpen} onclose={() => settingsOpen = false} />
</div>

<style>
  :global(html, body, #app) {
    margin: 0; padding: 0; height: 100%;
    background: #121216;
    color: #ddd;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, sans-serif;
    --wails-drop-target: drop;
  }
  .app { display: flex; flex-direction: column; height: 100vh; --wails-drop-target: drop; }
  .topbar {
    height: 44px; background: #0e0e12; border-bottom: 1px solid #2a2a33;
    display: flex; align-items: center; justify-content: space-between;
    padding: 0 1rem;
  }
  .brand { font-weight: 600; color: #c9a227; letter-spacing: 0.05em; }
  .cogwheel {
    background: transparent; border: none; color: #aaa;
    width: 32px; height: 32px; border-radius: 50%; cursor: pointer;
    font-size: 1.15rem;
  }
  .cogwheel:hover { background: #1f1f27; color: #c9a227; }
  .main { flex: 1; display: flex; min-height: 0; }
  .content { flex: 1; display: flex; flex-direction: column; min-height: 0; padding: 1.25rem 1.75rem 0; overflow: hidden; }
  .content :global(.news) { padding: 0; }
  .empty { padding: 2rem; color: #777; }
</style>
