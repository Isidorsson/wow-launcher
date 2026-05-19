<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { DetectInstalls, GetProfile } from '../../wailsjs/go/main/App';
  import { detectedInstalls, selectedServerId, includeOptional } from '../stores';

  export let open = false;
  const dispatch = createEventDispatcher();

  let profileRoot = '';
  let profileLocale = '';

  $: if (open && $selectedServerId) refreshProfile();

  async function refreshProfile() {
    const p = await GetProfile($selectedServerId);
    profileRoot = p.root ?? '';
    profileLocale = p.locale ?? '';
  }

  async function rescan() {
    const hits = await DetectInstalls();
    detectedInstalls.set(hits);
  }

  function close() {
    open = false;
    dispatch('close');
  }

  function onBackdrop(e: MouseEvent) {
    if (e.target === e.currentTarget) close();
  }
</script>

{#if open}
  <div class="backdrop" on:click={onBackdrop} role="presentation">
    <div class="modal" role="dialog" aria-labelledby="settings-title">
      <header>
        <h2 id="settings-title">Settings</h2>
        <button class="close" on:click={close} aria-label="Close">×</button>
      </header>

      <section>
        <h3>Current profile</h3>
        {#if profileRoot}
          <dl>
            <dt>Install path</dt><dd><code>{profileRoot}</code></dd>
            <dt>Locale</dt><dd>{profileLocale}</dd>
          </dl>
        {:else}
          <p class="muted">No profile created for this realm yet. Set one up from the main screen.</p>
        {/if}
      </section>

      <section>
        <h3>Detected installs</h3>
        <button class="ghost" on:click={rescan}>Rescan</button>
        {#if $detectedInstalls.length === 0}
          <p class="muted">None detected automatically. Drag Wow.exe onto the launcher window, or browse from the main screen.</p>
        {:else}
          <ul class="list">
            {#each $detectedInstalls as inst}
              <li><code>{inst.root}</code> <span class="locale">{inst.locale}</span></li>
            {/each}
          </ul>
        {/if}
      </section>

      <section>
        <h3>Downloads</h3>
        <label class="toggle">
          <input type="checkbox" bind:checked={$includeOptional} />
          Include optional packs (HD textures, cosmetics, …)
        </label>
        <p class="muted small">Applies on the next Update / Install. Optional packs can add several gigabytes.</p>
      </section>

      <footer>
        <button class="primary" on:click={close}>Done</button>
      </footer>
    </div>
  </div>
{/if}

<style>
  .backdrop {
    position: fixed; inset: 0;
    background: rgba(0, 0, 0, 0.55);
    display: flex; align-items: center; justify-content: center;
    z-index: 100;
  }
  .modal {
    width: min(560px, 90vw); max-height: 86vh; overflow-y: auto;
    background: #15151a; border: 1px solid #2a2a33; border-radius: 10px;
    box-shadow: 0 16px 48px rgba(0,0,0,0.5);
    display: flex; flex-direction: column; gap: 0.25rem;
  }
  header { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.25rem; border-bottom: 1px solid #1f1f27; }
  header h2 { margin: 0; color: #fff; font-size: 1.1rem; }
  .close { background: transparent; border: none; color: #888; font-size: 1.6rem; cursor: pointer; line-height: 1; }
  .close:hover { color: #fff; }
  section { padding: 0.9rem 1.25rem; border-bottom: 1px solid #1f1f27; display: flex; flex-direction: column; gap: 0.5rem; }
  section:last-of-type { border-bottom: none; }
  section h3 { margin: 0 0 0.25rem; font-size: 0.8rem; color: #aaa; text-transform: uppercase; letter-spacing: 0.08em; }
  dl { display: grid; grid-template-columns: 7rem 1fr; gap: 0.25rem 0.75rem; margin: 0; }
  dt { color: #888; font-size: 0.85rem; }
  dd { margin: 0; color: #ddd; font-size: 0.85rem; word-break: break-all; }
  code { background: #1f1f27; padding: 0.1rem 0.35rem; border-radius: 3px; font-size: 0.8rem; }
  .list { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 0.25rem; }
  .list li { display: flex; gap: 0.5rem; align-items: center; }
  .locale { font-size: 0.7rem; color: #888; background: #1f1f27; padding: 0.05rem 0.4rem; border-radius: 3px; }
  .toggle { display: flex; gap: 0.5rem; align-items: center; color: #ddd; cursor: pointer; }
  .muted { color: #888; font-size: 0.85rem; margin: 0; }
  .muted.small { font-size: 0.75rem; }
  footer { padding: 1rem 1.25rem; display: flex; justify-content: flex-end; }
  button { padding: 0.5rem 1rem; border: 1px solid #2a2a33; background: #1f1f27; color: #ddd; border-radius: 4px; cursor: pointer; }
  button:hover { background: #28283230; }
  button.primary { background: #c9a227; color: #111; border-color: #c9a227; font-weight: 600; }
  button.primary:hover { background: #ffd966; }
  button.ghost { background: transparent; align-self: flex-start; }
</style>
