import { writable } from 'svelte/store';
import type { config } from '../wailsjs/go/models';

export type Phase = 'idle' | 'detecting' | 'syncing' | 'launching' | 'error';

export const phase = writable<Phase>('idle');
export const statusMsg = writable<string>('');
export const currentFile = writable<string>('');
export const overallPct = writable<number>(0);
export const bytesPerSec = writable<number>(0);
export const servers = writable<config.Server[]>([]);
export const selectedServerId = writable<string>('');
export const detectedInstalls = writable<Array<{ root: string; locale: string }>>([]);
export const errorMsg = writable<string>('');
export const includeOptional = writable<boolean>(false);

export function humanBytes(n: number): string {
  if (n < 1024) return `${n} B`;
  const units = ['KB', 'MB', 'GB', 'TB'];
  let v = n / 1024;
  let i = 0;
  while (v >= 1024 && i < units.length - 1) { v /= 1024; i++; }
  return `${v.toFixed(1)} ${units[i]}`;
}
