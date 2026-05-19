<script lang="ts">
  import {
    DetectInstalls, GetProfile, ValidateInstall, CreateProfile, SyncServer
  } from '../../wailsjs/go/main/App';
  import {
    detectedInstalls, selectedServerId, servers, includeOptional,
    phase, errorMsg, statusMsg
  } from '../stores';

  let {
    open = $bindable(false),
    initialTab = 'installation',
    onclose,
    oninstalled,
  }: {
    open?: boolean;
    initialTab?: 'installation' | 'downloads' | 'profile';
    onclose?: () => void;
    oninstalled?: () => void;
  } = $props();

  type Tab = 'installation' | 'downloads' | 'profile';
  let tab = $state<Tab>('installation');

  let profileRoot = $state('');
  let profileLocale = $state('');
  let manualPath = $state('');
  let chosenBase = $state('');

  const srv = $derived($servers.find(s => s.id === $selectedServerId));
  const profileExists = $derived(!!profileRoot);

  $effect(() => {
    if (open && $selectedServerId) refreshProfile();
  });

  $effect(() => {
    if (open) tab = initialTab;
  });

  $effect(() => {
    if ($detectedInstalls.length > 0 && !chosenBase) {
      chosenBase = $detectedInstalls[0].root;
    }
  });

  async function refreshProfile() {
    const p = await GetProfile($selectedServerId);
    profileRoot = p.root ?? '';
    profileLocale = p.locale ?? '';
  }

  async function rescan() {
    const hits = await DetectInstalls();
    detectedInstalls.set(hits);
  }

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
      await refreshProfile();
      oninstalled?.();
      close();
    } catch (e: any) {
      phase.set('error');
      errorMsg.set(e?.message ?? String(e));
    }
  }

  function close() {
    open = false;
    onclose?.();
  }

  function onBackdrop(e: MouseEvent) {
    if (e.target === e.currentTarget) close();
  }

  function onBackdropKey(e: KeyboardEvent) {
    if (e.key === 'Escape') close();
  }
</script>

