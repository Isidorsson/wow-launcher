<script lang="ts">
  import { servers, selectedServerId } from '../stores';

  function select(id: string) {
    selectedServerId.set(id);
  }
</script>

<aside class="sidebar">
  <h2>Realms</h2>
  <ul>
    {#each $servers as srv}
      <li>
        <button
          class:active={$selectedServerId === srv.id}
          on:click={() => select(srv.id)}>
          <span class="name">{srv.name}</span>
          <span class="id">{srv.id}</span>
        </button>
      </li>
    {/each}
  </ul>
</aside>

<style>
  .sidebar {
    width: 220px;
    background: #15151a;
    border-right: 1px solid #2a2a33;
    padding: 1rem 0.5rem;
    overflow-y: auto;
  }
  h2 {
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: #888;
    margin: 0 0.75rem 0.75rem;
  }
  ul { list-style: none; padding: 0; margin: 0; }
  li { margin: 0.25rem 0; }
  button {
    width: 100%;
    text-align: left;
    background: transparent;
    border: 1px solid transparent;
    color: #ddd;
    padding: 0.6rem 0.75rem;
    border-radius: 6px;
    cursor: pointer;
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }
  button:hover { background: #1f1f27; }
  button.active { background: #232331; border-color: #c9a227; }
  .name { font-weight: 600; }
  .id { font-size: 0.7rem; color: #777; }
</style>
