<script lang="ts">
  import { FetchNews } from '../../wailsjs/go/main/App';
  import type { news } from '../../wailsjs/go/models';
  import { selectedServerId, servers } from '../stores';

  let items: news.Item[] = [];
  let loading = false;
  let error = '';

  $: srv = $servers.find(s => s.id === $selectedServerId);
  $: if ($selectedServerId) load();

  async function load() {
    if (!$selectedServerId) return;
    loading = true;
    error = '';
    items = [];
    try {
      const result = await FetchNews($selectedServerId);
      items = result ?? [];
    } catch (e: any) {
      error = e?.message ?? String(e);
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
    <button class="refresh" on:click={load} disabled={loading} title="Refresh">↻</button>
  </header>

  {#if !srv?.newsFeedUrl}
    <div class="empty">
      <p>No news feed configured for this realm.</p>
    </div>
  {:else if loading}
    <div class="empty"><p>Loading…</p></div>
  {:else if error}
    <div class="empty error">
      <p>Failed to load news:</p>
      <pre>{error}</pre>
    </div>
  {:else if items.length === 0}
    <div class="empty"><p>No news yet.</p></div>
  {:else}
    <ul class="items">
      {#each items as item}
        <article class="item">
          <div class="meta">
            {#if item.category}
              <span class={categoryClass(item.category)}>{item.category}</span>
            {/if}
            <time>{item.date}</time>
          </div>
          <h3>{item.title}</h3>
          <p class="body">{item.body}</p>
          {#if item.url}
            <a href={item.url} target="_blank" rel="noopener">Read more →</a>
          {/if}
        </article>
      {/each}
    </ul>
  {/if}
</section>

<style>
  .news { flex: 1; padding: 1.25rem 1.75rem; overflow-y: auto; }
  header { display: flex; align-items: center; justify-content: space-between; margin-bottom: 1rem; }
  header h2 { margin: 0; font-size: 1.25rem; color: #fff; }
  .refresh {
    background: transparent; border: 1px solid #2a2a33; color: #aaa;
    width: 32px; height: 32px; border-radius: 50%; cursor: pointer;
    font-size: 1rem; line-height: 1;
  }
  .refresh:hover:not(:disabled) { background: #1f1f27; color: #c9a227; }
  .refresh:disabled { opacity: 0.4; cursor: not-allowed; }

  .empty { padding: 2rem; text-align: center; color: #777; }
  .empty.error pre {
    background: #1f1f27; padding: 0.75rem; border-radius: 4px;
    color: #ff9c9c; font-size: 0.8rem; text-align: left; white-space: pre-wrap;
  }

  .items { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 1rem; }
  .item {
    background: #15151a; border: 1px solid #1f1f27; border-radius: 8px;
    padding: 1rem 1.25rem;
  }
  .item h3 { margin: 0.25rem 0 0.5rem; color: #fff; font-size: 1.05rem; }
  .meta { display: flex; gap: 0.6rem; align-items: center; font-size: 0.75rem; color: #888; }
  .cat {
    padding: 0.1rem 0.5rem; border-radius: 10px; background: #2a2a33;
    text-transform: uppercase; letter-spacing: 0.05em; font-weight: 600;
  }
  .cat-patch  { background: #2d2310; color: #c9a227; }
  .cat-event  { background: #112d12; color: #6cd06c; }
  .cat-news   { background: #11202d; color: #6caac9; }
  .body { color: #c5c5c5; line-height: 1.5; margin: 0.25rem 0 0.5rem; white-space: pre-wrap; }
  a { color: #c9a227; font-size: 0.85rem; text-decoration: none; }
  a:hover { text-decoration: underline; }
</style>
