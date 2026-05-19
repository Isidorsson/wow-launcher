<script lang="ts">
  import { onMount } from 'svelte';
  import { GetServers, GetBranding } from '../wailsjs/go/main/App';
  import { servers, selectedServerId } from './stores';
  import ServerList from './lib/ServerList.svelte';
  import PlayPanel from './lib/PlayPanel.svelte';

  let launcherName = 'WoW Launcher';

  onMount(async () => {
    const b = await GetBranding();
    launcherName = b.launcherName;
    const list = await GetServers();
    servers.set(list);
    if (list.length > 0) selectedServerId.set(list[0].id);
  });
</script>

<div class="app">
  <header class="topbar">
    <span class="brand">{launcherName}</span>
  </header>
  <main class="main">
    <ServerList />
    <PlayPanel />
  </main>
</div>

<style>
  :global(html, body, #app) {
    margin: 0; padding: 0; height: 100%;
    background: #121216;
    color: #ddd;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, sans-serif;
  }
  .app { display: flex; flex-direction: column; height: 100vh; }
  .topbar {
    height: 44px; background: #0e0e12; border-bottom: 1px solid #2a2a33;
    display: flex; align-items: center; padding: 0 1rem;
  }
  .brand { font-weight: 600; color: #c9a227; letter-spacing: 0.05em; }
  .main { flex: 1; display: flex; min-height: 0; }
</style>
