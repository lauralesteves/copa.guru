export type Confederation =
  | 'AFC'
  | 'CAF'
  | 'CONCACAF'
  | 'CONMEBOL'
  | 'OFC'
  | 'UEFA';

export type GroupName =
  | 'A'
  | 'B'
  | 'C'
  | 'D'
  | 'E'
  | 'F'
  | 'G'
  | 'H'
  | 'I'
  | 'J'
  | 'K'
  | 'L';

export type Stage =
  | 'group'
  | 'round32'
  | 'round16'
  | 'quarter'
  | 'semi'
  | 'third'
  | 'final';

export interface Team {
  name: string;
  code: string;
  confederation: Confederation;
  group: GroupName;
  primaryColor: string;
  secondaryColor: string;
  isHost: boolean;
  isDebutant: boolean;
}

export interface Match {
  id: number;
  team1: string;
  team2: string;
  goals1: number | null;
  goals2: number | null;
  date: string;
  time: string;
  stadium: string;
  city: string;
  group: GroupName | null;
  stage: Stage;
}

export interface GroupStanding {
  teamCode: string;
  played: number;
  won: number;
  drawn: number;
  lost: number;
  goalsFor: number;
  goalsAgainst: number;
  goalDifference: number;
  points: number;
}

export interface Group {
  name: GroupName;
  teams: string[];
}
