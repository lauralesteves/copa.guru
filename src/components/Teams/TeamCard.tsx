import { useState } from 'react';
import type { Team } from '../../types/worldcup';
import { Flag } from '../ui/Flag';

interface TeamCardProps {
  team: Team;
}

export function TeamCard({ team }: TeamCardProps) {
  const [flipped, setFlipped] = useState(false);

  return (
    <div
      className="perspective-1000 cursor-pointer"
      onClick={() => setFlipped(!flipped)}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') setFlipped(!flipped);
      }}
      role="button"
      tabIndex={0}
    >
      <div
        className={`relative w-full aspect-[4/3] transition-transform duration-500 preserve-3d ${
          flipped ? 'rotate-y-180' : ''
        }`}
        style={{ transformStyle: 'preserve-3d' }}
      >
        {/* Front */}
        <div
          className="absolute inset-0 bg-white/5 border border-white/10 rounded-xl flex flex-col items-center justify-center gap-2 hover:border-copa-gold/30 hover:-translate-y-1 hover:shadow-lg hover:shadow-copa-gold/5 transition-all"
          style={{ backfaceVisibility: 'hidden' }}
        >
          <Flag code={team.code} size={96} className="rounded" />
          <span className="text-white text-sm font-semibold">{team.name}</span>
          <span className="text-white/40 text-[10px] uppercase tracking-widest">
            Grupo {team.group}
          </span>
          {team.isHost && (
            <span className="text-[9px] bg-copa-gold/20 text-copa-gold px-2 py-0.5 rounded-full">
              SEDE
            </span>
          )}
          {team.isDebutant && (
            <span className="text-[9px] bg-green-500/20 text-green-400 px-2 py-0.5 rounded-full">
              ESTREANTE
            </span>
          )}
        </div>

        {/* Back */}
        <div
          className="absolute inset-0 bg-gradient-to-br from-copa-purple-light to-copa-purple border border-white/10 rounded-xl flex flex-col items-center justify-center gap-1.5 p-3 rotate-y-180"
          style={{ backfaceVisibility: 'hidden', transform: 'rotateY(180deg)' }}
        >
          <span className="text-copa-gold font-display text-lg tracking-wider">
            {team.code}
          </span>
          <div className="flex gap-1">
            <div
              className="w-4 h-4 rounded-full border border-white/20"
              style={{ backgroundColor: team.primaryColor }}
            />
            <div
              className="w-4 h-4 rounded-full border border-white/20"
              style={{ backgroundColor: team.secondaryColor }}
            />
          </div>
          <span className="text-white/60 text-[10px] uppercase">
            {team.confederation}
          </span>
          <span className="text-white/40 text-[10px]">
            Grupo {team.group}
          </span>
        </div>
      </div>
    </div>
  );
}
