import { useCallback, useMemo, useState } from 'react';
import { groups } from '../data/groups';
import { createInitialMatches } from '../data/matches';
import { calculateStandings } from '../lib/calculations';
import type { GroupName, GroupStanding, Match } from '../types/worldcup';

export function useWorldCupData() {
  const [matches, setMatches] = useState<Match[]>(() => createInitialMatches());

  const updateScore = useCallback(
    (matchId: number, goals1: number | null, goals2: number | null) => {
      setMatches((prev) =>
        prev.map((m) => (m.id === matchId ? { ...m, goals1, goals2 } : m)),
      );
    },
    [],
  );

  const getGroupMatches = useCallback(
    (groupName: GroupName): Match[] => {
      return matches.filter((m) => m.group === groupName);
    },
    [matches],
  );

  const getGroupStandings = useCallback(
    (groupName: GroupName): GroupStanding[] => {
      const group = groups.find((g) => g.name === groupName);
      if (!group) return [];
      const validTeams = group.teams.filter((t) => !t.startsWith('TBD'));
      const groupMatches = matches.filter((m) => m.group === groupName);
      return calculateStandings(validTeams, groupMatches);
    },
    [matches],
  );

  const allGroupStandings = useMemo(() => {
    const map = new Map<GroupName, GroupStanding[]>();
    for (const group of groups) {
      map.set(group.name, getGroupStandings(group.name));
    }
    return map;
  }, [getGroupStandings]);

  return {
    matches,
    updateScore,
    getGroupMatches,
    getGroupStandings,
    allGroupStandings,
  };
}
