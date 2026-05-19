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
      version: '1.4.0',
      date: '2026-05-15',
      title: 'Drag-drop installs across drives',
      entries: [
        { kind: 'added', text: 'Drop Wow.exe anywhere on the launcher to bind a realm install.' },
        { kind: 'added', text: 'Cross-drive base installs — game folder can live on a different disk from the launcher cache.' },
        { kind: 'fixed', text: 'Realm dropdown no longer flickers when switching while a sync is in progress.' },
        { kind: 'changed', text: 'News card spacing tightened on narrow windows.' },
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
