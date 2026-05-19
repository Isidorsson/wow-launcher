<script lang="ts">
  import { FetchNews } from '../../wailsjs/go/main/App';
  import type { news } from '../../wailsjs/go/models';
  import { selectedServerId, servers } from '../stores';

  let items = $state<news.Item[]>([]);
  let loading = $state(false);
  let error = $state('');

  const srv = $derived($servers.find(s => s.id === $selectedServerId));

  const placeholder: news.Item[] = [
    {
      title: 'Patch 3.3.5a — Fall of the Lich King',
      date: '2026-05-12',
      body: 'The Frozen Throne stands shattered. New seasonal affixes rotate weekly across Icecrown Citadel. Loot tables rebalanced; Shadowmourne quest chain re-enabled for the second wave of champions.',
      url: 'https://example.invalid/patch-3-3-5a',
      category: 'patch',
    },
    {
      title: 'Hallow\'s End returns to Azeroth',
      date: '2026-05-08',
      body: 'The Headless Horseman rides once more from Scarlet Monastery. Daily candy buckets, wickerman bonfires across capital cities, and limited-time toy drops are live until the 22nd.',
      url: 'https://example.invalid/hallows-end',
      category: 'event',
    },
    {
      title: 'Realm restart Thursday 04:00 server time',
      date: '2026-05-05',
      body: 'Scheduled hardware swap on the auth shard. Expected downtime ~30 minutes. Active characters will be safely logged out beforehand. Vendor refund window extended by 2 hours after restart.',
      category: 'news',
    },
    {
      title: 'Cross-faction battlegrounds — beta opt-in',
      date: '2026-04-29',
      body: 'Premade groups can now queue mixed Horde/Alliance compositions in random BG rotation. Rated arenas unaffected. Feedback thread pinned on the forums.',
      url: 'https://example.invalid/xfaction-bg',
      category: 'news',
    },
    {
      title: 'Hotfix — Ruby Sanctum trash respawn',
      date: '2026-04-24',
      body: 'Halion\'s antechamber trash no longer resets mid-pull when the raid wipes near the threshold. Twilight Drakes properly award reputation to the full raid group. Minor server stability improvements.',
      category: 'patch',
    },
    {
      title: 'Noblegarden — eggs hidden across the realm',
      date: '2026-04-19',
      body: 'Brightly painted eggs nestle in the starting zones of every race. Hunt the Spring Rabbit\'s Foot, race the chocolate-fueled mount quest, and don the Noble Garments transmog set, available all week.',
      url: 'https://example.invalid/noblegarden',
      category: 'event',
    },
    {
      title: 'Arena Season 8 concludes — Wrathful gear unlocked',
      date: '2026-04-14',
      body: 'Final standings locked. Top 0.5% receive Vanquisher titles and the Relentless Gladiator\'s Frost Wyrm. Wrathful Gladiator pieces now purchasable with honor at reduced cost. Season 9 ladder opens next reset.',
      url: 'https://example.invalid/arena-s8',
      category: 'news',
    },
    {
      title: 'Lorewalker chronicle — The Lich King\'s Final Hours',
      date: '2026-04-08',
      body: 'A new in-character chronicle from the Argent Crusade archivists recounts the assault on the Frozen Throne. Read it in the Storyteller tab inside Dalaran, or pick up the bound codex from the city vendor.',
      url: 'https://example.invalid/lich-king-chronicle',
      category: 'lore',
    },
    {
      title: 'Anti-cheat sweep — 412 accounts actioned',
      date: '2026-04-02',
      body: 'Recent wave of botting and teleport-hack detections resulted in permanent bans across all realms. Appeal window remains open for 14 days. Report suspicious behavior via the in-game ticket system.',
      category: 'news',
    },
    {
      title: 'Mid-Summer Fire Festival — early calendar notice',
      date: '2026-03-28',
      body: 'Mark your calendars: the Burning Blossom returns mid-June. Flame wardens are being recruited now; honor the bonfires of every capital and desecrate the rival faction\'s flames for unique tabards and the Brazier of Dancing Flames.',
      category: 'event',
    },
  ];

  const displayItems = $derived(items.length > 0 ? items : placeholder);
  const isPlaceholder = $derived(items.length === 0 && !loading && !error);

  const featured = $derived(displayItems[0]);
  const rest = $derived(displayItems.slice(1));

  $effect(() => {
    if ($selectedServerId) load();
  });

  async function load() {
    if (!$selectedServerId) return;
    loading = true;
    error = '';
    items = [];
    try {
      const result = await FetchNews($selectedServerId);
      items = result ?? [];
    } catch (e: any) {
      const msg = e?.message ?? String(e);
      if (/\b404\b|not found/i.test(msg)) {
        error = '';
      } else {
        error = msg;
      }
    } finally {
      loading = false;
    }
  }

  function categoryClass(cat?: string): string {
    if (!cat) return 'cat';
    return `cat cat-${cat.toLowerCase()}`;
  }

  function formatDate(d: string): string {
    if (!d) return '';
    try {
      const dt = new Date(d);
      if (isNaN(dt.getTime())) return d;
      return dt.toLocaleDateString(undefined, { month: 'short', day: 'numeric', year: 'numeric' });
    } catch {
      return d;
    }
  }
