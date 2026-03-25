import type { GroupStanding, Match } from '../types/worldcup';

export function calculateStandings(
  teamCodes: string[],
  matches: Match[],
): GroupStanding[] {
  const standings = new Map<string, GroupStanding>();

  for (const code of teamCodes) {
    standings.set(code, {
      teamCode: code,
      played: 0,
      won: 0,
      drawn: 0,
      lost: 0,
      goalsFor: 0,
      goalsAgainst: 0,
      goalDifference: 0,
      points: 0,
    });
  }

  for (const match of matches) {
    if (match.goals1 === null || match.goals2 === null) continue;

    const team1 = standings.get(match.team1);
    const team2 = standings.get(match.team2);
    if (!team1 || !team2) continue;

    team1.played++;
    team2.played++;
    team1.goalsFor += match.goals1;
    team1.goalsAgainst += match.goals2;
    team2.goalsFor += match.goals2;
    team2.goalsAgainst += match.goals1;

    if (match.goals1 > match.goals2) {
      team1.won++;
      team1.points += 3;
      team2.lost++;
    } else if (match.goals1 < match.goals2) {
      team2.won++;
      team2.points += 3;
      team1.lost++;
    } else {
      team1.drawn++;
      team2.drawn++;
      team1.points += 1;
      team2.points += 1;
    }
  }

  for (const s of standings.values()) {
    s.goalDifference = s.goalsFor - s.goalsAgainst;
  }

  return Array.from(standings.values()).sort((a, b) => {
    if (b.points !== a.points) return b.points - a.points;
    if (b.goalDifference !== a.goalDifference)
      return b.goalDifference - a.goalDifference;
    if (b.goalsFor !== a.goalsFor) return b.goalsFor - a.goalsFor;
    return 0;
  });
}
