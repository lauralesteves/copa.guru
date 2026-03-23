import type { Team } from '../types/worldcup';

// In-memory cache populated by fetchWorldCupData()
let teamsCache: Team[] = [];

export function setTeams(data: Team[]) {
  teamsCache = data;
}

export function getTeams(): Team[] {
  return teamsCache;
}

export function getTeamByCode(code: string): Team | undefined {
  return teamsCache.find((t) => t.code === code);
}

export function getTeamsByGroup(group: string): Team[] {
  return teamsCache.filter((t) => t.group === group);
}

export function getTeamsByConfederation(conf: string): Team[] {
  return teamsCache.filter((t) => t.confederation === conf);
}
