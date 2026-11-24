import { Suspense, lazy, useMemo, useState } from 'react';
import { useScrollReveal } from '../../hooks/useScrollReveal';
import { useLocale } from '../../i18n/LocaleContext';
import { buildKnockoutTree } from '../../lib/knockout';
import type { Match } from '../../types/worldcup';
import { Bracket } from './Bracket';
import { MoleculeGraph } from './MoleculeGraph';
import { Podium } from './Podium';

const Bracket3D = lazy(() =>
  import('./Bracket3D').then((m) => ({ default: m.Bracket3D })),
);

interface KnockoutStageProps {
  matches: Match[];
}

type View = 'bracket' | 'molecule' | '3d';

export function KnockoutStage({ matches }: KnockoutStageProps) {
  const { t } = useLocale();
  const tree = useMemo(() => buildKnockoutTree(matches), [matches]);
  const [view, setView] = useState<View>('molecule');
  const titleRef = useScrollReveal<HTMLHeadingElement>({ y: 30 });
  const subtitleRef = useScrollReveal<HTMLParagraphElement>({
    y: 20,
    delay: 0.15,
  });

  return (
    <section id="mata-mata" className="py-20 px-6">
      <div className="max-w-[1400px] mx-auto">
        <h2
          ref={titleRef}
          className="font-display text-4xl sm:text-5xl md:text-6xl text-copa-gold text-center tracking-wider mb-4"
        >
          {t.knockout.title}
        </h2>
        <p
          ref={subtitleRef}
          className="text-white text-center mb-8 text-sm uppercase tracking-widest"
        >
          {t.knockout.subtitle}
        </p>

        <Podium tree={tree} />

        {/* View toggle */}
        <div className="flex justify-center gap-2 mb-8">
          <button
            type="button"
            onClick={() => setView('bracket')}
            className={`px-4 py-1.5 rounded-full text-xs font-semibold uppercase tracking-wider transition-all ${
              view === 'bracket'
                ? 'bg-copa-gold text-copa-dark'
                : 'bg-white/5 text-white/60 border border-white/10 hover:border-copa-gold/30'
            }`}
          >
            {t.knockout.bracket}
          </button>
          <button
            type="button"
            onClick={() => setView('molecule')}
            className={`px-4 py-1.5 rounded-full text-xs font-semibold uppercase tracking-wider transition-all ${
              view === 'molecule'
                ? 'bg-copa-gold text-copa-dark'
                : 'bg-white/5 text-white/60 border border-white/10 hover:border-copa-gold/30'
            }`}
          >
            {t.knockout.graph}
          </button>
          <button
            type="button"
            onClick={() => setView('3d')}
            className={`px-4 py-1.5 rounded-full text-xs font-semibold uppercase tracking-wider transition-all ${
              view === '3d'
                ? 'bg-copa-gold text-copa-dark'
                : 'bg-white/5 text-white/60 border border-white/10 hover:border-copa-gold/30'
            }`}
          >
            {t.knockout.threeD}
          </button>
        </div>

        {view === 'bracket' && <Bracket tree={tree} />}
        {view === 'molecule' && <MoleculeGraph tree={tree} />}
        {view === '3d' && (
          <Suspense
            fallback={
              <div className="w-full h-[500px] flex items-center justify-center">
                <div className="skeleton w-64 h-64 rounded-xl" />
              </div>
            }
          >
            <Bracket3D tree={tree} />
          </Suspense>
        )}
      </div>
    </section>
  );
}
