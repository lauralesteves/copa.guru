import type { Group } from '../types/worldcup';

export const groups: Group[] = [
  { name: 'A', teams: ['MEX', 'KOR', 'RSA', 'TBD1'] },
  { name: 'B', teams: ['CAN', 'SUI', 'QAT', 'TBD2'] },
  { name: 'C', teams: ['BRA', 'MAR', 'SCO', 'HAI'] },
  { name: 'D', teams: ['USA', 'PAR', 'AUS', 'TBD3'] },
  { name: 'E', teams: ['GER', 'ECU', 'CIV', 'CUW'] },
  { name: 'F', teams: ['NED', 'JPN', 'TUN', 'TBD4'] },
  { name: 'G', teams: ['BEL', 'IRN', 'EGY', 'NZL'] },
  { name: 'H', teams: ['ESP', 'URU', 'KSA', 'CPV'] },
  { name: 'I', teams: ['FRA', 'SEN', 'NOR', 'TBD5'] },
  { name: 'J', teams: ['ARG', 'AUT', 'ALG', 'JOR'] },
  { name: 'K', teams: ['POR', 'COL', 'UZB', 'TBD6'] },
  { name: 'L', teams: ['ENG', 'CRO', 'PAN', 'GHA'] },
];

export function getGroupByName(name: string): Group | undefined {
  return groups.find((g) => g.name === name);
}
