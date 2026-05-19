<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { ValidateInstall, CreateProfile, SyncServer } from '../../wailsjs/go/main/App';
  import {
    detectedInstalls, selectedServerId, servers, includeOptional,
    phase, errorMsg, statusMsg
  } from '../stores';

  const dispatch = createEventDispatcher<{ installed: void }>();

  let manualPath = '';
  let chosenBase = '';

  $: srv = $servers.find(s => s.id === $selectedServerId);
  $: if ($detectedInstalls.length > 0 && !chosenBase) chosenBase = $detectedInstalls[0].root;

  async function validateManual() {
    if (!manualPath) return;
    try {
      const inst = await ValidateInstall(manualPath);
      detectedInstalls.update(list => {
        if (list.some(i => i.root === inst!.root)) return list;
        return [...list, inst!];
      });
      chosenBase = inst!.root;
      errorMsg.set('');
    } catch (e: any) {
      errorMsg.set(`Invalid install: ${e?.message ?? e}`);
    }
  }

  async function install() {
    if (!chosenBase || !srv) return;
    phase.set('syncing');
    errorMsg.set('');
    try {
      await CreateProfile(srv.id, chosenBase);
      await SyncServer(srv.id, $includeOptional);
      phase.set('idle');
      statusMsg.set('Ready to play');
      dispatch('installed');
    } catch (e: any) {
      phase.set('error');
      errorMsg.set(e?.message ?? String(e));
    }
  }
</script>

<div class="card" class:dropping={$phase === 'idle'}>
  <div class="head">
    <span class="badge">Setup required</span>
    <h3>Install {srv?.name ?? 'realm'}</h3>
  </div>

  <p class="hint">
    Drop <code>Wow.exe</code> anywhere on this window, or pick an install below.
  </p>

  {#if $detectedInstalls.length > 0}
    <label for="detected-select">Detected installs</label>
    <select id="detected-select" bind:value={chosenBase}>
      {#each $detectedInstalls as inst}
        <option value={inst.root}>{inst.root} ({inst.locale})</option>
      {/each}
    </select>
  {/if}

  <label for="manual-input">Or enter path manually</label>
  <div class="row">
    <input id="manual-input" type="text" bind:value={manualPath} placeholder="C:\WoW 3.3.5a" />
    <button on:click={validateManual}>Validate</button>
  </div>

  <button class="primary" disabled={!chosenBase || $phase === 'syncing'} on:click={install}>
    {$phase === 'syncing' ? 'Installing…' : 'Install'}
  </button>
</div>

<style>
  .card {
    background: linear-gradient(180deg, #1f1a10, #15151a);
    border: 1px solid #2d2310; border-radius: 8px;
    padding: 1rem 1.25rem; margin-bottom: 1rem;
    display: flex; flex-direction: column; gap: 0.5rem;
  }
  .head { display: flex; align-items: center; gap: 0.6rem; }
  .badge {
    padding: 0.15rem 0.5rem; border-radius: 10px;
    background: #2d2310; color: #c9a227; font-size: 0.7rem;
    text-transform: uppercase; letter-spacing: 0.06em; font-weight: 600;
  }
  h3 { margin: 0; color: #fff; font-size: 1rem; }
  .hint { color: #aaa; margin: 0; font-size: 0.85rem; }
  code { background: #1f1f27; padding: 0.05rem 0.3rem; border-radius: 3px; font-size: 0.8rem; }
  label { font-size: 0.75rem; color: #888; text-transform: uppercase; letter-spacing: 0.06em; margin-top: 0.25rem; }
  select, input[type=text] {
    padding: 0.45rem 0.65rem; background: #1f1f27;
    border: 1px solid #2a2a33; color: #ddd; border-radius: 4px;
  }
  .row { display: flex; gap: 0.5rem; }
  .row input { flex: 1; }
  button { padding: 0.5rem 1rem; border: 1px solid #2a2a33; background: #1f1f27; color: #ddd; border-radius: 4px; cursor: pointer; }
  button:hover:not(:disabled) { background: #28283230; }
  button:disabled { opacity: 0.5; cursor: not-allowed; }
  button.primary { background: #c9a227; color: #111; border-color: #c9a227; font-weight: 600; margin-top: 0.5rem; }
  button.primary:hover:not(:disabled) { background: #ffd966; }
</style>
