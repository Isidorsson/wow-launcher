<script lang="ts">
  import { servers, selectedServerId } from '../stores';

  let open = $state(false);
  let triggerEl: HTMLButtonElement | undefined = $state();
  let listEl: HTMLUListElement | undefined = $state();

  const current = $derived($servers.find(s => s.id === $selectedServerId));
  const activeIndex = $derived(
    Math.max(0, $servers.findIndex(s => s.id === $selectedServerId))
  );

  function toggle() {
    open = !open;
    if (open) queueMicrotask(focusActive);
  }

  function close() {
    open = false;
    triggerEl?.focus();
  }

  function focusActive() {
    const el = listEl?.querySelector<HTMLElement>(`[data-idx="${activeIndex}"]`);
    el?.focus();
  }

  function pick(id: string) {
    selectedServerId.set(id);
    close();
  }

  function onTriggerKey(e: KeyboardEvent) {
    if (e.key === 'ArrowDown' || e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      open = true;
      queueMicrotask(focusActive);
    }
  }

  function onListKey(e: KeyboardEvent, idx: number) {
    if (e.key === 'Escape') { e.preventDefault(); close(); return; }
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      pick($servers[idx].id);
      return;
    }
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      const next = (idx + 1) % $servers.length;
      listEl?.querySelector<HTMLElement>(`[data-idx="${next}"]`)?.focus();
    }
    if (e.key === 'ArrowUp') {
      e.preventDefault();
      const prev = (idx - 1 + $servers.length) % $servers.length;
      listEl?.querySelector<HTMLElement>(`[data-idx="${prev}"]`)?.focus();
    }
    if (e.key === 'Home') {
      e.preventDefault();
      listEl?.querySelector<HTMLElement>(`[data-idx="0"]`)?.focus();
    }
    if (e.key === 'End') {
      e.preventDefault();
      listEl?.querySelector<HTMLElement>(`[data-idx="${$servers.length - 1}"]`)?.focus();
    }
  }

  function onDocClick(e: MouseEvent) {
    if (!open) return;
    const target = e.target as Node;
    if (triggerEl?.contains(target)) return;
    if (listEl?.contains(target)) return;
    open = false;
  }

  $effect(() => {
    if (!open) return;
    document.addEventListener('mousedown', onDocClick);
    return () => document.removeEventListener('mousedown', onDocClick);
  });
</script>

<div class="wrap">
  <span class="lbl" id="realm-label">Realm</span>
  <div class="field">
    <button
      bind:this={triggerEl}
      type="button"
      class="trigger"
      class:open
      aria-haspopup="listbox"
      aria-expanded={open}
      aria-labelledby="realm-label"
      onclick={toggle}
      onkeydown={onTriggerKey}
      disabled={$servers.length === 0}
    >
      <span class="value">
        {current?.name ?? ($servers.length === 0 ? '— No realms configured —' : 'Select realm')}
      </span>
      <svg class="chev" class:flip={open} viewBox="0 0 24 24" width="13" height="13" aria-hidden="true">
        <path d="M6 9l6 6 6-6" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </button>

    {#if open}
      <ul
        bind:this={listEl}
        class="list"
        role="listbox"
        aria-labelledby="realm-label"
        tabindex="-1"
      >
        {#each $servers as srv, i}
          <li
            role="option"
            tabindex="-1"
            data-idx={i}
            class="opt"
            class:selected={srv.id === $selectedServerId}
            aria-selected={srv.id === $selectedServerId}
            onclick={() => pick(srv.id)}
            onkeydown={(e) => onListKey(e, i)}
          >
            <span class="opt-name">{srv.name}</span>
            {#if srv.id === $selectedServerId}
              <svg class="check" viewBox="0 0 24 24" width="12" height="12" aria-hidden="true">
                <path d="M5 13l4 4L19 7" fill="none" stroke="currentColor" stroke-width="2.4" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            {/if}
          </li>
        {/each}
      </ul>
    {/if}
  </div>
</div>

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

  .trigger {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-3);
    padding: 0.4rem 0.55rem 0.4rem var(--space-3);
    background: var(--bg-raised);
    border: 1px solid var(--border-default);
    border-radius: var(--radius-sm);
    color: var(--fg-default);
    font-family: var(--font-heading);
    font-size: var(--fs-sm);
    font-weight: 500;
    text-align: left;
    cursor: pointer;
    transition:
      border-color var(--dur-fast) var(--ease-out),
      background var(--dur-fast) var(--ease-out);
  }
  .trigger:hover:not(:disabled) {
    border-color: var(--border-strong);
    background: var(--bg-elevated);
  }
  .trigger.open {
    border-color: var(--accent);
    background: var(--bg-elevated);
  }
  .trigger:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .value {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    min-width: 0;
  }
  .chev {
    color: var(--fg-mute);
    flex-shrink: 0;
    transition: transform var(--dur-fast) var(--ease-out);
  }
  .chev.flip { transform: rotate(180deg); color: var(--accent); }

  .list {
    position: absolute;
    top: calc(100% + 4px);
    left: 0;
    right: 0;
    margin: 0;
    padding: 4px;
    list-style: none;
    background: var(--bg-raised);
    border: 1px solid var(--border-default);
    border-radius: var(--radius-sm);
    box-shadow: var(--shadow-md);
    z-index: 30;
    max-height: 240px;
    overflow-y: auto;
  }
  .opt {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-2);
    padding: 0.5rem 0.6rem;
    border-radius: var(--radius-sm);
    color: var(--fg-default);
    font-family: var(--font-heading);
    font-size: var(--fs-sm);
    font-weight: 500;
    cursor: pointer;
    outline: none;
  }
  .opt:hover,
  .opt:focus-visible {
    background: var(--bg-elevated);
    color: var(--fg-bright);
  }
  .opt:focus-visible {
    box-shadow: inset 0 0 0 1px var(--accent);
  }
  .opt.selected {
    color: var(--fg-bright);
  }
  .opt-name {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    min-width: 0;
  }
  .check {
    color: var(--accent);
    flex-shrink: 0;
  }
</style>
