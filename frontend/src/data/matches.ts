import type { Match } from '../types/worldcup';

let matchesCache: Match[] = [];

export function setMatches(data: Match[]) {
  matchesCache = data;
}

export function getInitialMatches(): Match[] {
  return matchesCache;
}
