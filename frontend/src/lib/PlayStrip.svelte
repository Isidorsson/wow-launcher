<script lang="ts">
  import { SyncServer, Play } from '../../wailsjs/go/main/App';
  import {
    selectedServerId, includeOptional, phase, errorMsg, statusMsg,
    overallPct, currentFile, bytesPerSec, humanBytes
  } from '../stores';

  export let profileExists = false;

  async function update() {
    if (!profileExists) return;
    phase.set('syncing');
    errorMsg.set('');
    try {
      await SyncServer($selectedServerId, $includeOptional);
      phase.set('idle');
      statusMsg.set('Up to date');
    } catch (e: any) {
      phase.set('error');
      errorMsg.set(e?.message ?? String(e));
    }
  }

  async function play() {
    if (!profileExists) return;
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

<footer class="strip">
  {#if $phase === 'syncing'}
    <div class="progress">
      <div class="meta">
        <span class="file">{$currentFile || 'Preparing…'}</span>
        <span class="speed">{$bytesPerSec ? `${humanBytes($bytesPerSec)}/s` : ''}</span>
      </div>
      <div class="bar"><div class="fill" style="width: {$overallPct}%"></div></div>
    </div>
  {:else if $errorMsg}
    <div class="status err">{$errorMsg}</div>
  {:else if $statusMsg}
    <div class="status ok">{$statusMsg}</div>
  {:else}
    <div class="status muted">{profileExists ? 'Ready' : 'Not installed for this realm'}</div>
  {/if}

  <div class="actions">
    <button class="ghost" on:click={update} disabled={!profileExists || $phase === 'syncing'}>Check update</button>
    <button class="play" on:click={play} disabled={!profileExists || $phase !== 'idle'}>Play</button>
  </div>
</footer>

<style>
  .strip {
    display: flex; align-items: center; gap: 1rem;
    padding: 0.75rem 1.25rem;
    background: #0e0e12; border-top: 1px solid #2a2a33;
  }
  .progress { flex: 1; display: flex; flex-direction: column; gap: 0.25rem; }
  .meta { display: flex; justify-content: space-between; font-size: 0.75rem; color: #aaa; }
  .file { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 75%; }
  .bar { height: 6px; background: #2a2a33; border-radius: 3px; overflow: hidden; }
  .fill { height: 100%; background: linear-gradient(90deg, #c9a227, #ffd966); transition: width 200ms ease-out; }

  .status { flex: 1; font-size: 0.9rem; }
  .status.ok { color: #8be78b; }
  .status.err { color: #ff8b8b; }
  .status.muted { color: #777; }

  .actions { display: flex; gap: 0.5rem; align-items: center; }
  button {
    padding: 0.5rem 1rem; border: 1px solid #2a2a33;
    background: #1f1f27; color: #ddd; border-radius: 4px; cursor: pointer;
  }
  button:hover:not(:disabled) { background: #28283230; }
  button:disabled { opacity: 0.45; cursor: not-allowed; }
  button.ghost { background: transparent; }
  button.play {
    padding: 0.5rem 2.25rem; background: #c9a227; color: #111;
    border-color: #c9a227; font-weight: 700; font-size: 1rem;
  }
  button.play:hover:not(:disabled) { background: #ffd966; }
</style>
