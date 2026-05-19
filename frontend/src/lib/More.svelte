<script lang="ts">
  import { selectedServerId, servers } from '../stores';
  const srv = $derived($servers.find(s => s.id === $selectedServerId));

  const links = $derived([
    { label: 'Website',  href: srv?.website ?? '' },
    { label: 'Discord',  href: '' },
    { label: 'Forums',   href: '' },
    { label: 'Donate',   href: '' },
    { label: 'Wiki',     href: '' },
    { label: 'Vote',     href: '' },
  ]);
</script>

<section class="more" aria-label="More links">
  <h3>Quick links</h3>
  <div class="grid">
    {#each links as l}
      <a class="tile" class:disabled={!l.href} href={l.href || '#'} target="_blank" rel="noopener">
        <span class="dot" aria-hidden="true"></span>
        <span class="label">{l.label}</span>
        {#if !l.href}<span class="badge">soon</span>{/if}
      </a>
    {/each}
  </div>
</section>

<style>
  .more { padding: var(--space-2) 0; }
  h3 {
    font-family: var(--font-display);
    font-size: var(--fs-xs);
    font-weight: 700;
    letter-spacing: 0.22em;
    text-transform: uppercase;
    color: var(--fg-mute);
    margin: 0 0 var(--space-4);
  }
  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: var(--space-3);
  }
  .tile {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    padding: var(--space-4);
    background: var(--bg-raised);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-md);
    color: var(--fg-default);
    text-decoration: none;
    transition: border-color var(--dur-fast) var(--ease-out),
                transform var(--dur-fast) var(--ease-out),
                background var(--dur-fast) var(--ease-out);
  }
  .tile:hover {
    border-color: var(--accent);
    background: var(--bg-hover);
    transform: translateY(-1px);
  }
  .tile.disabled {
    opacity: 0.5;
    pointer-events: none;
  }
  .dot {
    width: 8px; height: 8px; border-radius: 50%;
    background: var(--accent);
    box-shadow: 0 0 10px var(--accent-glow);
    flex-shrink: 0;
  }
  .label {
    flex: 1;
    font-family: var(--font-display);
    font-size: var(--fs-sm);
    letter-spacing: 0.1em;
    text-transform: uppercase;
    font-weight: 600;
  }
  .badge {
    font-size: 0.62rem;
    letter-spacing: 0.18em;
    text-transform: uppercase;
    color: var(--fg-mute);
    padding: 2px 6px;
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-sm);
  }
</style>
