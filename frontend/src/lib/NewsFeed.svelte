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
        // No feed configured for this realm — fall back to sample data silently.
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
</script>

<section class="news">
  <header>
    <h2>{srv?.name ?? 'News'}</h2>
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
    <ul class="items">
      {#each displayItems as item, i}
        <article class="item" style="animation-delay: {i * 60}ms;">
          <div class="meta">
            {#if item.category}
              <span class={categoryClass(item.category)}>{item.category}</span>
            {/if}
            <time>{item.date}</time>
          </div>
          <h3>{item.title}</h3>
          <p class="body">{item.body}</p>
          {#if item.url}
            <a href={item.url} target="_blank" rel="noopener">Read the full scroll →</a>
          {/if}
        </article>
      {/each}
    </ul>
  {/if}
</section>

<style>
  .news { padding: 0; }
  header {
    display: flex; align-items: center; justify-content: space-between;
    margin-bottom: var(--space-4);
  }
  header h2 {
    margin: 0;
    font-family: var(--font-display);
    font-weight: 700;
    font-size: var(--fs-lg);
    letter-spacing: 0.18em;
    text-transform: uppercase;
    color: var(--fg-bright);
  }

  .refresh {
    background: transparent;
    border: 1px solid var(--rune-line);
    color: var(--text-soft);
    width: 32px; height: 32px;
    border-radius: 50%;
    display: flex; align-items: center; justify-content: center;
    transition: all 180ms ease;
  }
  .refresh:hover:not(:disabled) {
    border-color: var(--gold);
    color: var(--gold-bright);
    transform: rotate(-90deg);
  }
  .refresh:disabled { opacity: 0.35; cursor: not-allowed; }

  .empty {
    padding: 2.25rem 1.5rem;
    text-align: center;
    color: var(--text-mute);
    font-style: italic;
    font-family: var(--font-script);
    font-size: 1rem;
  }
  .empty.error pre {
    background: rgba(192, 57, 43, 0.08);
    padding: 0.85rem 1rem;
    border-radius: 2px;
    border: 1px solid rgba(192, 57, 43, 0.4);
    color: var(--blood-glow);
    font-family: 'Consolas', monospace;
    font-size: 0.82rem;
    text-align: left;
    white-space: pre-wrap;
    margin: 0.6rem 0 0;
    font-style: normal;
  }

  .sample-note {
    font-family: var(--font-script);
    font-style: italic;
    font-size: 0.78rem;
    color: var(--text-mute);
    text-align: center;
    padding: 0.3rem 0.6rem;
    margin-bottom: 0.75rem;
    border: 1px dashed var(--rune-line);
    border-radius: 2px;
    background: rgba(0,0,0,0.25);
  }

  .items { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 0.9rem; }

  .item {
    position: relative;
    background:
      linear-gradient(180deg, rgba(19, 26, 48, 0.6), rgba(8, 11, 20, 0.78));
    border: 1px solid var(--rune-line);
    padding: 1rem 1.25rem 1.05rem;
    border-radius: 2px;
    box-shadow:
      inset 0 1px 0 rgba(143, 205, 255, 0.04),
      0 4px 14px rgba(0,0,0,0.35);
    opacity: 0;
    transform: translateY(8px);
    animation: rise 420ms cubic-bezier(0.2, 0.8, 0.2, 1) forwards;
    transition: border-color 180ms, box-shadow 180ms;
  }
  .item::before {
    /* left gilt rule */
    content: '';
    position: absolute; left: 0; top: 8px; bottom: 8px;
    width: 2px;
    background: linear-gradient(180deg, transparent, var(--gold), transparent);
    opacity: 0.5;
  }
  .item:hover {
    border-color: var(--rune-line-2);
    box-shadow:
      inset 0 1px 0 rgba(143, 205, 255, 0.08),
      0 6px 22px rgba(0,0,0,0.5),
      0 0 0 1px rgba(78, 164, 255, 0.15);
  }
  @keyframes rise {
    to { opacity: 1; transform: translateY(0); }
  }

  .item h3 {
    margin: 0.4rem 0 0.55rem;
    color: var(--text-bright);
    font-family: var(--font-display);
    font-weight: 600;
    font-size: 1.02rem;
    letter-spacing: 0.06em;
  }
  .meta {
    display: flex; gap: 0.6rem; align-items: center;
    font-size: 0.72rem;
    color: var(--text-mute);
    font-family: var(--font-ui);
  }
  .cat {
    padding: 0.15rem 0.55rem;
    border-radius: 2px;
    background: rgba(0,0,0,0.4);
    border: 1px solid var(--rune-line);
    text-transform: uppercase;
    letter-spacing: 0.12em;
    font-weight: 700;
    font-size: 0.66rem;
    color: var(--text-soft);
  }
  .cat-patch  { color: var(--gold-bright); border-color: rgba(78, 164, 255, 0.5); background: rgba(78, 164, 255, 0.1); }
  .cat-event  { color: var(--fel-glow);    border-color: rgba(76, 175, 80, 0.4);  background: rgba(76, 175, 80, 0.08); }
  .cat-news   { color: var(--arcane);      border-color: rgba(106, 169, 216, 0.4); background: rgba(106, 169, 216, 0.08); }
  .cat-lore   { color: #c9a86a;            border-color: rgba(201, 168, 106, 0.4);  background: rgba(201, 168, 106, 0.08); }

  time { font-style: italic; }

  .body {
    color: var(--text);
    line-height: 1.6;
    margin: 0.25rem 0 0.6rem;
    white-space: pre-wrap;
    font-size: 0.96rem;
  }
  a {
    color: var(--gold-bright);
    font-family: var(--font-display);
    font-size: 0.74rem;
    text-decoration: none;
    letter-spacing: 0.18em;
    text-transform: uppercase;
    font-weight: 600;
    transition: color 150ms, text-shadow 150ms;
  }
  a:hover {
    color: var(--gold-flash);
    text-shadow: 0 0 10px rgba(143, 205, 255, 0.5);
  }
</style>
