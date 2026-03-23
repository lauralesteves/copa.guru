import { useRef, useState } from 'react';
import { getTeamByCode } from '../../data/teams';
import type { Match } from '../../types/worldcup';
import { Flag } from '../ui/Flag';
import { GoalTimeline } from './GoalTimeline';

interface MatchCardProps {
  match: Match;
  onScoreChange: (
    matchId: number,
    goals1: number | null,
    goals2: number | null,
  ) => void;
}

function ScoreInput({
  value,
  onChange,
  teamCode,
}: {
  value: number | null;
  onChange: (v: number | null) => void;
  teamCode: string;
}) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [animating, setAnimating] = useState(false);

  function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    const raw = e.target.value;
    if (raw === '') {
      onChange(null);
      return;
    }
    const n = Number.parseInt(raw, 10);
    if (!Number.isNaN(n) && n >= 0 && n <= 99) {
      setAnimating(true);
      onChange(n);
      setTimeout(() => setAnimating(false), 300);
    }
  }

  function increment() {
    const next = (value ?? 0) + 1;
    if (next <= 99) {
      setAnimating(true);
      onChange(next);
      setTimeout(() => setAnimating(false), 300);
    }
  }

  function decrement() {
    if (value !== null && value > 0) {
      setAnimating(true);
      onChange(value - 1);
      setTimeout(() => setAnimating(false), 300);
    }
  }

  return (
    <div className="flex flex-col items-center gap-1">
      <button
        type="button"
        onClick={increment}
        className="text-white/30 hover:text-copa-gold transition-colors text-xs leading-none"
        aria-label={`Incrementar gols ${teamCode}`}
      >
        <svg className="w-3 h-3" viewBox="0 0 12 12" fill="currentColor" aria-hidden="true">
          <path d="M6 2l4 5H2z" />
        </svg>
      </button>

      <div className="relative">
        <input
          ref={inputRef}
          type="text"
          inputMode="numeric"
          value={value ?? ''}
          onChange={handleChange}
          placeholder="-"
          className={`w-10 h-10 text-center text-lg font-bold bg-white/5 border border-white/10 rounded-lg text-white placeholder-white/20 focus:border-copa-gold focus:outline-none focus:ring-1 focus:ring-copa-gold/50 transition-all ${
            animating ? 'scale-110' : 'scale-100'
          }`}
          style={{ transition: 'transform 0.15s ease-out, border-color 0.2s' }}
          aria-label={`Gols ${teamCode}`}
        />
      </div>

      <button
        type="button"
        onClick={decrement}
        className="text-white/30 hover:text-copa-gold transition-colors text-xs leading-none"
        aria-label={`Decrementar gols ${teamCode}`}
      >
        <svg className="w-3 h-3" viewBox="0 0 12 12" fill="currentColor" aria-hidden="true">
          <path d="M6 10l4-5H2z" />
        </svg>
      </button>
    </div>
  );
}

export function MatchCard({ match, onScoreChange }: MatchCardProps) {
  const team1 = getTeamByCode(match.team1);
  const team2 = getTeamByCode(match.team2);

  const hasResult = match.goals1 !== null && match.goals2 !== null;
  let resultClass = '';
  if (hasResult) {
    if (match.goals1! > match.goals2!) resultClass = 'team1-win';
    else if (match.goals1! < match.goals2!) resultClass = 'team2-win';
    else resultClass = 'draw';
  }

  return (
    <div
      className={`py-2 px-3 rounded-lg transition-all ${
        hasResult
          ? 'bg-white/5 border border-white/10'
          : 'bg-transparent border border-transparent hover:border-white/5'
      }`}
    >
      <div className="flex items-center gap-2">
      {/* Team 1 */}
      <div className="flex items-center gap-2 flex-1 justify-end min-w-0">
        <span
          className={`text-xs font-medium truncate ${
            resultClass === 'team1-win'
              ? 'text-copa-gold'
              : resultClass === 'team2-win'
                ? 'text-white/40'
                : 'text-white/80'
          }`}
        >
          {team1?.name ?? match.team1}
        </span>
        <Flag code={match.team1} size={48} className="rounded-sm shrink-0" />
      </div>

      {/* Score inputs */}
      <div className="flex items-center gap-1.5 shrink-0">
        <ScoreInput
          value={match.goals1}
          onChange={(v) => onScoreChange(match.id, v, match.goals2)}
          teamCode={match.team1}
        />
        <span className="text-white/20 text-xs font-bold">x</span>
        <ScoreInput
          value={match.goals2}
          onChange={(v) => onScoreChange(match.id, match.goals1, v)}
          teamCode={match.team2}
        />
      </div>

      {/* Team 2 */}
      <div className="flex items-center gap-2 flex-1 min-w-0">
        <Flag code={match.team2} size={48} className="rounded-sm shrink-0" />
        <span
          className={`text-xs font-medium truncate ${
            resultClass === 'team2-win'
              ? 'text-copa-gold'
              : resultClass === 'team1-win'
                ? 'text-white/40'
                : 'text-white/80'
          }`}
        >
          {team2?.name ?? match.team2}
        </span>
      </div>
      </div>
      <GoalTimeline match={match} />
    </div>
  );
}
