import { groups } from '../../data/groups';
import { useScrollReveal, useStaggerReveal } from '../../hooks/useScrollReveal';
import type { GroupStanding, Match } from '../../types/worldcup';
import type { GroupName } from '../../types/worldcup';
import { GroupCard } from './GroupCard';

interface GroupStageProps {
  getGroupMatches: (group: GroupName) => Match[];
  allGroupStandings: Map<GroupName, GroupStanding[]>;
  onScoreChange: (
    matchId: number,
    goals1: number | null,
    goals2: number | null,
  ) => void;
}

export function GroupStage({
  getGroupMatches,
  allGroupStandings,
  onScoreChange,
}: GroupStageProps) {
  const titleRef = useScrollReveal<HTMLHeadingElement>({ y: 30 });
  const subtitleRef = useScrollReveal<HTMLParagraphElement>({
    y: 20,
    delay: 0.15,
  });
  const gridRef = useStaggerReveal<HTMLDivElement>({
    y: 40,
    stagger: 0.08,
    duration: 0.7,
  });

  return (
    <section id="grupos" className="py-20 px-6">
      <div className="max-w-7xl mx-auto">
        <h2
          ref={titleRef}
          className="font-display text-4xl sm:text-5xl md:text-6xl text-copa-gold text-center tracking-wider mb-4"
        >
          FASE DE GRUPOS
        </h2>
        <p
          ref={subtitleRef}
          className="text-white text-center mb-4 text-sm uppercase tracking-widest"
        >
          12 grupos · 48 seleções · Insira seus palpites
        </p>
        <div className="flex flex-wrap justify-center gap-x-6 gap-y-1 mb-12 text-xs text-white/60">
          <span><strong className="text-white/60">Tabela</strong> — classificação ao vivo</span>
          <span><strong className="text-white/60">Jogos</strong> — insira seus placares</span>
          <span><strong className="text-white/60">Gráfico</strong> — diagrama de confrontos</span>
          <span><strong className="text-white/60">3D</strong> — visualização orbital</span>
        </div>

        <div
          ref={gridRef}
          className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6"
        >
          {groups.map((group) => (
            <GroupCard
              key={group.name}
              group={group}
              matches={getGroupMatches(group.name)}
              standings={allGroupStandings.get(group.name) ?? []}
              onScoreChange={onScoreChange}
            />
          ))}
        </div>
      </div>
    </section>
  );
}
