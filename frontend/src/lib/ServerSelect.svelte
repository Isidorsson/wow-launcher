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
    <svg class="chev" viewBox="0 0 24 24" width="13" height="13" aria-hidden="true">
      <path d="M6 9l6 6 6-6" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
    </svg>
  </div>
</label>

<style>
  .wrap {
    display: flex;
    flex-direction: column;
    gap: 3px;
    min-width: 200px;
  }
  .lbl {
    font-family: var(--font-heading);
    font-size: var(--fs-2xs);
    font-weight: 500;
    letter-spacing: var(--tracking-wider);
    text-transform: uppercase;
    color: var(--fg-faint);
    line-height: 1;
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
    padding: 0.4rem var(--space-7) 0.4rem var(--space-3);
    background: var(--bg-raised);
    border: 1px solid var(--border-default);
    border-radius: var(--radius-sm);
    color: var(--fg-default);
    font-family: var(--font-heading);
    font-size: var(--fs-sm);
    font-weight: 500;
    letter-spacing: 0;
    cursor: pointer;
    transition: border-color var(--dur-fast) var(--ease-out),
                background var(--dur-fast) var(--ease-out);
  }
  select:hover {
    border-color: var(--border-strong);
    background: var(--bg-elevated);
  }
  select option {
    background: var(--bg-raised);
    color: var(--fg-default);
    font-family: var(--font-body);
  }
  .chev {
    position: absolute;
    right: var(--space-3);
    color: var(--fg-mute);
    pointer-events: none;
  }
</style>
