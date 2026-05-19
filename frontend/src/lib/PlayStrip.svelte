<script lang="ts">
  import { SyncServer, Play } from '../../wailsjs/go/main/App';
  import {
    selectedServerId, includeOptional, phase, errorMsg, statusMsg,
    overallPct, currentFile, bytesPerSec, humanBytes
  } from '../stores';

  let { profileExists = false }: { profileExists?: boolean } = $props();

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
        <span class="file">{$currentFile || 'Preparing the rites…'}</span>
        <span class="speed">{$bytesPerSec ? `${humanBytes($bytesPerSec)}/s` : ''}</span>
      </div>
      <div class="bar">
        <div class="fill" style="width: {$overallPct}%"></div>
        <div class="shimmer" style="left: {$overallPct}%"></div>
      </div>
    </div>
  {:else if $errorMsg}
    <div class="status err">{$errorMsg}</div>
  {:else if $statusMsg}
    <div class="status ok">{$statusMsg}</div>
  {:else}
    <div class="status muted">{profileExists ? 'Ready for battle' : 'Not installed for this realm'}</div>
  {/if}

  <div class="actions">
    <button class="ghost" onclick={update} disabled={!profileExists || $phase === 'syncing'}>
      Check Update
    </button>
    <button class="play" onclick={play} disabled={!profileExists || $phase !== 'idle'}>
      <span class="play-label">Play</span>
    </button>
  </div>
</footer>

<style>
  .strip {
    display: flex; align-items: center; gap: 1.25rem;
    padding: 0.85rem 1.5rem;
    background:
      linear-gradient(180deg, rgba(6, 9, 18, 0.95), rgba(19, 26, 48, 0.95));
    border-top: 1px solid var(--rune-line-2);
    position: relative;
    box-shadow: 0 -6px 16px rgba(0,0,0,0.45);
  }
  .strip::before {
    content: '';
    position: absolute; left: 0; right: 0; top: -1px;
    height: 1px;
    background: linear-gradient(90deg, transparent 5%, var(--gold) 50%, transparent 95%);
    opacity: 0.6;
  }

  .progress {
    flex: 1; display: flex; flex-direction: column; gap: 0.35rem;
    min-width: 0;
  }
  .meta {
    display: flex; justify-content: space-between;
    font-size: 0.78rem; color: var(--text-soft);
    font-family: var(--font-script); font-style: italic;
  }
  .file {
    overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
    max-width: 75%;
  }
  .speed { color: var(--gold-bright); font-style: normal; font-family: var(--font-ui); }

  .bar {
    height: 8px;
    background: rgba(0,0,0,0.55);
    border: 1px solid var(--rune-line);
    border-radius: 2px;
    overflow: hidden;
    position: relative;
    box-shadow: inset 0 1px 3px rgba(0,0,0,0.6);
  }
  .fill {
    height: 100%;
    background:
      linear-gradient(90deg, var(--gold-deep), var(--gold) 50%, var(--gold-bright));
    box-shadow:
      0 0 12px rgba(78, 164, 255, 0.65),
      inset 0 1px 0 rgba(216, 236, 255, 0.45);
    transition: width 200ms ease-out;
  }
  .shimmer {
    position: absolute; top: 0; bottom: 0;
    width: 30px; transform: translateX(-30px);
    background: linear-gradient(90deg, transparent, rgba(216, 236, 255, 0.55), transparent);
    transition: left 200ms ease-out;
    pointer-events: none;
  }

  .status {
    flex: 1;
    font-family: var(--font-script);
    font-size: 0.98rem;
    font-style: italic;
    letter-spacing: 0.03em;
  }
  .status.ok    { color: var(--fel-glow); text-shadow: 0 0 8px rgba(124, 226, 129, 0.25); }
  .status.err   { color: var(--blood-glow); text-shadow: 0 0 8px rgba(255, 122, 106, 0.3); }
  .status.muted { color: var(--text-mute); }

  .actions { display: flex; gap: 0.6rem; align-items: center; }

  button {
    padding: 0.6rem 1.1rem;
    border: 1px solid var(--rune-line-2);
    background: linear-gradient(180deg, var(--slate-3), var(--slate-2));
    color: var(--text);
    font-family: var(--font-display);
    font-size: 0.74rem;
    letter-spacing: 0.18em;
    text-transform: uppercase;
    font-weight: 600;
    border-radius: 2px;
    transition: all 160ms ease;
  }
  button:hover:not(:disabled) {
    border-color: var(--gold);
    color: var(--gold-bright);
    box-shadow: inset 0 0 0 1px rgba(78, 164, 255, 0.3);
  }
  button:disabled { opacity: 0.4; cursor: not-allowed; }
  button.ghost { background: transparent; }

  button.play {
    position: relative;
    padding: 0.75rem 2.6rem;
    background:
      linear-gradient(180deg, #7ec4ff 0%, #4ea4ff 45%, #1e60b8 100%);
    color: #04060c;
    border: 1px solid #0d3b78;
    font-family: var(--font-display);
    font-weight: 900;
    font-size: 0.95rem;
    letter-spacing: 0.32em;
    text-transform: uppercase;
    text-shadow: 0 1px 0 rgba(216, 236, 255, 0.55);
    box-shadow:
      inset 0 1px 0 rgba(216, 236, 255, 0.65),
      inset 0 -2px 0 rgba(0,0,0,0.35),
      0 6px 22px rgba(78, 164, 255, 0.45);
    overflow: hidden;
  }
  button.play::before {
    content: '';
    position: absolute; inset: 0;
    background:
      radial-gradient(80% 120% at 50% -10%, rgba(216, 236, 255, 0.6), transparent 60%);
    pointer-events: none;
    opacity: 0.55;
    transition: opacity 200ms ease;
  }
  button.play::after {
    content: '';
    position: absolute; top: 0; bottom: 0;
    left: -40%; width: 30%;
    background: linear-gradient(90deg, transparent, rgba(216, 236, 255, 0.55), transparent);
    transform: skewX(-20deg);
    transition: left 600ms ease;
    pointer-events: none;
  }
  button.play:hover:not(:disabled) {
    background: linear-gradient(180deg, #cfe9ff, #7ec4ff 55%, #2a7ad8);
    color: #02030a;
    box-shadow:
      inset 0 1px 0 rgba(255, 255, 255, 0.8),
      inset 0 -2px 0 rgba(0,0,0,0.4),
      0 0 0 1px rgba(216, 236, 255, 0.4),
      0 12px 36px rgba(143, 205, 255, 0.65);
    transform: translateY(-1px);
  }
  button.play:hover:not(:disabled)::before { opacity: 1; }
  button.play:hover:not(:disabled)::after { left: 130%; }
  button.play:not(:disabled) .play-label {
    animation: glow 2.6s ease-in-out infinite;
    position: relative;
    z-index: 1;
  }
  @keyframes glow {
    0%, 100% { text-shadow: 0 1px 0 rgba(216, 236, 255, 0.55); }
    50%      { text-shadow: 0 1px 0 rgba(216, 236, 255, 0.55), 0 0 16px rgba(216, 236, 255, 0.75); }
  }
</style>