{#if open}
  <div
    class="backdrop"
    onclick={onBackdrop}
    onkeydown={onBackdropKey}
    role="button"
    tabindex="-1"
    aria-label="Close settings"
  >
    <div class="modal" role="dialog" aria-labelledby="settings-title">
      <header>
        <div class="title-wrap">
          <h2 id="settings-title">Settings</h2>
        </div>
        <button class="close" onclick={close} aria-label="Close">×</button>
      </header>

      <div class="tabs" role="tablist">
        <button role="tab" class:active={tab === 'installation'} onclick={() => tab = 'installation'}>
          Installation
          {#if !profileExists}<span class="pip" aria-label="setup required"></span>{/if}
        </button>
        <button role="tab" class:active={tab === 'downloads'} onclick={() => tab = 'downloads'}>Downloads</button>
        <button role="tab" class:active={tab === 'profile'} onclick={() => tab = 'profile'}>Profile</button>
      </div>

      <div class="body">
        {#if tab === 'installation'}
          <section>
            <h3>Forge a new installation</h3>
            <p class="hint">
              Drop <code>Wow.exe</code> anywhere on the launcher, or pick a detected install.
            </p>

            {#if $detectedInstalls.length > 0}
              <label for="detected-select">Detected installs</label>
              <select id="detected-select" bind:value={chosenBase}>
                {#each $detectedInstalls as inst}
                  <option value={inst.root}>{inst.root} — {inst.locale}</option>
                {/each}
              </select>
            {:else}
              <p class="muted small">No installs detected automatically. Drop <code>Wow.exe</code> on the window or enter a path below.</p>
            {/if}

            <label for="manual-input">Or enter path manually</label>
            <div class="row">
              <input id="manual-input" type="text" bind:value={manualPath} placeholder="C:\WoW 3.3.5a" />
              <button class="ghost" onclick={validateManual}>Validate</button>
              <button class="ghost" onclick={rescan}>Rescan</button>
            </div>

            <button
              class="primary"
              disabled={!chosenBase || $phase === 'syncing'}
              onclick={install}
            >
              {$phase === 'syncing' ? 'Installing…' : profileExists ? 'Reinstall for this realm' : 'Install for this realm'}
            </button>
          </section>
        {:else if tab === 'downloads'}
          <section>
            <h3>Optional packs</h3>
            <label class="toggle">
              <input type="checkbox" bind:checked={$includeOptional} />
              <span>Include optional packs (HD textures, cosmetics, …)</span>
            </label>
            <p class="muted small">Applies on the next Update / Install. Optional packs can add several gigabytes.</p>
          </section>
        {:else}
          <section>
            <h3>Current profile</h3>
            {#if profileRoot}
              <dl>
                <dt>Realm</dt><dd>{srv?.name ?? '—'}</dd>
                <dt>Install path</dt><dd><code>{profileRoot}</code></dd>
                <dt>Locale</dt><dd>{profileLocale}</dd>
              </dl>
            {:else}
              <p class="muted">No profile bound to this realm yet. Visit the <button class="link" onclick={() => tab = 'installation'}>Installation</button> tab.</p>
            {/if}
          </section>
        {/if}
      </div>

      <footer>
        <button class="primary" onclick={close}>Close</button>
      </footer>
    </div>
  </div>
{/if}

<style>
  .backdrop {
    position: fixed; inset: 0;
    background:
      radial-gradient(60% 50% at 50% 40%, rgba(34, 23, 40, 0.6), rgba(0,0,0,0.82));
    display: flex; align-items: center; justify-content: center;
    z-index: 100;
    border: none;
    backdrop-filter: blur(3px);
    animation: fade 200ms ease-out;
  }
  @keyframes fade { from { opacity: 0; } to { opacity: 1; } }

  .modal {
    width: min(620px, 92vw); max-height: 88vh; overflow: hidden;
    display: flex; flex-direction: column;
    background: var(--bg-surface);
    border: 1px solid var(--border-default);
    box-shadow:
      0 0 0 1px rgba(0,0,0,0.5),
      0 24px 72px rgba(0,0,0,0.75);
    border-radius: var(--radius-md);
    position: relative;
    animation: rise 260ms cubic-bezier(0.2, 0.8, 0.2, 1);
  }
  @keyframes rise {
    from { opacity: 0; transform: translateY(12px) scale(0.985); }
    to   { opacity: 1; transform: translateY(0) scale(1); }
  }
  .modal::before { display: none; }

  header {
    display: flex; align-items: center; justify-content: space-between;
    padding: 1.1rem 1.5rem 0.9rem;
    border-bottom: 1px solid var(--rune-line);
    position: relative;
  }
  header::after { display: none; }
  .title-wrap { display: flex; align-items: center; gap: 0.65rem; }
  header h2 {
    margin: 0; color: var(--fg-bright);
    font-family: var(--font-heading); font-weight: 700;
    font-size: var(--fs-md); letter-spacing: var(--tracking-tight); text-transform: none;
  }
  .close {
    background: transparent; border: 1px solid transparent;
    color: var(--text-soft); font-size: 1.4rem; line-height: 1;
    width: 32px; height: 32px; border-radius: 50%;
  }
  .close:hover { color: var(--gold-bright); border-color: var(--rune-line); }

  .tabs {
    display: flex; gap: 0.25rem;
    padding: 0.6rem 1.5rem 0;
    border-bottom: 1px solid var(--rune-line);
  }
  .tabs button {
    background: transparent; border: none;
    color: var(--text-soft);
    padding: 0.55rem 0.9rem;
    font-family: var(--font-display); font-size: 0.72rem;
    letter-spacing: 0.18em; text-transform: uppercase; font-weight: 600;
    border-bottom: 2px solid transparent;
    margin-bottom: -1px;
    position: relative;
    transition: color 150ms;
  }
  .tabs button:hover { color: var(--gold-bright); }
  .tabs button.active {
    color: var(--rim-bright);
    border-bottom-color: var(--accent);
    text-shadow: 0 0 12px rgba(249, 115, 22, 0.45);
  }
  .pip {
    display: inline-block; width: 6px; height: 6px; margin-left: 6px;
    border-radius: 50%; background: var(--blood-glow);
    box-shadow: 0 0 8px var(--blood-glow);
    animation: pulse 1.4s ease-in-out infinite;
    vertical-align: middle;
  }
  @keyframes pulse {
    0%, 100% { opacity: 0.6; transform: scale(1); }
    50%      { opacity: 1;   transform: scale(1.25); }
  }

  .body { flex: 1; overflow-y: auto; padding: 1.1rem 1.5rem; }
  section { display: flex; flex-direction: column; gap: 0.55rem; }
  section h3 {
    margin: 0 0 0.35rem;
    font-family: var(--font-display); font-weight: 700;
    font-size: 0.78rem; color: var(--text-bright);
    text-transform: uppercase; letter-spacing: 0.18em;
  }
  .hint { color: var(--text-soft); margin: 0 0 0.25rem; font-size: 0.92rem; font-style: italic; }
  label:not(.toggle) {
    font-family: var(--font-display);
    font-size: 0.68rem; color: var(--text-mute);
    text-transform: uppercase; letter-spacing: 0.16em;
    margin-top: 0.45rem;
  }

  select, input[type=text] {
    padding: 0.55rem 0.75rem;
    background: rgba(0,0,0,0.4);
    border: 1px solid var(--rune-line);
    color: var(--text-bright);
    font-family: var(--font-ui); font-size: 0.92rem;
    border-radius: 2px;
    outline: none;
    transition: border-color 150ms, box-shadow 150ms;
  }
  select:focus, input[type=text]:focus {
    border-color: var(--accent);
    box-shadow: 0 0 0 3px rgba(229, 72, 100, 0.22);
  }

  .row { display: flex; gap: 0.5rem; flex-wrap: wrap; }
  .row input { flex: 1; min-width: 12rem; }

  code {
    background: rgba(0,0,0,0.5); padding: 0.1rem 0.4rem; border-radius: 2px;
    border: 1px solid var(--rune-line);
    color: var(--gold-bright);
    font-family: 'Consolas', 'Monaco', monospace; font-size: 0.82rem;
  }

  dl { display: grid; grid-template-columns: 8rem 1fr; gap: 0.4rem 1rem; margin: 0; }
  dt {
    color: var(--text-mute);
    font-family: var(--font-display); font-size: 0.7rem;
    text-transform: uppercase; letter-spacing: 0.14em;
  }
  dd { margin: 0; color: var(--text-bright); font-size: 0.92rem; word-break: break-all; }

  .toggle {
    display: flex; gap: 0.6rem; align-items: center; color: var(--text);
    cursor: pointer;
    padding: 0.4rem 0;
  }
  .toggle input { accent-color: var(--gold); width: 16px; height: 16px; }

  .muted { color: var(--text-mute); margin: 0; font-size: 0.92rem; font-style: italic; }
  .muted.small { font-size: 0.82rem; }
  .link {
    background: transparent; border: none; padding: 0;
    color: var(--gold-bright); text-decoration: underline;
    font: inherit; cursor: pointer;
  }

  footer {
    padding: 0.9rem 1.5rem;
    display: flex; justify-content: flex-end;
    border-top: 1px solid var(--rune-line);
    position: relative;
  }
  footer::before { display: none; }

  button {
    padding: 0.5rem 0.9rem;
    border: 1px solid var(--border-default);
    background: transparent;
    color: var(--fg-default);
    font-family: var(--font-heading);
    font-size: var(--fs-sm);
    letter-spacing: 0;
    text-transform: none;
    font-weight: 500;
    border-radius: var(--radius-sm);
    transition:
      color var(--dur-fast) var(--ease-out),
      border-color var(--dur-fast) var(--ease-out),
      background var(--dur-fast) var(--ease-out);
  }
  button:hover:not(:disabled) {
    border-color: var(--border-strong);
    color: var(--fg-bright);
    background: var(--bg-hover);
  }
  button:disabled { opacity: 0.4; cursor: not-allowed; }

  button.ghost {
    background: transparent;
  }

  button.primary {
    background: var(--fg-bright);
    color: var(--fg-on-accent);
    border-color: var(--fg-bright);
    text-shadow: none;
    box-shadow: none;
    margin-top: 0.4rem;
  }
  button.primary:hover:not(:disabled) {
    background: var(--c-bone-200);
    border-color: var(--c-bone-200);
    color: var(--fg-on-accent);
    box-shadow: none;
  }
</style>
