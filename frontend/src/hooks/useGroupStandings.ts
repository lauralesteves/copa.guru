import { useMemo } from 'react';
import type { GroupStanding, Match } from '../types/worldcup';
import { calculateStandings } from '../lib/calculations';

export function useGroupStandings(
  teamCodes: string[],
  matches: Match[],
): GroupStanding[] {
  return useMemo(
    () => calculateStandings(teamCodes, matches),
    [teamCodes, matches],
  );
}
