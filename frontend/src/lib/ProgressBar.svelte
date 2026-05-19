<script lang="ts">
  import { overallPct, currentFile, bytesPerSec, humanBytes } from '../stores';
</script>

<div class="wrap">
  <div class="meta">
    <span class="file">{$currentFile || ''}</span>
    <span class="speed">{$bytesPerSec ? `${humanBytes($bytesPerSec)}/s` : ''}</span>
  </div>
  <div class="bar">
    <div class="fill" style="width: {$overallPct}%"></div>
  </div>
  <div class="pct">{$overallPct.toFixed(1)}%</div>
</div>

<style>
  .wrap { display: flex; flex-direction: column; gap: 0.35rem; }
  .meta { display: flex; justify-content: space-between; font-size: 0.8rem; color: #aaa; }
  .file { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 70%; }
  .bar { height: 8px; background: #2a2a33; border-radius: 4px; overflow: hidden; }
  .fill { height: 100%; background: linear-gradient(90deg, #c9a227, #ffd966); transition: width 200ms ease-out; }
  .pct { font-size: 0.75rem; color: #888; text-align: right; }
</style>
