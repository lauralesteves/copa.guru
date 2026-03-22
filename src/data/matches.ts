import type { GroupName, Match } from '../types/worldcup';
import { groups } from './groups';

function generateGroupMatches(
  groupName: GroupName,
  teams: string[],
  startId: number,
): Match[] {
  const valid = teams.filter((t) => !t.startsWith('TBD'));
  const matches: Match[] = [];
  let id = startId;

  for (let i = 0; i < valid.length; i++) {
    for (let j = i + 1; j < valid.length; j++) {
      matches.push({
        id: id++,
        team1: valid[i],
        team2: valid[j],
        goals1: null,
        goals2: null,
        date: '',
        time: '',
        stadium: '',
        city: '',
        group: groupName,
        stage: 'group',
      });
    }
  }

  return matches;
}

export function createInitialMatches(): Match[] {
  const all: Match[] = [];
  let nextId = 1;

  for (const group of groups) {
    const gm = generateGroupMatches(group.name, group.teams, nextId);
    all.push(...gm);
    nextId += gm.length;
  }

  return all;
}
