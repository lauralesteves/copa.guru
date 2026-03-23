import { useState } from 'react';
import { teams } from '../../data/teams';
import { useScrollReveal } from '../../hooks/useScrollReveal';
import type { Prediction } from '../../hooks/usePredictions';
import { Flag } from '../ui/Flag';
import { Confetti } from './Confetti';
import { ShareCard } from './ShareCard';

interface PredictionSectionProps {
  predictions: Prediction[];
  champion: string | null;
  setChampion: (code: string | null) => void;
  stats: { total: number; filled: number; percentage: number };
  clearAll: () => void;
}

export function PredictionSection({
  predictions,
  champion,
  setChampion,
  stats,
  clearAll,
}: PredictionSectionProps) {
  const [showConfetti, setShowConfetti] = useState(false);
  const [showShareCard, setShowShareCard] = useState(false);
  const [confirmClear, setConfirmClear] = useState(false);

  const titleRef = useScrollReveal<HTMLHeadingElement>({ y: 30 });
  const contentRef = useScrollReveal<HTMLDivElement>({ y: 40, delay: 0.2 });

  function handleChampionSelect(code: string) {
    setChampion(code);
    setShowConfetti(true);
  }

  return (
    <section id="palpites" className="py-20 px-6">
      <div className="max-w-4xl mx-auto">
        <h2
          ref={titleRef}
          className="font-display text-4xl sm:text-5xl md:text-6xl text-copa-gold text-center tracking-wider mb-4"
        >
          PALPITES
        </h2>
        <p className="text-white text-center mb-12 text-sm uppercase tracking-widest">
          Faça seus palpites e compartilhe
        </p>

        <div ref={contentRef} className="space-y-8">
          {/* Progress */}
          <div className="bg-white/5 border border-white/10 rounded-xl p-6">
            <div className="flex items-center justify-between mb-3">
              <span className="text-sm text-white/60">Progresso dos palpites</span>
              <span className="text-sm font-bold text-copa-gold">
                {stats.filled}/{stats.total} jogos
              </span>
            </div>
            <div className="w-full h-2 bg-white/5 rounded-full overflow-hidden">
              <div
                className="h-full bg-gradient-to-r from-copa-gold to-copa-gold-light rounded-full transition-all duration-500"
                style={{ width: `${stats.percentage}%` }}
              />
            </div>
            <p className="text-xs text-white mt-2">
              <a href="#grupos" className="text-copa-gold hover:text-copa-gold-light transition-colors underline underline-offset-2">Ir para Fase de Grupos</a> e inserir placares na aba "Jogos" de cada grupo
            </p>
          </div>

          {/* Champion pick */}
          <div className="bg-white/5 border border-white/10 rounded-xl p-6">
            <h3 className="font-display text-xl text-copa-gold tracking-wider mb-4">
              QUEM SERÁ O CAMPEÃO?
            </h3>

            {champion ? (
              <div className="flex items-center gap-4">
                <div className="flex items-center gap-3 bg-copa-gold/10 border border-copa-gold/30 rounded-xl px-5 py-3 flex-1">
                  <Flag code={champion} size={96} className="rounded" />
                  <div>
                    <p className="text-copa-gold font-bold text-lg">
                      {teams.find((t) => t.code === champion)?.name ?? champion}
                    </p>
                    <p className="text-white/60 text-xs uppercase">
                      Seu palpite de campeão
                    </p>
                  </div>
                </div>
                <button
                  type="button"
                  onClick={() => setChampion(null)}
                  className="text-white/30 hover:text-white/60 text-xs transition-colors"
                >
                  Trocar
                </button>
              </div>
            ) : (
              <div className="grid grid-cols-4 sm:grid-cols-6 md:grid-cols-8 gap-2">
                {teams.map((team) => (
                  <button
                    key={team.code}
                    type="button"
                    onClick={() => handleChampionSelect(team.code)}
                    className="flex flex-col items-center gap-1 p-2 rounded-lg bg-white/3 border border-transparent hover:border-copa-gold/30 hover:bg-copa-gold/5 transition-all"
                  >
                    <Flag
                      code={team.code}
                      size={48}
                      className="rounded-sm"
                    />
                    <span className="text-[9px] text-white/60 truncate w-full text-center">
                      {team.name}
                    </span>
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Actions */}
          <div className="flex flex-col sm:flex-row gap-3">
            <button
              type="button"
              onClick={() => setShowShareCard(true)}
              disabled={!champion && stats.filled === 0}
              className="flex-1 bg-copa-gold text-copa-dark font-semibold py-3 rounded-lg hover:bg-copa-gold-light transition-colors text-sm uppercase tracking-wider disabled:opacity-30 disabled:cursor-not-allowed"
            >
              Compartilhar Palpites
            </button>
            <button
              type="button"
              onClick={() => {
                if (confirmClear) {
                  clearAll();
                  setConfirmClear(false);
                } else {
                  setConfirmClear(true);
                  setTimeout(() => setConfirmClear(false), 3000);
                }
              }}
              className="px-6 py-3 border border-white/10 text-white/60 rounded-lg hover:border-red-500/30 hover:text-red-400 transition-colors text-sm uppercase tracking-wider"
            >
              {confirmClear ? 'Confirmar reset?' : 'Limpar tudo'}
            </button>
          </div>
        </div>
      </div>

      <Confetti active={showConfetti} onComplete={() => setShowConfetti(false)} />

      {showShareCard && (
        <ShareCard
          champion={champion}
          predictions={predictions}
          stats={stats}
          onClose={() => setShowShareCard(false)}
        />
      )}
    </section>
  );
}
