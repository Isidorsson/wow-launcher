<script lang="ts">
  type EntryKind = 'added' | 'fixed' | 'changed' | 'removed';
  type Entry = { kind: EntryKind; text: string };
  type Release = {
    version: string;
    date: string;
    title?: string;
    entries: Entry[];
  };

  const releases: Release[] = [
    {
      version: '1.5.0-beta',
      date: '2026-05-19',
      title: 'Stress-test fixtures & layout hardening',
      entries: [
        { kind: 'added', text: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat — long-form entry used to verify multi-line wrapping behaviour inside the entry grid does not break baseline alignment with the kind chip.' },
        { kind: 'added', text: 'Pellentesque-habitant-morbi-tristique-senectus-et-netus-et-malesuada-fames-ac-turpis-egestas — single unbroken token to test overflow handling and word-break rules on narrow viewports.' },
        { kind: 'changed', text: 'Short.' },
        { kind: 'fixed', text: 'Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.' },
        { kind: 'removed', text: 'Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum, plus several additional deprecated hooks that nobody had touched since the original prototype shipped two releases ago.' },
      ],
    },
    {
      version: '1.4.0',
      date: '2026-05-15',
      title: 'Drag-drop installs across drives',
      entries: [
        { kind: 'added', text: 'Drop Wow.exe anywhere on the launcher to bind a realm install.' },
        { kind: 'added', text: 'Cross-drive base installs — game folder can live on a different disk from the launcher cache.' },
        { kind: 'fixed', text: 'Realm dropdown no longer flickers when switching while a sync is in progress.' },
        { kind: 'changed', text: 'News card spacing tightened on narrow windows.' },
        { kind: 'added', text: 'Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, eaque ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo.' },
        { kind: 'changed', text: 'Drop-target overlay opacity lowered from 0.85 to 0.6 to keep the underlying news card legible during the hover state.' },
      ],
    },
    {
      version: '1.3.0',
      date: '2026-04-30',
      title: 'Profile materialization',
      entries: [
        { kind: 'added', text: 'Per-realm profile directory with realmlist.wtf auto-written on first launch.' },
        { kind: 'added', text: 'Patch manifest verification before launch — corrupt files re-pulled automatically.' },
        { kind: 'fixed', text: 'Sync progress could stall at 99.9% when the final file was a small MPQ.' },
      ],
    },
    {
      version: '1.2.1',
      date: '2026-04-12',
      title: 'Hotfix — auth handshake',
      entries: [
        { kind: 'fixed', text: 'Login attempts against realms with TLS 1.3-only auth no longer fail with "handshake refused".' },
        { kind: 'fixed', text: 'Settings modal would not close on Escape if a tab transition was mid-animation.' },
      ],
    },
    {
      version: '1.2.0',
      date: '2026-04-02',
      title: 'News-first layout',
      entries: [
        { kind: 'added', text: 'Realm news feed pulled from operator-hosted JSON.' },
        { kind: 'added', text: 'Settings modal with Installation / Downloads / Profile tabs.' },
        { kind: 'changed', text: 'Play strip moved to bottom — progress bar replaces status text during sync.' },
        { kind: 'removed', text: 'Legacy "Setup Card" panel — superseded by the setup-call in the news pane.' },
      ],
    },
    {
      version: '1.1.0',
      date: '2026-03-18',
      title: 'Patcher rewrite — chunked downloads',
      entries: [
        { kind: 'added', text: 'Parallel chunked downloads with resumable transfer state. Bandwidth ceiling configurable per realm; sane defaults pull from system network class.' },
        { kind: 'added', text: 'SHA-256 verification on every patch chunk before commit to the install tree, with automatic retry on mismatch up to three attempts before surfacing the error to the user.' },
        { kind: 'changed', text: 'Patch progress now reported as bytes-transferred rather than file-count, which gives a more honest ETA on mixed-size manifests where a handful of multi-gigabyte MPQs dominate total transfer time.' },
        { kind: 'fixed', text: 'Patcher would occasionally exit with status 0 while leaving the install in a partially-written state — temp files now atomically renamed on success and cleaned on failure.' },
        { kind: 'fixed', text: 'Tooltip clipping on the realm selector in the bottom-right corner when running at sub-1280px widths.' },
        { kind: 'removed', text: 'Single-threaded "classic" downloader codepath behind the deprecated --legacy-patcher flag. The flag is now ignored with a console warning and will be removed entirely in 2.0.' },
      ],
    },
    {
      version: '1.0.2',
      date: '2026-03-05',
      title: 'First-run polish',
      entries: [
        { kind: 'fixed', text: 'First-launch wizard could deadlock on systems with no writable AppData path (rare, surfaced by a handful of corporate-managed Windows installs).' },
        { kind: 'fixed', text: 'Realm list cache TTL respected on cold start instead of always refetching.' },
        { kind: 'changed', text: 'Default install root suggestion now respects the largest fixed-disk free-space heuristic rather than always picking C:.' },
      ],
    },
    {
      version: '1.0.1',
      date: '2026-02-22',
      title: 'Day-one hotfix',
      entries: [
        { kind: 'fixed', text: 'Crash when launching with no realms configured.' },
        { kind: 'fixed', text: 'Magnam aliquam quaerat voluptatem ut enim ad minima veniam, quis nostrum exercitationem ullam corporis suscipit laboriosam, nisi ut aliquid ex ea commodi consequatur — fixed a race in the early bootstrap path that occasionally left the splash screen visible behind the main window.' },
      ],
    },
    {
      version: '1.0.0',
      date: '2026-02-20',
      title: 'Initial release',
      entries: [
        { kind: 'added', text: 'Initial public release of the launcher with realm selection, news feed, patcher, and one-click play.' },
        { kind: 'added', text: 'Neque porro quisquam est, qui dolorem ipsum quia dolor sit amet, consectetur, adipisci velit, sed quia non numquam eius modi tempora incidunt ut labore et dolore magnam aliquam quaerat voluptatem. This entry exists to verify that the very first release in the timeline still renders its rise animation with a sensible delay rather than waiting noticeably long for the staggered cascade to complete.' },
        { kind: 'added', text: 'Dark cathedral theme with gilt rune accents.' },
        { kind: 'added', text: 'Wails + Svelte 5 frontend bundled into a single executable under 18 MB.' },
      ],
    },
  ];

  function kindLabel(k: EntryKind): string {
    return { added: 'Added', fixed: 'Fixed', changed: 'Changed', removed: 'Removed' }[k];
  }
</script>

<section class="changelog">
  <header>
    <h2>Release notes</h2>
  </header>
  <div class="sample-note">Sample data — wire to real release notes when available.</div>

  <ol class="releases">
    {#each releases as r, i}
      <li class="release" style="animation-delay: {i * 70}ms;">
        <div class="ver-row">
          <span class="ver">v{r.version}</span>
          <time>{r.date}</time>
        </div>
        {#if r.title}
          <h3>{r.title}</h3>
        {/if}
        <ul class="entries">
          {#each r.entries as e}
            <li class="entry">
              <span class="kind kind-{e.kind}">{kindLabel(e.kind)}</span>
              <span class="text">{e.text}</span>
            </li>
          {/each}
        </ul>
      </li>
    {/each}
  </ol>
</section>

<style>
  .changelog { padding: 0; }

  header { margin-bottom: var(--space-4); }
  header h2 {
    margin: 0;
    font-family: var(--font-display);
    font-weight: 700;
    font-size: var(--fs-lg);
    letter-spacing: 0.18em;
    text-transform: uppercase;
    color: var(--fg-bright);
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

  .releases { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 0.9rem; }

  .release {
    position: relative;
    background: linear-gradient(180deg, rgba(19, 26, 48, 0.6), rgba(8, 11, 20, 0.78));
    border: 1px solid var(--rune-line);
    padding: 0.9rem 1.1rem 1rem;
    border-radius: 2px;
    box-shadow:
      inset 0 1px 0 rgba(143, 205, 255, 0.04),
      0 4px 14px rgba(0,0,0,0.35);
    opacity: 0;
    transform: translateY(8px);
    animation: rise 420ms cubic-bezier(0.2, 0.8, 0.2, 1) forwards;
  }
  .release::before {
    content: '';
    position: absolute; left: 0; top: 8px; bottom: 8px;
    width: 2px;
    background: linear-gradient(180deg, transparent, var(--gold), transparent);
    opacity: 0.5;
  }
  @keyframes rise {
    to { opacity: 1; transform: translateY(0); }
  }

  .ver-row {
    display: flex; align-items: center; gap: 0.7rem;
    font-family: var(--font-ui);
    font-size: 0.74rem;
  }
  .ver {
    color: var(--gold-bright);
    font-family: var(--font-display);
    font-weight: 700;
    letter-spacing: 0.14em;
    text-transform: uppercase;
    padding: 0.15rem 0.55rem;
    border: 1px solid rgba(78, 164, 255, 0.45);
    background: rgba(78, 164, 255, 0.1);
    border-radius: 2px;
    font-size: 0.72rem;
  }
  time { color: var(--text-mute); font-style: italic; }

  .release h3 {
    margin: 0.45rem 0 0.55rem;
    color: var(--text-bright);
    font-family: var(--font-display);
    font-weight: 600;
    font-size: 1rem;
    letter-spacing: 0.06em;
  }

  .entries { list-style: none; padding: 0; margin: 0.2rem 0 0; display: flex; flex-direction: column; gap: 0.35rem; }
  .entry {
    display: grid;
    grid-template-columns: 76px 1fr;
    gap: 0.7rem;
    align-items: baseline;
    color: var(--text);
    font-size: 0.92rem;
    line-height: 1.5;
  }
  .kind {
    text-transform: uppercase;
    letter-spacing: 0.12em;
    font-weight: 700;
    font-size: 0.64rem;
    padding: 0.12rem 0.4rem;
    border-radius: 2px;
    border: 1px solid var(--rune-line);
    background: rgba(0,0,0,0.4);
    color: var(--text-soft);
    text-align: center;
  }
  .kind-added   { color: var(--fel-glow);    border-color: rgba(76, 175, 80, 0.4);   background: rgba(76, 175, 80, 0.08); }
  .kind-fixed   { color: var(--arcane);      border-color: rgba(106, 169, 216, 0.4); background: rgba(106, 169, 216, 0.08); }
  .kind-changed { color: var(--gold-bright); border-color: rgba(78, 164, 255, 0.5);  background: rgba(78, 164, 255, 0.1); }
  .kind-removed { color: var(--blood-glow);  border-color: rgba(192, 57, 43, 0.4);   background: rgba(192, 57, 43, 0.08); }

  .text { color: var(--text); }
</style>
