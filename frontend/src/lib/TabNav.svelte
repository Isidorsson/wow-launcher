<script lang="ts">
  type Tab = { id: string; label: string; icon?: string };

  let {
    tabs,
    active = $bindable<string>(''),
  }: {
    tabs: Tab[];
    active: string;
  } = $props();

  function select(id: string) {
    active = id;
  }

  function onKey(e: KeyboardEvent, idx: number) {
    if (e.key === 'Home') { e.preventDefault(); active = tabs[0].id; focusTab(tabs[0].id); return; }
    if (e.key === 'End')  { e.preventDefault(); active = tabs[tabs.length - 1].id; focusTab(tabs[tabs.length - 1].id); return; }
    if (e.key !== 'ArrowRight' && e.key !== 'ArrowLeft') return;
    e.preventDefault();
    const next = e.key === 'ArrowRight'
      ? (idx + 1) % tabs.length
      : (idx - 1 + tabs.length) % tabs.length;
    active = tabs[next].id;
    focusTab(tabs[next].id);
  }

  function focusTab(id: string) {
    document.getElementById(`tab-${id}`)?.focus();
  }
</script>

<div class="tabs" role="tablist" aria-label="Content sections">
  {#each tabs as t, i}
    <button
      id="tab-{t.id}"
      role="tab"
      type="button"
      aria-selected={active === t.id}
      aria-controls="panel-{t.id}"
      tabindex={active === t.id ? 0 : -1}
      class:active={active === t.id}
      onclick={() => select(t.id)}
      onkeydown={(e) => onKey(e, i)}
    >
      <span class="label">{t.label}</span>
    </button>
  {/each}
</div>

<style>
  .tabs {
    position: relative;
    display: inline-flex;
    align-items: stretch;
    gap: 0;
  }

  button {
    position: relative;
    background: transparent;
    border: 0;
    color: var(--fg-mute);
    padding: 0.65rem 1rem;
    font-family: var(--font-heading);
    font-size: var(--fs-sm);
    font-weight: 500;
    letter-spacing: 0;
    border-radius: 0;
    transition:
      color var(--dur-fast) var(--ease-out);
  }
  button:hover {
    color: var(--fg-default);
  }
  button.active {
    color: var(--fg-bright);
  }
  button.active::after {
    content: '';
    position: absolute;
    left: 1rem;
    right: 1rem;
    bottom: 0;
    height: 2px;
    background: var(--accent);
  }

  .label { display: inline-block; line-height: 1; }
</style>
