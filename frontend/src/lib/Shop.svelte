<script lang="ts">
  import { selectedServerId, servers } from '../stores';
  const srv = $derived($servers.find(s => s.id === $selectedServerId));
</script>

<section class="shop" aria-label="Shop">
  <header class="head">
    <div class="head-text">
      <span class="eyebrow">Marketplace</span>
      <h2>Shop</h2>
    </div>
  </header>

  <div class="placeholder">
    <div class="icon-wrap" aria-hidden="true">
      <svg viewBox="0 0 24 24" width="32" height="32" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <path d="M3 7h18l-2 12H5L3 7z"/>
        <path d="M8 7V5a4 4 0 0 1 8 0v2"/>
      </svg>
    </div>
    <h3>Coming soon</h3>
    <p>Cash-shop integration for <strong>{srv?.name ?? 'this realm'}</strong> not yet wired.</p>
    <p class="hint">Operator can supply a JSON catalog endpoint to populate this view.</p>
  </div>
</section>

<style>
  .shop { padding: 0; }

  .head {
    display: flex; align-items: flex-end; justify-content: space-between;
    margin-bottom: var(--space-6);
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

  .placeholder {
    text-align: center;
    color: var(--fg-mute);
    max-width: 460px;
    margin: 0 auto;
    padding: var(--space-12) var(--space-6);
    background: var(--bg-surface);
    border: 1px dashed var(--border-default);
    border-radius: var(--radius-md);
  }
  .icon-wrap {
    width: 56px; height: 56px;
    margin: 0 auto var(--space-4);
    border-radius: var(--radius-md);
    display: flex; align-items: center; justify-content: center;
    color: var(--fg-mute);
    background: var(--bg-elevated);
    border: 1px solid var(--border-subtle);
  }
  h3 {
    font-family: var(--font-heading);
    font-size: var(--fs-md);
    font-weight: 700;
    letter-spacing: var(--tracking-tight);
    color: var(--fg-bright);
    margin: 0 0 var(--space-3);
  }
  p { margin: var(--space-2) 0; font-size: var(--fs-sm); color: var(--fg-soft); }
  strong { color: var(--fg-bright); font-weight: 600; }
  .hint { color: var(--fg-faint); font-size: var(--fs-xs); }
</style>
