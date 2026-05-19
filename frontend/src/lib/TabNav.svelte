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
    if (e.key !== 'ArrowRight' && e.key !== 'ArrowLeft') return;
    e.preventDefault();
    const next = e.key === 'ArrowRight'
      ? (idx + 1) % tabs.length
      : (idx - 1 + tabs.length) % tabs.length;
    active = tabs[next].id;
    const el = document.getElementById(`tab-${tabs[next].id}`);
    el?.focus();
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
      {t.label}
    </button>
  {/each}
</div>

<style>
  .tabs {
    display: flex;
    gap: var(--space-1);
    border-bottom: 1px solid var(--border-subtle);
    padding: 0 var(--space-2);
  }
  button {
    position: relative;
    background: transparent;
    border: 0;
    color: var(--fg-soft);
    padding: var(--space-3) var(--space-4);
    font-family: var(--font-display);
    font-size: var(--fs-xs);
    font-weight: 600;
    letter-spacing: 0.22em;
    text-transform: uppercase;
    border-bottom: 2px solid transparent;
    margin-bottom: -1px;
    transition: color var(--dur-fast) var(--ease-out),
                border-color var(--dur-fast) var(--ease-out);
  }
  button:hover {
    color: var(--accent-bright);
  }
  button.active {
    color: var(--fg-bright);
    border-bottom-color: var(--accent);
  }
  button.active::after {
    content: '';
    position: absolute;
    left: 50%; bottom: -2px;
    width: 60%;
    height: 2px;
    transform: translateX(-50%);
    background: var(--accent);
    box-shadow: 0 0 12px var(--accent-glow);
    border-radius: var(--radius-pill);
  }
</style>
