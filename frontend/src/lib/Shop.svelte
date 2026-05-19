<script lang="ts">
  import { selectedServerId, servers } from '../stores';
  const srv = $derived($servers.find(s => s.id === $selectedServerId));
</script>

<section class="shop" aria-label="Shop">
  <div class="placeholder">
    <svg viewBox="0 0 24 24" width="40" height="40" fill="none" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
      <path d="M3 7h18l-2 12H5L3 7z"/>
      <path d="M8 7V5a4 4 0 0 1 8 0v2"/>
    </svg>
    <h3>Shop</h3>
    <p>Cash-shop integration for {srv?.name ?? 'this realm'} not yet wired.</p>
    <p class="hint">Operator can supply a JSON catalog endpoint to populate this view.</p>
  </div>
</section>

<style>
  .shop { display: flex; align-items: center; justify-content: center; min-height: 320px; }
  .placeholder {
    text-align: center;
    color: var(--fg-mute);
    max-width: 380px;
    padding: var(--space-8);
  }
  .placeholder svg { color: var(--accent); opacity: 0.6; margin-bottom: var(--space-3); }
  h3 {
    font-family: var(--font-display);
    font-size: var(--fs-lg);
    letter-spacing: 0.22em;
    text-transform: uppercase;
    color: var(--fg-bright);
    margin: 0 0 var(--space-3);
  }
  p { margin: var(--space-2) 0; font-size: var(--fs-sm); }
  .hint { font-style: italic; color: var(--fg-mute); font-size: var(--fs-xs); }
</style>
