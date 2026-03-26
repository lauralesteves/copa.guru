import { useRef } from 'react';
import { getTeams } from '../../data/teams';
import { useLocale } from '../../i18n/LocaleContext';
import type { Prediction } from '../../hooks/usePredictions';
import { Flag } from '../ui/Flag';

interface ShareCardProps {
  champion: string | null;
  predictions: Prediction[];
  stats: { total: number; filled: number; percentage: number };
  onClose: () => void;
}

export function ShareCard({
  champion,
  predictions,
  stats,
  onClose,
}: ShareCardProps) {
  const { t } = useLocale();
  const cardRef = useRef<HTMLDivElement>(null);

  const teams = getTeams();
  const championTeam = champion
    ? teams.find((t) => t.code === champion)
    : null;

  // Count predictions by result type
  const homeWins = predictions.filter((p) => p.goals1 > p.goals2).length;
  const draws = predictions.filter((p) => p.goals1 === p.goals2).length;
  const awayWins = predictions.filter((p) => p.goals1 < p.goals2).length;
  const totalGoals = predictions.reduce(
    (sum, p) => sum + p.goals1 + p.goals2,
    0,
  );

  function handleCopyText() {
    const sl = t.shareCard.shareLines;
    const lines = [
      sl.header,
      '',
      champion ? `${sl.champion}: ${championTeam?.name ?? champion}` : '',
      `${sl.filled}: ${stats.filled}/${stats.total}`,
      `${sl.goals}: ${totalGoals}`,
      sl.results.replace('%h', String(homeWins)).replace('%d', String(draws)).replace('%a', String(awayWins)),
      '',
      sl.cta,
    ]
      .filter(Boolean)
      .join('\n');

    navigator.clipboard.writeText(lines);
  }

  return (
    <div
      className="fixed inset-0 z-[150] flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm"
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose();
      }}
      onKeyDown={(e) => {
        if (e.key === 'Escape') onClose();
      }}
      role="dialog"
      aria-modal="true"
    >
      <div
        ref={cardRef}
        className="bg-gradient-to-br from-copa-purple to-copa-dark border border-copa-gold/20 rounded-2xl p-8 max-w-md w-full shadow-2xl"
      >
        {/* Header */}
        <div className="text-center mb-6">
          <h3 className="font-display text-3xl text-copa-gold tracking-wider">
            COPA.GURU
          </h3>
          <p className="text-white/60 text-xs uppercase tracking-widest mt-1">
            {t.shareCard.title}
          </p>
        </div>

        {/* Champion */}
        {championTeam && (
          <div className="flex items-center justify-center gap-3 bg-copa-gold/10 border border-copa-gold/20 rounded-xl py-4 px-6 mb-6">
            <Flag code={championTeam.code} size={96} className="rounded" />
            <div>
              <p className="text-[10px] text-white/60 uppercase tracking-widest">
                {t.shareCard.champion}
              </p>
              <p className="text-copa-gold font-bold text-xl">
                {championTeam.name}
              </p>
            </div>
          </div>
        )}

        {/* Stats grid */}
        <div className="grid grid-cols-2 gap-3 mb-6">
          <div className="bg-white/5 rounded-lg p-3 text-center">
            <p className="text-2xl font-bold text-white">{stats.filled}</p>
            <p className="text-[10px] text-white/60 uppercase">{t.shareCard.games}</p>
          </div>
          <div className="bg-white/5 rounded-lg p-3 text-center">
            <p className="text-2xl font-bold text-copa-gold">{totalGoals}</p>
            <p className="text-[10px] text-white/60 uppercase">{t.shareCard.totalGoals}</p>
          </div>
          <div className="bg-white/5 rounded-lg p-3 text-center">
            <p className="text-2xl font-bold text-green-400">{homeWins}</p>
            <p className="text-[10px] text-white/60 uppercase">{t.shareCard.homeWins}</p>
          </div>
          <div className="bg-white/5 rounded-lg p-3 text-center">
            <p className="text-2xl font-bold text-blue-400">{draws}</p>
            <p className="text-[10px] text-white/60 uppercase">{t.shareCard.draws}</p>
          </div>
        </div>

        {/* Progress bar */}
        <div className="mb-6">
          <div className="flex justify-between text-xs text-white/60 mb-1">
            <span>{t.shareCard.progress}</span>
            <span>{stats.percentage}%</span>
          </div>
          <div className="w-full h-1.5 bg-white/5 rounded-full">
            <div
              className="h-full bg-copa-gold rounded-full"
              style={{ width: `${stats.percentage}%` }}
            />
          </div>
        </div>

        {/* Actions */}
        <div className="flex gap-3">
          <button
            type="button"
            onClick={handleCopyText}
            className="flex-1 bg-copa-gold text-copa-dark font-semibold py-2.5 rounded-lg hover:bg-copa-gold-light transition-colors text-sm"
          >
            {t.shareCard.copyText}
          </button>
          <button
            type="button"
            onClick={onClose}
            className="px-4 py-2.5 border border-white/10 text-white/60 rounded-lg hover:text-white transition-colors text-sm"
          >
            {t.shareCard.close}
          </button>
        </div>

        {/* Footer */}
        <p className="text-center text-[10px] text-white/20 mt-4">
          {t.footer.tagline}
        </p>
      </div>
    </div>
  );
}
