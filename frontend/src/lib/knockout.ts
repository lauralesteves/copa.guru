import type { Match, Stage } from '../types/worldcup';

export type KnockoutStageName = Exclude<Stage, 'group' | 'third'>;

export interface KnockoutTree {
  round32: Match[];
  round16: Match[];
  quarter: Match[];
  semi: Match[];
  final: Match[];
  third: Match | null;
  champion: string | null;
}

const ROUNDS: KnockoutStageName[] = [
  'round32',
  'round16',
  'quarter',
  'semi',
  'final',
];

export function getWinner(match: Match): string | null {
  if (match.goals1 === null || match.goals2 === null) return null;
  if (match.goals1 > match.goals2) return match.team1;
  if (match.goals1 < match.goals2) return match.team2;
  const p1 = match.penalties1 ?? null;
  const p2 = match.penalties2 ?? null;
  if (p1 === null || p2 === null || p1 === p2) return null;
  return p1 > p2 ? match.team1 : match.team2;
}

// Orders each round so that matches i*2 and i*2+1 feed match i of the next
// round, walking the bracket backwards from the final through each team's
// previous knockout match.
export function buildKnockoutTree(matches: Match[]): KnockoutTree {
  const byStage = new Map<KnockoutStageName, Match[]>();
  for (const stage of ROUNDS) {
    byStage.set(
      stage,
      matches.filter((m) => m.stage === stage).sort((a, b) => a.id - b.id),
    );
  }

  const ordered = new Map<KnockoutStageName, Match[]>();
  ordered.set('final', byStage.get('final') ?? []);

  for (let r = ROUNDS.length - 1; r > 0; r--) {
    const next = ordered.get(ROUNDS[r]) ?? [];
    const pool = byStage.get(ROUNDS[r - 1]) ?? [];
    const used = new Set<number>();
    const result: Match[] = [];

    for (const m of next) {
      for (const team of [m.team1, m.team2]) {
        const feeder = pool.find(
          (p) => !used.has(p.id) && getWinner(p) === team,
        );
        if (feeder) {
          used.add(feeder.id);
          result.push(feeder);
        }
      }
    }
    // Unresolved feeders keep FIFA match order
    for (const p of pool) {
      if (!used.has(p.id)) result.push(p);
    }
    ordered.set(ROUNDS[r - 1], result);
  }

  const finalMatch = ordered.get('final')?.[0];

  return {
    round32: ordered.get('round32') ?? [],
    round16: ordered.get('round16') ?? [],
    quarter: ordered.get('quarter') ?? [],
    semi: ordered.get('semi') ?? [],
    final: ordered.get('final') ?? [],
    third: matches.find((m) => m.stage === 'third') ?? null,
    champion: finalMatch ? getWinner(finalMatch) : null,
  };
}
