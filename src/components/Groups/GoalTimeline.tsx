import { useEffect, useRef } from 'react';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';
import { getTeamByCode } from '../../data/teams';
import type { Match } from '../../types/worldcup';

gsap.registerPlugin(ScrollTrigger);

interface GoalTimelineProps {
  match: Match;
}

export function GoalTimeline({ match }: GoalTimelineProps) {
  const containerRef = useRef<HTMLDivElement>(null);

  const hasResult = match.goals1 !== null && match.goals2 !== null;
  const team1 = getTeamByCode(match.team1);
  const team2 = getTeamByCode(match.team2);
  const totalGoals = (match.goals1 ?? 0) + (match.goals2 ?? 0);

  // Generate fake goal minutes for visualization
  const goalEvents: { minute: number; team: string; color: string }[] = [];
  if (hasResult && totalGoals > 0) {
    const g1 = match.goals1 ?? 0;
    const g2 = match.goals2 ?? 0;

    // Distribute goals across 90 minutes with some randomness
    const allGoals = [
      ...Array.from({ length: g1 }, () => ({
        team: match.team1,
        color: team1?.primaryColor ?? '#00ff9c',
      })),
      ...Array.from({ length: g2 }, () => ({
        team: match.team2,
        color: team2?.primaryColor ?? '#666',
      })),
    ];

    // Assign minutes
    const step = 85 / Math.max(allGoals.length, 1);
    for (let i = 0; i < allGoals.length; i++) {
      goalEvents.push({
        minute: Math.round(5 + i * step + Math.random() * (step * 0.5)),
        ...allGoals[i],
      });
    }
    goalEvents.sort((a, b) => a.minute - b.minute);
  }

  // Animate dots on scroll
  useEffect(() => {
    const el = containerRef.current;
    if (!el || goalEvents.length === 0) return;

    const dots = el.querySelectorAll('.goal-dot');
    gsap.set(dots, { scale: 0, opacity: 0 });

    const trigger = ScrollTrigger.create({
      trigger: el,
      start: 'top 90%',
      once: true,
      onEnter: () => {
        gsap.to(dots, {
          scale: 1,
          opacity: 1,
          duration: 0.4,
          stagger: 0.15,
          ease: 'back.out(2)',
        });
      },
    });

    return () => trigger.kill();
  }, [goalEvents.length]);

  if (!hasResult) return null;

  return (
    <div ref={containerRef} className="mt-2 px-1">
      <div className="relative h-6">
        {/* Timeline bar */}
        <div className="absolute top-1/2 left-0 right-0 h-px bg-white/10 -translate-y-1/2" />

        {/* Minute markers */}
        <div className="absolute top-1/2 left-0 w-1 h-1 rounded-full bg-white/20 -translate-y-1/2" />
        <div className="absolute top-1/2 left-1/2 w-1 h-1 rounded-full bg-white/20 -translate-y-1/2" />
        <div className="absolute top-1/2 right-0 w-1 h-1 rounded-full bg-white/20 -translate-y-1/2" />

        {/* Goal dots */}
        {goalEvents.map((goal, i) => {
          const position = `${(goal.minute / 90) * 100}%`;
          return (
            <div
              key={`goal-${goal.team}-${goal.minute}-${i}`}
              className="goal-dot absolute top-1/2 -translate-y-1/2 -translate-x-1/2 group"
              style={{ left: position }}
            >
              <div
                className="w-3 h-3 rounded-full border-2 border-copa-dark"
                style={{ backgroundColor: goal.color }}
              />
              {/* Tooltip */}
              <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-1 hidden group-hover:block">
                <div className="bg-copa-dark/90 border border-white/10 rounded px-2 py-0.5 whitespace-nowrap">
                  <span className="text-[9px] text-white/70">{goal.minute}'</span>
                </div>
              </div>
            </div>
          );
        })}
      </div>

      {/* Labels */}
      <div className="flex justify-between mt-0.5">
        <span className="text-[8px] text-white/20">0'</span>
        <span className="text-[8px] text-white/20">45'</span>
        <span className="text-[8px] text-white/20">90'</span>
      </div>
    </div>
  );
}
