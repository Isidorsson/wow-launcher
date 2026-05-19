<script lang="ts">
  import { onMount } from 'svelte';
  import { EventsOn } from '../../wailsjs/runtime/runtime';
  import {
    DetectInstalls, ValidateInstall, GetProfile, CreateProfile, SyncServer, Play
  } from '../../wailsjs/go/main/App';
  import {
    selectedServerId, servers, detectedInstalls, phase, statusMsg,
    currentFile, overallPct, bytesPerSec, errorMsg
  } from '../stores';
  import ProgressBar from './ProgressBar.svelte';

  let profileExists = false;
  let chosenBaseInstall = '';
  let includeOptional = false;
  let manualPath = '';

  $: srv = $servers.find(s => s.id === $selectedServerId);

  $: if ($selectedServerId) refreshProfile();

  async function refreshProfile() {
    const p = await GetProfile($selectedServerId);
    profileExists = p.exists;
  }

  onMount(async () => {
    const hits = await DetectInstalls();
    detectedInstalls.set(hits);
    if (hits.length > 0) chosenBaseInstall = hits[0].root;

    EventsOn('download:start', (p: any) => {
      phase.set('syncing');
      currentFile.set(p.file);
      statusMsg.set(`Downloading ${p.index}/${p.ofTotal}`);
    });
    EventsOn('download:progress', (p: any) => {
      overallPct.set(p.overallPct);
      bytesPerSec.set(p.bytesPerSec);
    });
    EventsOn('download:done', () => {});
    EventsOn('download:error', (p: any) => {
      phase.set('error');
      errorMsg.set(`${p.file}: ${p.message}`);
    });
    EventsOn('status:message', (msg: string) => statusMsg.set(msg));
  });

  async function browseManual() {
    if (!manualPath) return;
    try {
      const inst = await ValidateInstall(manualPath);
      detectedInstalls.update(list => [...list, inst!]);
      chosenBaseInstall = inst!.root;
      errorMsg.set('');
    } catch (e: any) {
      errorMsg.set(`Invalid install: ${e?.message ?? e}`);
    }
  }

  async function install() {
    if (!chosenBaseInstall) { errorMsg.set('Pick a base WoW install first'); return; }
    phase.set('syncing');
    errorMsg.set('');
    try {
      await CreateProfile($selectedServerId, chosenBaseInstall);
      profileExists = true;
      await SyncServer($selectedServerId, includeOptional);
      phase.set('idle');
      statusMsg.set('Ready to play');
    } catch (e: any) {
      phase.set('error');
      errorMsg.set(e?.message ?? String(e));
    }
  }

  async function update() {
    phase.set('syncing');
    errorMsg.set('');
    try {
      await SyncServer($selectedServerId, includeOptional);
      phase.set('idle');
      statusMsg.set('Up to date');
    } catch (e: any) {
      phase.set('error');
      errorMsg.set(e?.message ?? String(e));
    }
  }

  async function play() {
    phase.set('launching');
    try {
      await Play($selectedServerId);
      statusMsg.set('Wow.exe launched');
      phase.set('idle');
    } catch (e: any) {
      phase.set('error');
      errorMsg.set(e?.message ?? String(e));
    }
  }
</script>

<section class="panel">
  {#if !srv}
    <p class="muted">Select a realm on the left.</p>
  {:else}
    <header>
      <h1>{srv.name}</h1>
      {#if srv.website}<a href={srv.website} target="_blank">{srv.website}</a>{/if}
    </header>

    {#if !profileExists}
      <div class="install-section">
        <h3>Install base client</h3>
        {#if $detectedInstalls.length > 0}
          <label>Detected installs:</label>
          <select bind:value={chosenBaseInstall}>
            {#each $detectedInstalls as inst}
              <option value={inst.root}>{inst.root} ({inst.locale})</option>
            {/each}
          </select>
        {:else}
          <p class="muted">No installs detected automatically.</p>
        {/if}

        <label>Or enter path manually:</label>
        <div class="manual-row">
          <input type="text" bind:value={manualPath} placeholder="C:\WoW 3.3.5a" />
          <button on:click={browseManual}>Validate</button>
        </div>

        <label class="check">
          <input type="checkbox" bind:checked={includeOptional} />
          Include optional packs (HD textures, etc.)
        </label>

        <button class="primary" disabled={!chosenBaseInstall || $phase === 'syncing'} on:click={install}>
          {$phase === 'syncing' ? 'Installing…' : 'Install'}
        </button>
      </div>
    {:else}
      <div class="play-section">
        <button class="primary play" disabled={$phase === 'syncing'} on:click={play}>Play</button>
        <button on:click={update} disabled={$phase === 'syncing'}>Check for updates</button>
        <label class="check">
          <input type="checkbox" bind:checked={includeOptional} />
          Include optional packs
        </label>
      </div>
    {/if}

    {#if $phase === 'syncing'}
      <ProgressBar />
    {/if}

    {#if $statusMsg}<p class="status">{$statusMsg}</p>{/if}
    {#if $errorMsg}<p class="error">{$errorMsg}</p>{/if}
  {/if}
</section>

<style>
  .panel { flex: 1; padding: 1.5rem 2rem; display: flex; flex-direction: column; gap: 1rem; overflow-y: auto; }
  header { display: flex; flex-direction: column; gap: 0.25rem; }
  h1 { margin: 0; color: #fff; }
  a { color: #c9a227; font-size: 0.85rem; text-decoration: none; }
  a:hover { text-decoration: underline; }
  .muted { color: #777; }
  label { font-size: 0.85rem; color: #aaa; display: block; margin-top: 0.5rem; }
  select, input[type=text] {
    width: 100%; padding: 0.5rem 0.7rem; background: #1f1f27;
    border: 1px solid #2a2a33; color: #ddd; border-radius: 4px;
  }
  .manual-row { display: flex; gap: 0.5rem; }
  .manual-row input { flex: 1; }
  button {
    padding: 0.5rem 1rem; border: 1px solid #2a2a33; background: #1f1f27;
    color: #ddd; border-radius: 4px; cursor: pointer;
  }
  button:hover:not(:disabled) { background: #28283230; }
  button:disabled { opacity: 0.5; cursor: not-allowed; }
  button.primary { background: #c9a227; color: #111; border-color: #c9a227; font-weight: 600; }
  button.primary:hover:not(:disabled) { background: #ffd966; }
  button.play { font-size: 1.1rem; padding: 0.75rem 2rem; }
  .check { display: flex; align-items: center; gap: 0.5rem; }
  .install-section, .play-section { display: flex; flex-direction: column; gap: 0.5rem; }
  .play-section { flex-direction: row; align-items: center; }
  .status { color: #8be78b; margin: 0; }
  .error { color: #ff6b6b; margin: 0; white-space: pre-wrap; }
</style>
