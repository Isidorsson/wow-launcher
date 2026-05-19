<script lang="ts">
  import { selectedServerId, servers } from '../stores';
  const srv = $derived($servers.find(s => s.id === $selectedServerId));

  type LinkDef = { label: string; href: string; desc: string; icon: string };

  const links = $derived<LinkDef[]>([
    { label: 'Website', href: srv?.website ?? '', desc: 'Official realm portal',  icon: 'globe'   },
    { label: 'Discord', href: '',                  desc: 'Community chat & support', icon: 'chat'    },
    { label: 'Forums',  href: '',                  desc: 'Discussions & guides',    icon: 'book'    },
    { label: 'Donate',  href: '',                  desc: 'Support the realm',       icon: 'heart'   },
    { label: 'Wiki',    href: '',                  desc: 'Game knowledge base',     icon: 'scroll'  },
    { label: 'Vote',    href: '',                  desc: 'Rate us on listings',     icon: 'star'    },
  ]);
</script>

<section class="more" aria-label="Connect">
  <header class="head">
    <div class="head-text">
      <span class="eyebrow">Connect</span>
      <h2>Community &amp; Resources</h2>
    </div>
  </header>

  <div class="grid">
    {#each links as l, i}
      <a
        class="tile"
        class:disabled={!l.href}
        href={l.href || '#'}
        target="_blank"
        rel="noopener"
        style="animation-delay: {i * 50}ms;"
      >
        <span class="icon" aria-hidden="true">
          {#if l.icon === 'globe'}
            <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"/>
              <path d="M2 12h20M12 2a15 15 0 0 1 0 20M12 2a15 15 0 0 0 0 20"/>
            </svg>
          {:else if l.icon === 'chat'}
            <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
            </svg>
          {:else if l.icon === 'book'}
            <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
              <path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2zM22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/>
            </svg>
          {:else if l.icon === 'heart'}
            <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
              <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
            </svg>
          {:else if l.icon === 'scroll'}
            <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
              <path d="M8 21h12a2 2 0 0 0 2-2v-2H10"/>
              <path d="M10 17V5a2 2 0 0 0-2-2H4a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h2"/>
            </svg>
          {:else if l.icon === 'star'}
            <svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round">
              <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/>
            </svg>
          {/if}
        </span>
        <div class="text">
          <span class="label">{l.label}</span>
          <span class="desc">{l.desc}</span>
        </div>
        {#if l.href}
          <svg class="arrow" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
            <path d="M7 17L17 7M7 7h10v10"/>
          </svg>
        {:else}
          <span class="badge">soon</span>
        {/if}
      </a>
    {/each}
  </div>
</section>

<style>
  .more { padding: 0; }

  .head {
    display: flex; align-items: flex-end; justify-content: space-between;
    margin-bottom: var(--space-5);
  }
  .head-text { display: flex; flex-direction: column; gap: var(--space-1); }
  .eyebrow {
    font-family: var(--font-heading);
    font-size: var(--fs-2xs);
    font-weight: 600;
    color: var(--fg-faint);
    letter-spacing: var(--tracking-wider);
    text-transform: uppercase;
  }
  .head h2 {
    margin: 0;
    font-family: var(--font-heading);
    font-weight: 700;
    font-size: var(--fs-lg);
    letter-spacing: var(--tracking-tight);
    color: var(--fg-bright);
    line-height: 1.1;
  }

  .grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
    gap: var(--space-4);
  }

  .tile {
    position: relative;
    display: flex;
    align-items: center;
    gap: var(--space-4);
    padding: var(--space-4);
    background: var(--bg-surface);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-md);
    color: var(--fg-default);
    text-decoration: none;
    opacity: 0;
    transform: translateY(6px);
    animation: rise 360ms var(--ease-out) forwards;
    transition:
      border-color var(--dur-fast) var(--ease-out),
      background var(--dur-fast) var(--ease-out);
  }
  @keyframes rise {
    to { opacity: 1; transform: translateY(0); }
  }
  .tile:hover {
    border-color: var(--border-strong);
    background: var(--bg-raised);
  }
  .tile.disabled {
    opacity: 0.45;
    pointer-events: none;
  }

  .icon {
    width: 36px; height: 36px;
    border-radius: var(--radius-sm);
    display: flex; align-items: center; justify-content: center;
    color: var(--accent);
    background: var(--bg-elevated);
    border: 1px solid var(--border-subtle);
    flex-shrink: 0;
  }

  .text {
    display: flex;
    flex-direction: column;
    gap: 2px;
    flex: 1;
    min-width: 0;
  }
  .label {
    font-family: var(--font-heading);
    font-size: var(--fs-base);
    font-weight: 600;
    letter-spacing: 0;
    color: var(--fg-bright);
  }
  .desc {
    font-family: var(--font-body);
    font-size: var(--fs-xs);
    color: var(--fg-mute);
    letter-spacing: 0;
  }

  .arrow {
    color: var(--fg-faint);
    flex-shrink: 0;
    transition: transform var(--dur-fast) var(--ease-out),
                color var(--dur-fast) var(--ease-out);
  }
  .tile:hover .arrow { transform: translate(2px, -2px); color: var(--fg-default); }

  .badge {
    font-family: var(--font-heading);
    font-size: var(--fs-2xs);
    letter-spacing: var(--tracking-wider);
    text-transform: uppercase;
    color: var(--fg-faint);
    padding: 0.15rem 0.45rem;
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-sm);
    flex-shrink: 0;
  }
</style>
