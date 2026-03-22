import { useScrollReveal } from '../../hooks/useScrollReveal';
import type { GroupName, GroupStanding } from '../../types/worldcup';
import { Bracket } from './Bracket';

interface KnockoutStageProps {
  allGroupStandings: Map<GroupName, GroupStanding[]>;
}

export function KnockoutStage({ allGroupStandings }: KnockoutStageProps) {
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
          MATA-MATA
        </h2>
        <p
          ref={subtitleRef}
          className="text-white/50 text-center mb-12 text-sm uppercase tracking-widest"
        >
          Dos 32 avos ate a grande final
        </p>

        <Bracket allGroupStandings={allGroupStandings} />
      </div>
    </section>
  );
}
