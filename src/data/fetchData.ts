import type { Group, Match, Team } from '../types/worldcup';
import { getGroups, setGroups } from './groups';
import { getInitialMatches, setMatches } from './matches';
import { getTeams, setTeams } from './teams';

interface WorldCupData {
  teams: Team[];
  groups: Group[];
  matches: Match[];
}

let loaded = false;

export async function fetchWorldCupData(): Promise<WorldCupData> {
  if (loaded) {
    return {
      teams: getTeams(),
      groups: getGroups(),
      matches: getInitialMatches(),
    };
  }

  const res = await fetch('/data/worldcup2026.json');
  const data: WorldCupData = await res.json();

  setTeams(data.teams);
  setGroups(data.groups);
  setMatches(data.matches);
  loaded = true;

  return data;
}
