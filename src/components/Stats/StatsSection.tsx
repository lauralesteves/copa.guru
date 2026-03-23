import { useState } from 'react';
import { teams } from '../../data/teams';
import { useScrollReveal } from '../../hooks/useScrollReveal';
import { AnimatedCounter } from '../ui/AnimatedCounter';
import { Flag } from '../ui/Flag';
import { RadarChart } from './RadarChart';

// Simulated team stats (0-100 scale)
const teamStats: Record<string, number[]> = {
  BRA: [90, 72, 85, 88, 80, 35],
  ARG: [88, 80, 82, 85, 78, 40],
  FRA: [85, 82, 80, 83, 76, 38],
  ENG: [82, 78, 78, 80, 75, 42],
  ESP: [80, 75, 90, 92, 70, 30],
  GER: [83, 80, 82, 85, 78, 35],
  POR: [85, 72, 78, 80, 82, 38],
  NED: [78, 76, 84, 86, 72, 32],
  BEL: [76, 74, 80, 82, 74, 36],
  CRO: [72, 78, 82, 84, 68, 34],
  URU: [74, 82, 70, 72, 70, 45],
  COL: [76, 70, 76, 78, 74, 40],
  JPN: [70, 72, 78, 80, 68, 28],
  KOR: [68, 74, 72, 74, 66, 32],
  MEX: [72, 68, 74, 76, 70, 38],
  USA: [70, 70, 72, 74, 68, 35],
  MAR: [68, 78, 70, 72, 66, 36],
  SEN: [72, 70, 68, 70, 70, 42],
};

const STAT_LABELS = ['Ataque', 'Defesa', 'Posse', 'Passes', 'Finalizacoes', 'Faltas'];

const selectableTeams = teams.filter((t) => teamStats[t.code]);

export function StatsSection() {
  const [team1Code, setTeam1Code] = useState('BRA');
  const [team2Code, setTeam2Code] = useState('ARG');
  const titleRef = useScrollReveal<HTMLHeadingElement>({ y: 30 });

  const t1 = teams.find((t) => t.code === team1Code);
  const t2 = teams.find((t) => t.code === team2Code);

  const datasets = [
    {
      label: t1?.name ?? team1Code,
      values: teamStats[team1Code] ?? [50, 50, 50, 50, 50, 50],
      color: t1?.primaryColor ?? '#d4a843',
    },
    {
      label: t2?.name ?? team2Code,
      values: teamStats[team2Code] ?? [50, 50, 50, 50, 50, 50],
      color: t2?.primaryColor ?? '#0047a0',
    },
  ];

  return (
    <section className="py-20 px-6">
      <div className="max-w-4xl mx-auto">
        <h2
          ref={titleRef}
          className="font-display text-4xl sm:text-5xl md:text-6xl text-copa-gold text-center tracking-wider mb-4"
        >
          COMPARAR SELECOES
        </h2>
        <p className="text-white/50 text-center mb-10 text-sm uppercase tracking-widest">
          Radar de atributos lado a lado
        </p>

        {/* Team selectors */}
        <div className="flex items-center justify-center gap-6 mb-8">
          <TeamSelector
            value={team1Code}
            onChange={setTeam1Code}
            teams={selectableTeams}
            color={t1?.primaryColor}
          />
          <span className="text-copa-gold font-display text-2xl">VS</span>
          <TeamSelector
            value={team2Code}
            onChange={setTeam2Code}
            teams={selectableTeams}
            color={t2?.primaryColor}
          />
        </div>

        {/* Radar chart */}
        <RadarChart labels={STAT_LABELS} datasets={datasets} />

        {/* Stats summary */}
        <div className="grid grid-cols-3 gap-4 mt-10">
          <div className="bg-white/5 border border-white/10 rounded-xl p-4 text-center">
            <AnimatedCounter
              value={teams.length}
              className="text-3xl font-bold text-copa-gold"
            />
            <p className="text-white/40 text-xs mt-1 uppercase">Selecoes</p>
          </div>
          <div className="bg-white/5 border border-white/10 rounded-xl p-4 text-center">
            <AnimatedCounter
              value={12}
              className="text-3xl font-bold text-white"
            />
            <p className="text-white/40 text-xs mt-1 uppercase">Grupos</p>
          </div>
          <div className="bg-white/5 border border-white/10 rounded-xl p-4 text-center">
            <AnimatedCounter
              value={104}
              className="text-3xl font-bold text-copa-gold"
            />
            <p className="text-white/40 text-xs mt-1 uppercase">Jogos</p>
          </div>
        </div>
      </div>
    </section>
  );
}

function TeamSelector({
  value,
  onChange,
  teams: teamList,
  color,
}: {
  value: string;
  onChange: (v: string) => void;
  teams: typeof teams;
  color?: string;
}) {
  return (
    <div className="flex items-center gap-2">
      <Flag code={value} size={48} className="rounded" />
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="bg-white/5 border border-white/10 text-white text-sm rounded-lg px-3 py-1.5 focus:border-copa-gold focus:outline-none appearance-none cursor-pointer"
        style={{ borderColor: color ? `${color}40` : undefined }}
      >
        {teamList.map((t) => (
          <option key={t.code} value={t.code} className="bg-copa-dark text-white">
            {t.name}
          </option>
        ))}
      </select>
    </div>
  );
}
