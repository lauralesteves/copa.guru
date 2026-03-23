import { useEffect, useRef } from 'react';
import gsap from 'gsap';
import { Flip } from 'gsap/dist/Flip';
import { getTeamByCode } from '../../data/teams';
import type { GroupStanding } from '../../types/worldcup';
import { Flag } from '../ui/Flag';

gsap.registerPlugin(Flip);

interface StandingsTableProps {
  standings: GroupStanding[];
}

export function StandingsTable({ standings }: StandingsTableProps) {
  const tbodyRef = useRef<HTMLTableSectionElement>(null);
  const prevOrderRef = useRef<string[]>([]);

  useEffect(() => {
    const tbody = tbodyRef.current;
    if (!tbody) return;

    const currentOrder = standings.map((s) => s.teamCode);
    const prevOrder = prevOrderRef.current;

    // Only animate if the order actually changed (not on first render)
    if (
      prevOrder.length > 0 &&
      prevOrder.some((code, i) => code !== currentOrder[i])
    ) {
      const state = Flip.getState(tbody.children);

      // Re-order DOM to match new standings
      for (const code of currentOrder) {
        const row = tbody.querySelector(`[data-team="${code}"]`);
        if (row) tbody.appendChild(row);
      }

      Flip.from(state, {
        duration: 0.5,
        ease: 'power2.inOut',
        stagger: 0.05,
        absolute: true,
      });
    }

    prevOrderRef.current = currentOrder;
  }, [standings]);

  return (
    <div className="mt-4 border-t border-white/10 pt-3">
      <table className="w-full text-xs">
        <thead>
          <tr className="text-white/40 uppercase tracking-wider">
            <th className="text-left pb-2 font-medium w-5">#</th>
            <th className="text-left pb-2 font-medium">Seleção</th>
            <th className="text-center pb-2 font-medium w-7">PG</th>
            <th className="text-center pb-2 font-medium w-5">J</th>
            <th className="text-center pb-2 font-medium w-5">V</th>
            <th className="text-center pb-2 font-medium w-5">E</th>
            <th className="text-center pb-2 font-medium w-5">D</th>
            <th className="text-center pb-2 font-medium w-7">SG</th>
          </tr>
        </thead>
        <tbody ref={tbodyRef}>
          {standings.map((s, i) => {
            const team = getTeamByCode(s.teamCode);
            const isQualified = i < 2;

            return (
              <tr
                key={s.teamCode}
                data-team={s.teamCode}
                className={`border-t border-white/5 transition-colors duration-300 hover:bg-white/5 ${
                  isQualified ? 'text-white' : 'text-white/60'
                }`}
                style={{
                  backgroundColor:
                    team && s.points > 0
                      ? `${team.primaryColor}08`
                      : undefined,
                }}
              >
                <td className="py-1.5">
                  <span
                    className={`inline-flex items-center justify-center w-4 h-4 rounded-full text-[10px] font-bold transition-colors duration-300 ${
                      isQualified
                        ? 'bg-copa-gold/20 text-copa-gold'
                        : 'text-white/30'
                    }`}
                  >
                    {i + 1}
                  </span>
                </td>
                <td className="py-1.5">
                  <div className="flex items-center gap-1.5">
                    <Flag
                      code={s.teamCode}
                      size={48}
                      className="rounded-sm w-5 h-3.5"
                    />
                    <span className="font-medium truncate">
                      {team?.name ?? s.teamCode}
                    </span>
                  </div>
                </td>
                <td className="py-1.5 text-center font-bold text-copa-gold">
                  {s.points}
                </td>
                <td className="py-1.5 text-center">{s.played}</td>
                <td className="py-1.5 text-center">{s.won}</td>
                <td className="py-1.5 text-center">{s.drawn}</td>
                <td className="py-1.5 text-center">{s.lost}</td>
                <td className="py-1.5 text-center">
                  <span
                    className={
                      s.goalDifference > 0
                        ? 'text-green-400'
                        : s.goalDifference < 0
                          ? 'text-red-400'
                          : ''
                    }
                  >
                    {s.goalDifference > 0 ? '+' : ''}
                    {s.goalDifference}
                  </span>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
