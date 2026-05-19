<script lang="ts">
  import { SyncServer, Play } from '../../wailsjs/go/main/App';
  import {
    selectedServerId, includeOptional, phase, errorMsg, statusMsg,
    overallPct, currentFile, bytesPerSec, humanBytes, servers
  } from '../stores';

  let { profileExists = false }: { profileExists?: boolean } = $props();

  const srv = $derived($servers.find(s => s.id === $selectedServerId));

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
  <div class="strip-inner">
    <div class="dock-left">
      <div class="realm-chip">
        <span class="chip-sigil" aria-hidden="true">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 2 L15 9 L22 10 L17 15 L18 22 L12 19 L6 22 L7 15 L2 10 L9 9 Z"/>
          </svg>
        </span>
        <div class="chip-text">
          <span class="chip-label">Realm</span>
          <span class="chip-name">{srv?.name ?? '—'}</span>
        </div>
      </div>
    </div>

    <div class="dock-center">
      {#if $phase === 'syncing'}
        <div class="progress">
          <div class="meta">
            <span class="file" title={$currentFile}>{$currentFile || 'Preparing patch manifest…'}</span>
            <span class="speed">
              <span class="pct">{Math.round($overallPct)}%</span>
              {#if $bytesPerSec}<span class="sep">·</span><span>{humanBytes($bytesPerSec)}/s</span>{/if}
            </span>
          </div>
          <div class="bar">
            <div class="fill" style="width: {$overallPct}%"></div>
            <div class="shimmer" style="left: {$overallPct}%"></div>
          </div>
        </div>
      {:else if $errorMsg}
        <div class="status err">
          <span class="dot" aria-hidden="true"></span>
          <span class="status-text">{$errorMsg}</span>
        </div>
      {:else if $statusMsg}
        <div class="status ok">
          <span class="dot" aria-hidden="true"></span>
          <span class="status-text">{$statusMsg}</span>
        </div>
      {:else}
        <div class="status muted">
          <span class="dot" aria-hidden="true"></span>
          <span class="status-text">{profileExists ? 'Ready for battle' : 'Setup required for this realm'}</span>
        </div>
      {/if}
    </div>

    <div class="dock-right">
      <button class="ghost" onclick={update} disabled={!profileExists || $phase === 'syncing'} aria-label="Check for updates">
        <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <polyline points="23 4 23 10 17 10"/>
          <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>
        </svg>
        <span>Check Update</span>
      </button>
      <button class="play" onclick={play} disabled={!profileExists || $phase !== 'idle'}>
        <span class="play-glow" aria-hidden="true"></span>
        <span class="play-label">Play</span>
        <svg class="play-arrow" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
          <polygon points="6 4 20 12 6 20 6 4" fill="currentColor"/>
        </svg>
      </button>
    </div>
  </div>
</footer>

<style>
  .strip {
    position: relative;
    padding: var(--space-4) var(--space-6);
    background: var(--bg-surface);
    border-top: 1px solid var(--border-subtle);
    z-index: 4;
  }

  .strip-inner {
    max-width: 1280px;
    margin: 0 auto;
    display: grid;
    grid-template-columns: minmax(140px, 1fr) minmax(0, 2fr) auto;
    align-items: center;
    gap: var(--space-5);
  }

  /* ---------- Realm chip ---------- */
  .realm-chip {
    display: inline-flex;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-2) var(--space-3);
    background: transparent;
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-sm);
    min-width: 0;
    max-width: 280px;
  }
  .chip-sigil {
    width: 24px; height: 24px;
    border-radius: var(--radius-sm);
    background: var(--bg-elevated);
    border: 1px solid var(--border-subtle);
    display: flex; align-items: center; justify-content: center;
    color: var(--accent);
    flex-shrink: 0;
  }
  .chip-text { display: flex; flex-direction: column; min-width: 0; gap: 1px; }
  .chip-label {
    font-family: var(--font-heading);
    font-size: var(--fs-2xs);
    font-weight: 500;
    color: var(--fg-faint);
    letter-spacing: var(--tracking-wider);
    text-transform: uppercase;
    line-height: 1;
  }
  .chip-name {
    font-family: var(--font-heading);
    font-size: var(--fs-sm);
    font-weight: 700;
    color: var(--fg-bright);
    letter-spacing: var(--tracking-tight);
    line-height: 1.2;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  /* ---------- Center: status / progress ---------- */
  .dock-center { min-width: 0; }

  .progress {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
    min-width: 0;
  }
  .meta {
    display: flex;
    justify-content: space-between;
    align-items: baseline;
    gap: var(--space-3);
    font-family: var(--font-mono);
    font-size: var(--fs-xs);
    color: var(--fg-soft);
  }
  .file {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 60%;
    color: var(--fg-soft);
  }
  .speed {
    display: inline-flex;
    align-items: baseline;
    gap: var(--space-2);
    color: var(--fg-default);
    flex-shrink: 0;
  }
  .speed .pct {
    font-family: var(--font-mono);
    font-weight: 600;
    font-size: var(--fs-sm);
    color: var(--fg-bright);
    letter-spacing: 0;
  }
  .speed .sep { color: var(--fg-faint); }

  .bar {
    height: 4px;
    background: var(--bg-sunken);
    border: 0;
    border-radius: var(--radius-sm);
    overflow: hidden;
    position: relative;
  }
  .fill {
    height: 100%;
    background: var(--accent);
    border-radius: inherit;
    transition: width 180ms var(--ease-out);
  }
  .shimmer { display: none; }

  .status {
    display: inline-flex;
    align-items: center;
    gap: var(--space-3);
    font-family: var(--font-heading);
    font-size: var(--fs-sm);
    font-weight: 500;
    letter-spacing: 0.01em;
  }
  .status .dot {
    width: 8px; height: 8px;
    border-radius: 50%;
    flex-shrink: 0;
  }
  .status-text {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  .status.ok { color: var(--fg-default); }
  .status.ok .dot { background: var(--status-success); }
  .status.err { color: var(--status-error); }
  .status.err .dot { background: var(--status-error); }
  .status.muted { color: var(--fg-mute); }
  .status.muted .dot { background: var(--fg-faint); }

  /* ---------- Right: actions ---------- */
  .dock-right {
    display: flex;
    gap: var(--space-3);
    align-items: center;
  }

  button {
    padding: 0.55rem 1rem;
    border: 1px solid var(--border-default);
    background: transparent;
    color: var(--fg-default);
    font-family: var(--font-heading);
    font-size: var(--fs-sm);
    letter-spacing: 0;
    font-weight: 500;
    border-radius: var(--radius-sm);
    transition:
      color var(--dur-fast) var(--ease-out),
      border-color var(--dur-fast) var(--ease-out),
      background var(--dur-fast) var(--ease-out);
    display: inline-flex;
    align-items: center;
    gap: var(--space-2);
  }
  button:hover:not(:disabled) {
    border-color: var(--border-strong);
    color: var(--fg-bright);
    background: var(--bg-hover);
  }
  button:disabled { opacity: 0.4; cursor: not-allowed; }
  button.ghost { background: transparent; }

  /* ---------- Play CTA — solid bone, sharp ---------- */
  button.play {
    position: relative;
    padding: 0.7rem 2.2rem;
    background: var(--fg-bright);
    color: var(--fg-on-accent);
    border: 1px solid var(--fg-bright);
    font-family: var(--font-heading);
    font-weight: 700;
    font-size: var(--fs-base);
    letter-spacing: var(--tracking-wider);
    text-transform: uppercase;
    border-radius: var(--radius-sm);
    transition:
      background var(--dur-fast) var(--ease-out),
      border-color var(--dur-fast) var(--ease-out),
      transform var(--dur-fast) var(--ease-out);
    overflow: hidden;
  }
  button.play .play-glow { display: none; }
  button.play .play-arrow {
    color: var(--fg-on-accent);
  }
  button.play:hover:not(:disabled) {
    background: var(--c-bone-200);
    border-color: var(--c-bone-200);
  }
  button.play:active:not(:disabled) { transform: translateY(1px); }
  button.play:disabled {
    background: var(--bg-elevated);
    border-color: var(--border-default);
    color: var(--fg-faint);
  }
  button.play:disabled .play-arrow { color: var(--fg-faint); }

  @media (max-width: 880px) {
    .strip-inner { grid-template-columns: 1fr auto; }
    .dock-left { display: none; }
  }
</style>
