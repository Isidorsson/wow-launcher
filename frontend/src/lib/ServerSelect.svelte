<script lang="ts">
  import { servers, selectedServerId } from '../stores';

  function onChange(e: Event) {
    const v = (e.currentTarget as HTMLSelectElement).value;
    selectedServerId.set(v);
  }
</script>

<label class="wrap">
  <span class="lbl">Realm</span>
  <div class="field">
    <select value={$selectedServerId} onchange={onChange} aria-label="Select realm">
      {#if $servers.length === 0}
        <option value="">— No realms configured —</option>
      {/if}
      {#each $servers as srv}
        <option value={srv.id}>{srv.name}</option>
      {/each}
    </select>
    <svg class="chev" viewBox="0 0 24 24" width="14" height="14" aria-hidden="true">
      <path d="M6 9l6 6 6-6" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>
  </div>
</label>

<style>
  .wrap {
    display: flex;
    flex-direction: column;
    gap: var(--space-1);
    min-width: 240px;
  }
  .lbl {
    font-family: var(--font-display);
    font-size: 0.62rem;
    letter-spacing: 0.22em;
    text-transform: uppercase;
    color: var(--fg-mute);
  }
  .field {
    position: relative;
    display: flex;
    align-items: center;
  }
  select {
    appearance: none;
    -webkit-appearance: none;
    width: 100%;
    padding: var(--space-2) var(--space-8) var(--space-2) var(--space-3);
    background: var(--bg-raised);
    border: 1px solid var(--border-default);
    border-radius: var(--radius-md);
    color: var(--fg-bright);
    font-family: var(--font-display);
    font-size: var(--fs-sm);
    letter-spacing: 0.12em;
    text-transform: uppercase;
    font-weight: 600;
    cursor: pointer;
    transition: border-color var(--dur-fast) var(--ease-out),
                box-shadow var(--dur-fast) var(--ease-out);
  }
  select:hover {
    border-color: var(--accent);
    box-shadow: 0 0 0 1px var(--accent-glow);
  }
  select option {
    background: var(--bg-sunken);
    color: var(--fg-default);
  }
  .chev {
    position: absolute;
    right: var(--space-3);
    color: var(--accent-soft);
    pointer-events: none;
  }
</style>