</script>

<section class="news">
  <header class="head">
    <div class="head-text">
      <span class="eyebrow">Latest</span>
      <h2>News &amp; Announcements</h2>
    </div>
    <button class="refresh" onclick={load} disabled={loading} title="Refresh" aria-label="Refresh news">
      <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <path d="M21 12a9 9 0 1 1-3-6.7"/>
        <polyline points="21 4 21 9 16 9"/>
      </svg>
    </button>
  </header>

  {#if loading}
    <div class="empty"><p>Consulting the seers…</p></div>
  {:else if error}
    <div class="empty error">
      <p>The ravens return with ill tidings:</p>
      <pre>{error}</pre>
    </div>
  {:else}
    {#if isPlaceholder}
      <div class="sample-note">Sample data — no live feed configured.</div>
    {/if}

    {#if featured}
      <article class="featured" style="animation-delay: 0ms;">
        <div class="featured-body">
          <div class="meta">
            <span class="featured-pill">Featured</span>
            {#if featured.category}
              <span class={categoryClass(featured.category)}>{featured.category}</span>
            {/if}
            <time>{formatDate(featured.date)}</time>
          </div>
          <h3>{featured.title}</h3>
          <p class="body">{featured.body}</p>
          {#if featured.url}
            <a href={featured.url} target="_blank" rel="noopener">
              Read the full article
              <svg viewBox="0 0 24 24" width="12" height="12" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                <path d="M5 12h14M13 6l6 6-6 6"/>
              </svg>
            </a>
          {/if}
        </div>
      </article>
    {/if}

    {#if rest.length > 0}
      <ul class="items">
        {#each rest as item, i}
          <article class="item" style="animation-delay: {(i + 1) * 60}ms;">
            <div class="meta">
              {#if item.category}
                <span class={categoryClass(item.category)}>{item.category}</span>
              {/if}
              <time>{formatDate(item.date)}</time>
            </div>
            <h3>{item.title}</h3>
            <p class="body">{item.body}</p>
            {#if item.url}
              <a href={item.url} target="_blank" rel="noopener">
                Read more
                <svg viewBox="0 0 24 24" width="11" height="11" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
                  <path d="M5 12h14M13 6l6 6-6 6"/>
                </svg>
              </a>
            {/if}
          </article>
        {/each}
      </ul>
    {/if}
  {/if}
</section>

<style>
  .news { padding: 0; }

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

  .refresh {
    background: transparent;
    border: 1px solid var(--border-subtle);
    color: var(--fg-soft);
    width: 32px; height: 32px;
    border-radius: var(--radius-sm);
    display: flex; align-items: center; justify-content: center;
    transition:
      color var(--dur-fast) var(--ease-out),
      border-color var(--dur-fast) var(--ease-out),
      background var(--dur-fast) var(--ease-out);
  }
  .refresh:hover:not(:disabled) {
    border-color: var(--border-strong);
    color: var(--fg-bright);
    background: var(--bg-hover);
  }
  .refresh:disabled { opacity: 0.35; cursor: not-allowed; }

  .empty {
    padding: var(--space-10) var(--space-6);
    text-align: center;
    color: var(--fg-mute);
    font-family: var(--font-heading);
    font-size: var(--fs-base);
  }
  .empty.error pre {
    background: var(--bg-sunken);
    padding: var(--space-3) var(--space-4);
    border-radius: var(--radius-sm);
    border: 1px solid rgba(214, 84, 84, 0.4);
    color: var(--status-error);
    font-family: var(--font-mono);
    font-size: var(--fs-xs);
    text-align: left;
    white-space: pre-wrap;
    margin: var(--space-3) 0 0;
    font-style: normal;
  }

  .sample-note {
    font-family: var(--font-mono);
    font-size: var(--fs-2xs);
    color: var(--fg-faint);
    text-align: left;
    padding: var(--space-2) var(--space-3);
    margin-bottom: var(--space-4);
    border: 1px dashed var(--border-subtle);
    border-radius: var(--radius-sm);
    background: transparent;
  }

  /* ---------- Featured card ---------- */
  .featured {
    background: var(--bg-surface);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-md);
    overflow: hidden;
    margin-bottom: var(--space-5);
    opacity: 0;
    transform: translateY(8px);
    animation: rise 400ms var(--ease-out) forwards;
    position: relative;
  }
  .featured::before {
    /* thin accent stripe on left edge */
    content: '';
    position: absolute; left: 0; top: 0; bottom: 0;
    width: 2px;
    background: var(--accent);
  }

  .featured-body {
    padding: var(--space-6) var(--space-7);
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
    min-width: 0;
  }
  .featured-body h3 {
    margin: 0;
    font-family: var(--font-heading);
    font-weight: 700;
    font-size: var(--fs-xl);
    letter-spacing: var(--tracking-tight);
    line-height: 1.2;
    color: var(--fg-bright);
    text-wrap: balance;
  }
  .featured-body .body {
    color: var(--fg-soft);
    font-size: var(--fs-base);
    line-height: 1.6;
    margin: 0;
    max-width: 70ch;
  }
  .featured-pill {
    font-family: var(--font-heading);
    font-size: var(--fs-2xs);
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: var(--tracking-wider);
    color: var(--accent);
    padding: 0;
    background: transparent;
    border: 0;
    border-radius: 0;
  }

  /* ---------- Card grid ---------- */
  .items {
    list-style: none;
    padding: 0;
    margin: 0;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: var(--space-4);
  }

  .item {
    position: relative;
    background: var(--bg-surface);
    border: 1px solid var(--border-subtle);
    padding: var(--space-5);
    border-radius: var(--radius-md);
    opacity: 0;
    transform: translateY(6px);
    animation: rise 360ms var(--ease-out) forwards;
    transition:
      border-color var(--dur-fast) var(--ease-out),
      background var(--dur-fast) var(--ease-out);
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
    overflow: hidden;
  }
  .item:hover {
    border-color: var(--border-strong);
    background: var(--bg-raised);
  }

  @keyframes rise {
    to { opacity: 1; transform: translateY(0); }
  }

  .item h3 {
    margin: var(--space-1) 0 0;
    color: var(--fg-bright);
    font-family: var(--font-heading);
    font-weight: 700;
    font-size: var(--fs-md);
    letter-spacing: var(--tracking-tight);
    line-height: 1.3;
    text-wrap: balance;
  }
  .meta {
    display: flex; gap: var(--space-3);
    align-items: center;
    font-size: var(--fs-2xs);
    color: var(--fg-faint);
    font-family: var(--font-mono);
    flex-wrap: wrap;
  }
  .cat {
    padding: 0;
    border-radius: 0;
    background: transparent;
    border: 0;
    text-transform: uppercase;
    letter-spacing: var(--tracking-wider);
    font-weight: 600;
    font-size: var(--fs-2xs);
    color: var(--fg-mute);
    font-family: var(--font-heading);
  }
  .cat-patch  { color: var(--accent); }
  .cat-event  { color: var(--c-green-400); }
  .cat-news   { color: var(--fg-soft); }
  .cat-lore   { color: var(--c-amber-400); }

  time { font-style: normal; color: var(--fg-faint); }

  .body {
    color: var(--fg-soft);
    line-height: 1.55;
    margin: 0;
    white-space: pre-wrap;
    font-size: var(--fs-sm);
    display: -webkit-box;
    -webkit-line-clamp: 3;
    line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .featured-body .body {
    display: block;
    -webkit-line-clamp: unset;
    line-clamp: unset;
  }

  a {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    color: var(--accent);
    font-family: var(--font-heading);
    font-size: var(--fs-xs);
    text-decoration: none;
    letter-spacing: 0;
    font-weight: 600;
    margin-top: auto;
    padding-top: var(--space-2);
    transition: color var(--dur-fast) var(--ease-out);
  }
  a svg { transition: transform var(--dur-fast) var(--ease-out); }
  a:hover {
    color: var(--accent-bright);
  }
  a:hover svg { transform: translateX(2px); }

  .featured-body a {
    align-self: flex-start;
    padding: 0.5rem 0.9rem;
    border: 1px solid var(--border-default);
    background: transparent;
    border-radius: var(--radius-sm);
    margin-top: var(--space-2);
    color: var(--fg-bright);
  }
  .featured-body a:hover {
    background: var(--bg-hover);
    border-color: var(--border-strong);
  }
</style>
