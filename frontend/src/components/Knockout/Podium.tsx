import { getTeamByCode } from '../../data/teams';
import { useLocale } from '../../i18n/LocaleContext';
import { type KnockoutTree, getWinner } from '../../lib/knockout';
import type { Match } from '../../types/worldcup';
import { Flag } from '../ui/Flag';

interface PodiumProps {
  tree: KnockoutTree;
}

function loserOf(match: Match | null | undefined): string | null {
  if (!match) return null;
  const winner = getWinner(match);
  if (!winner) return null;
  return winner === match.team1 ? match.team2 : match.team1;
}

function formatScore(
  match: Match,
  penaltiesLabel: string,
  extraTimeLabel: string,
): string {
  let score = `${match.goals1} × ${match.goals2}`;
  if (match.penalties1 != null && match.penalties2 != null) {
    score += ` (${match.penalties1}-${match.penalties2} ${penaltiesLabel})`;
  } else if (match.extraTime) {
    score += ` (${extraTimeLabel})`;
  }
  return score;
}

export function Podium({ tree }: PodiumProps) {
  const { t } = useLocale();
  const finalMatch = tree.final[0];
  if (!tree.champion || !finalMatch) return null;

  const runnerUp = loserOf(finalMatch);
  const third = tree.third ? getWinner(tree.third) : null;
  const champion = getTeamByCode(tree.champion);

  const places = [
    { code: runnerUp, label: t.knockout.runnerUp, height: 'h-20' },
    { code: tree.champion, label: t.knockout.champion, height: 'h-28' },
    { code: third, label: t.knockout.thirdPlace, height: 'h-14' },
  ];

  return (
    <div className="flex flex-col items-center mb-12">
      <p className="text-white/50 text-xs uppercase tracking-widest mb-2">
        {t.knockout.final} · {finalMatch.stadium}, {finalMatch.city}
      </p>
      <div className="flex items-center gap-3 mb-8 text-white">
        <Flag code={finalMatch.team1} />
        <span className="font-display text-3xl tracking-wider text-copa-gold">
          {formatScore(finalMatch, t.knockout.penalties, t.knockout.extraTime)}
        </span>
        <Flag code={finalMatch.team2} />
      </div>

      <div className="flex items-end gap-3 sm:gap-6">
        {places.map(
          (place) =>
            place.code && (
              <div
                key={place.label}
                className="flex flex-col items-center w-24 sm:w-32"
              >
                <Flag code={place.code} size={96} className="mb-2" />
                <span
                  className={`text-center text-sm font-semibold mb-2 ${
                    place.code === tree.champion
                      ? 'text-copa-gold'
                      : 'text-white/80'
                  }`}
                >
                  {getTeamByCode(place.code)?.name ?? place.code}
                </span>
                <div
                  className={`w-full ${place.height} rounded-t-lg flex items-start justify-center pt-2 ${
                    place.code === tree.champion
                      ? 'bg-copa-gold/20 border border-copa-gold/40'
                      : 'bg-white/5 border border-white/10'
                  }`}
                >
                  <span className="text-[10px] uppercase tracking-widest text-white/60">
                    {place.label}
                  </span>
                </div>
              </div>
            ),
        )}
      </div>

      {champion && (
        <p className="mt-6 font-display text-2xl sm:text-3xl tracking-wider text-copa-gold">
          {champion.name.toUpperCase()} · {t.knockout.champion}
        </p>
      )}
    </div>
  );
}
