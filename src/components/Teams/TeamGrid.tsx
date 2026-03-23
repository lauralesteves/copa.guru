import { useEffect, useRef, useState } from 'react';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';
import { getTeams } from '../../data/teams';
import type { Confederation } from '../../types/worldcup';
import { useScrollReveal } from '../../hooks/useScrollReveal';
import { TeamCard } from './TeamCard';

gsap.registerPlugin(ScrollTrigger);

const confederations: { label: string; value: Confederation | 'ALL' }[] = [
  { label: 'Todas', value: 'ALL' },
  { label: 'UEFA', value: 'UEFA' },
  { label: 'CONMEBOL', value: 'CONMEBOL' },
  { label: 'CAF', value: 'CAF' },
  { label: 'AFC', value: 'AFC' },
  { label: 'CONCACAF', value: 'CONCACAF' },
  { label: 'OFC', value: 'OFC' },
];

export function TeamGrid() {
  const [filter, setFilter] = useState<Confederation | 'ALL'>('ALL');
  const gridRef = useRef<HTMLDivElement>(null);
  const titleRef = useScrollReveal<HTMLHeadingElement>({ y: 30 });
  const subtitleRef = useScrollReveal<HTMLParagraphElement>({
    y: 20,
    delay: 0.15,
  });

  const filtered =
    filter === 'ALL' ? getTeams() : getTeams().filter((t) => t.confederation === filter);

  useEffect(() => {
    if (!gridRef.current) return;
    const children = gridRef.current.children;
    gsap.fromTo(
      children,
      { y: 20, opacity: 0, scale: 0.95 },
      {
        y: 0,
        opacity: 1,
        scale: 1,
        duration: 0.4,
        stagger: 0.03,
        ease: 'power2.out',
      },
    );
  }, [filtered]);

  return (
    <section id="selecoes" className="py-20 px-6">
      <div className="max-w-7xl mx-auto">
        <h2
          ref={titleRef}
          className="font-display text-4xl sm:text-5xl md:text-6xl text-copa-gold text-center tracking-wider mb-4"
        >
          SELEÇÕES
        </h2>
        <p
          ref={subtitleRef}
          className="text-white text-center mb-8 text-sm uppercase tracking-widest"
        >
          {getTeams().length} seleções classificadas
        </p>

        {/* Filter pills */}
        <div className="flex flex-wrap justify-center gap-2 mb-10">
          {confederations.map((conf) => (
            <button
              key={conf.value}
              type="button"
              onClick={() => setFilter(conf.value)}
              className={`px-4 py-1.5 rounded-full text-xs font-semibold uppercase tracking-wider transition-all ${
                filter === conf.value
                  ? 'bg-copa-gold text-copa-dark'
                  : 'bg-white/5 text-white/60 border border-white/10 hover:border-copa-gold/30 hover:text-white'
              }`}
            >
              {conf.label}
            </button>
          ))}
        </div>

        {/* Grid */}
        <div
          ref={gridRef}
          className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4"
        >
          {filtered.map((team) => (
            <TeamCard key={team.code} team={team} />
          ))}
        </div>
      </div>
    </section>
  );
}
