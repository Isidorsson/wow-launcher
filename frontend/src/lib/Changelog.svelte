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
        { kind: 'fixed', text: 'Login attempts against realms with TLS 1.3-only auth no longer fail with “handshake refused”.' },
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
        { kind: 'removed', text: 'Legacy “Setup Card” panel — superseded by the setup-call in the news pane.' },
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
        { kind: 'removed', text: 'Single-threaded “classic” downloader codepath behind the deprecated --legacy-patcher flag. The flag is now ignored with a console warning and will be removed entirely in 2.0.' },
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

  const KIND_ORDER: EntryKind[] = ['added', 'changed', 'fixed', 'removed'];
  const KIND_LABEL: Record<EntryKind, string> = {
    added: 'Added',
    changed: 'Changed',
    fixed: 'Fixed',
    removed: 'Removed',
  };

  function groupByKind(entries: Entry[]): { kind: EntryKind; items: string[] }[] {
    const buckets = new Map<EntryKind, string[]>();
    for (const e of entries) {
      const arr = buckets.get(e.kind) ?? [];
      arr.push(e.text);
      buckets.set(e.kind, arr);
    }
    return KIND_ORDER
      .filter(k => buckets.has(k))
      .map(k => ({ kind: k, items: buckets.get(k)! }));
  }
</script>

<section class="changelog">
  <header class="head">
    <div class="head-text">
      <span class="eyebrow">Release history</span>
      <h2>Patch Notes</h2>
    </div>
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

        <div class="groups">
          {#each groupByKind(r.entries) as g}
            <section class="group group-{g.kind}" aria-labelledby="g-{r.version}-{g.kind}">
              <h4 class="group-head" id="g-{r.version}-{g.kind}">
                <span class="ico" aria-hidden="true">
                  {#if g.kind === 'added'}
                    <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M5 12h14"/><path d="M12 5v14"/>
                    </svg>
                  {:else if g.kind === 'fixed'}
                    <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
                    </svg>
                  {:else if g.kind === 'changed'}
                    <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/>
                      <path d="M21 3v5h-5"/>
                      <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/>
                      <path d="M8 16H3v5"/>
                    </svg>
                  {:else}
                    <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M5 12h14"/>
                    </svg>
                  {/if}
                </span>
                <span class="group-label">{KIND_LABEL[g.kind]}</span>
                <span class="group-count" aria-hidden="true">{g.items.length}</span>
              </h4>
              <ul class="entries">
                {#each g.items as text}
                  <li class="entry">{text}</li>
                {/each}
              </ul>
            </section>
          {/each}
        </div>
      </li>
    {/each}
  </ol>
</section>

<style>
  .changelog { padding: 0; }

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
    text-wrap: balance;
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

  .releases { list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: var(--space-4); }

  .release {
    position: relative;
    background: var(--bg-surface);
    border: 1px solid var(--border-subtle);
    padding: var(--space-5) var(--space-6);
    border-radius: var(--radius-md);
    opacity: 0;
    transform: translateY(6px);
    animation: rise 360ms var(--ease-out) forwards;
    transition: border-color var(--dur-fast) var(--ease-out);
  }
  .release:hover { border-color: var(--border-strong); }
  @keyframes rise {
    to { opacity: 1; transform: translateY(0); }
  }
  @media (prefers-reduced-motion: reduce) {
    .release {
      animation: none;
      opacity: 1;
      transform: none;
    }
  }

  .ver-row {
    display: flex; align-items: center; gap: var(--space-3);
    font-family: var(--font-mono);
    font-size: var(--fs-2xs);
  }
  .ver {
    color: var(--fg-bright);
    font-family: var(--font-mono);
    font-weight: 600;
    letter-spacing: 0;
    padding: 0.15rem 0.45rem;
    border: 1px solid var(--border-default);
    background: var(--bg-raised);
    border-radius: var(--radius-sm);
    font-size: var(--fs-2xs);
    font-variant-numeric: tabular-nums;
  }
  time {
    color: var(--fg-faint);
    font-style: normal;
    font-family: var(--font-mono);
    font-size: var(--fs-2xs);
    font-variant-numeric: tabular-nums;
  }

  .release h3 {
    margin: var(--space-3) 0 var(--space-2);
    color: var(--fg-bright);
    font-family: var(--font-heading);
    font-weight: 700;
    font-size: var(--fs-md);
    letter-spacing: var(--tracking-tight);
    text-wrap: balance;
  }

  .groups {
    margin-top: var(--space-3);
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
  }

  .group {
    border-left: 2px solid var(--border-subtle);
    padding-left: var(--space-4);
  }
  .group-added   { border-left-color: var(--c-green-400); }
  .group-fixed   { border-left-color: var(--accent); }
  .group-changed { border-left-color: var(--c-amber-400); }
  .group-removed { border-left-color: var(--status-error); }

  .group-head {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    margin: 0 0 var(--space-2);
    font-family: var(--font-heading);
    font-size: var(--fs-2xs);
    font-weight: 700;
    letter-spacing: var(--tracking-wider);
    text-transform: uppercase;
    color: var(--fg-mute);
    line-height: 1;
  }
  .ico {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 18px;
    height: 18px;
    border-radius: var(--radius-sm);
    flex-shrink: 0;
  }
  .group-added   .ico { color: var(--c-green-400); background: color-mix(in srgb, var(--c-green-400) 12%, transparent); }
  .group-fixed   .ico { color: var(--accent);      background: color-mix(in srgb, var(--accent) 12%, transparent); }
  .group-changed .ico { color: var(--c-amber-400); background: color-mix(in srgb, var(--c-amber-400) 12%, transparent); }
  .group-removed .ico { color: var(--status-error);background: color-mix(in srgb, var(--status-error) 12%, transparent); }

  .group-added   .group-label { color: var(--c-green-400); }
  .group-fixed   .group-label { color: var(--accent); }
  .group-changed .group-label { color: var(--c-amber-400); }
  .group-removed .group-label { color: var(--status-error); }

  .group-count {
    color: var(--fg-faint);
    font-family: var(--font-mono);
    font-size: var(--fs-2xs);
    font-weight: 500;
    letter-spacing: 0;
    text-transform: none;
    font-variant-numeric: tabular-nums;
  }

  .entries {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }
  .entry {
    position: relative;
    padding-left: var(--space-3);
    color: var(--fg-default);
    font-size: var(--fs-sm);
    line-height: 1.55;
    overflow-wrap: anywhere;
  }
  .entry::before {
    content: '';
    position: absolute;
    left: 0;
    top: 0.62em;
    width: 4px;
    height: 4px;
    border-radius: 50%;
    background: var(--fg-faint);
  }
</style>
