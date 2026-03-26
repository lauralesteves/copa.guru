import { useEffect, useRef, useState } from 'react';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';
import type { GroupName, GroupStanding } from '../../types/worldcup';
import { type BracketMatch, MatchNode } from './MatchNode';

gsap.registerPlugin(ScrollTrigger);

interface BracketProps {
  allGroupStandings: Map<GroupName, GroupStanding[]>;
}

// 2026 World Cup bracket:
// 48 teams -> 32 advance from groups (top 2 per group + 8 best 3rd)
// Round of 32 -> Round of 16 -> QF -> SF -> Final
// For simplicity, we show the top-2 per group advancing to R32

// R32 matchups (FIFA confirmed structure):
// The bracket pairs groups: A vs B, C vs D, E vs F, G vs H, I vs J, K vs L
function buildBracketMatches(
  allGroupStandings: Map<GroupName, GroupStanding[]>,
): BracketMatch[] {
  function getTeam(group: GroupName, pos: number): string | null {
    const standings = allGroupStandings.get(group);
    if (!standings || !standings[pos]) return null;
    const s = standings[pos];
    return s.points > 0 ? s.teamCode : null;
  }

  // Round of 32 (16 matches) - group pairings
  const r32: BracketMatch[] = [
    { id: 'r32-1', team1: getTeam('A', 0), team2: getTeam('B', 1), label: 'R32', stage: 'round32' },
    { id: 'r32-2', team1: getTeam('B', 0), team2: getTeam('A', 1), label: 'R32', stage: 'round32' },
    { id: 'r32-3', team1: getTeam('C', 0), team2: getTeam('D', 1), label: 'R32', stage: 'round32' },
    { id: 'r32-4', team1: getTeam('D', 0), team2: getTeam('C', 1), label: 'R32', stage: 'round32' },
    { id: 'r32-5', team1: getTeam('E', 0), team2: getTeam('F', 1), label: 'R32', stage: 'round32' },
    { id: 'r32-6', team1: getTeam('F', 0), team2: getTeam('E', 1), label: 'R32', stage: 'round32' },
    { id: 'r32-7', team1: getTeam('G', 0), team2: getTeam('H', 1), label: 'R32', stage: 'round32' },
    { id: 'r32-8', team1: getTeam('H', 0), team2: getTeam('G', 1), label: 'R32', stage: 'round32' },
    { id: 'r32-9', team1: getTeam('I', 0), team2: getTeam('J', 1), label: 'R32', stage: 'round32' },
    { id: 'r32-10', team1: getTeam('J', 0), team2: getTeam('I', 1), label: 'R32', stage: 'round32' },
    { id: 'r32-11', team1: getTeam('K', 0), team2: getTeam('L', 1), label: 'R32', stage: 'round32' },
    { id: 'r32-12', team1: getTeam('L', 0), team2: getTeam('K', 1), label: 'R32', stage: 'round32' },
    // 3rd place qualifiers (4 slots)
    { id: 'r32-13', team1: null, team2: null, label: 'R32', stage: 'round32' },
    { id: 'r32-14', team1: null, team2: null, label: 'R32', stage: 'round32' },
    { id: 'r32-15', team1: null, team2: null, label: 'R32', stage: 'round32' },
    { id: 'r32-16', team1: null, team2: null, label: 'R32', stage: 'round32' },
  ];

  // R16 placeholders
  const r16: BracketMatch[] = Array.from({ length: 8 }, (_, i) => ({
    id: `r16-${i + 1}`,
    team1: null,
    team2: null,
    label: 'OITAVAS',
    stage: 'round16',
  }));

  // QF placeholders
  const qf: BracketMatch[] = Array.from({ length: 4 }, (_, i) => ({
    id: `qf-${i + 1}`,
    team1: null,
    team2: null,
    label: 'QUARTAS',
    stage: 'quarter',
  }));

  // SF placeholders
  const sf: BracketMatch[] = Array.from({ length: 2 }, (_, i) => ({
    id: `sf-${i + 1}`,
    team1: null,
    team2: null,
    label: 'SEMI',
    stage: 'semi',
  }));

  // Final
  const final: BracketMatch[] = [
    { id: 'final', team1: null, team2: null, label: 'FINAL', stage: 'final' },
  ];

  return [...r32, ...r16, ...qf, ...sf, ...final];
}

// Layout config
const NODE_W = 150;
const NODE_H = 52;
const H_GAP = 50;
const V_GAP = 12;

interface RoundLayout {
  x: number;
  startY: number;
  spacing: number;
  count: number;
}

function computeLayout(): RoundLayout[] {
  // R32: 16 matches, left 8 + right 8
  // We'll show as a single left-to-right bracket for the top half (8 R32)
  // and bottom half handled by scrolling
  // Simplified: show 8 R32 on left side, flowing into R16 -> QF -> SF -> Final
  const rounds: RoundLayout[] = [];
  const baseX = 20;

  // Round of 32 (8 visible in top bracket)
  rounds.push({
    x: baseX,
    startY: 30,
    spacing: NODE_H + V_GAP,
    count: 8,
  });

  // Round of 16
  rounds.push({
    x: baseX + NODE_W + H_GAP,
    startY: 30 + (NODE_H + V_GAP) / 2,
    spacing: (NODE_H + V_GAP) * 2,
    count: 4,
  });

  // Quarter-finals
  rounds.push({
    x: baseX + (NODE_W + H_GAP) * 2,
    startY: 30 + (NODE_H + V_GAP) * 1.5,
    spacing: (NODE_H + V_GAP) * 4,
    count: 2,
  });

  // Semi-final
  rounds.push({
    x: baseX + (NODE_W + H_GAP) * 3,
    startY: 30 + (NODE_H + V_GAP) * 3.5,
    spacing: 0,
    count: 1,
  });

  // Final
  rounds.push({
    x: baseX + (NODE_W + H_GAP) * 4,
    startY: 30 + (NODE_H + V_GAP) * 3.5,
    spacing: 0,
    count: 1,
  });

  return rounds;
}

interface NodePos {
  match: BracketMatch;
  x: number;
  y: number;
}

interface ConnectorLine {
  x1: number;
  y1: number;
  x2: number;
  y2: number;
}

export function Bracket({ allGroupStandings }: BracketProps) {
  const svgRef = useRef<SVGSVGElement>(null);
  const linesRef = useRef<SVGGElement>(null);
  const [highlightedTeam, setHighlightedTeam] = useState<string | null>(null);

  const bracketMatches = buildBracketMatches(allGroupStandings);
  const rounds = computeLayout();

  // Build positioned nodes for top-half bracket
  const topR32 = bracketMatches.filter((m) => m.stage === 'round32').slice(0, 8);
  const r16Matches = bracketMatches.filter((m) => m.stage === 'round16').slice(0, 4);
  const qfMatches = bracketMatches.filter((m) => m.stage === 'quarter').slice(0, 2);
  const sfMatches = bracketMatches.filter((m) => m.stage === 'semi').slice(0, 1);
  const finalMatch = bracketMatches.filter((m) => m.stage === 'final');

  // Bottom-half R32
  const bottomR32 = bracketMatches.filter((m) => m.stage === 'round32').slice(8, 16);
  const r16Bottom = bracketMatches.filter((m) => m.stage === 'round16').slice(4, 8);
  const qfBottom = bracketMatches.filter((m) => m.stage === 'quarter').slice(2, 4);
  const sfBottom = bracketMatches.filter((m) => m.stage === 'semi').slice(1, 2);

  const allRounds = [
    { matches: topR32, round: rounds[0] },
    { matches: r16Matches, round: rounds[1] },
    { matches: qfMatches, round: rounds[2] },
    { matches: sfMatches, round: rounds[3] },
    { matches: finalMatch, round: rounds[4] },
  ];

  // Bottom half offset
  const bottomOffset = (NODE_H + V_GAP) * 8 + 60;
  const allRoundsBottom = [
    { matches: bottomR32, round: { ...rounds[0], startY: rounds[0].startY + bottomOffset } },
    { matches: r16Bottom, round: { ...rounds[1], startY: rounds[1].startY + bottomOffset } },
    { matches: qfBottom, round: { ...rounds[2], startY: rounds[2].startY + bottomOffset } },
    { matches: sfBottom, round: { ...rounds[3], startY: rounds[3].startY + bottomOffset } },
  ];

  const nodePositions: NodePos[] = [];
  const connectors: ConnectorLine[] = [];

  function layoutRound(
    matches: BracketMatch[],
    round: RoundLayout,
  ) {
    matches.forEach((match, i) => {
      const y = round.startY + i * round.spacing;
      nodePositions.push({ match, x: round.x, y });
    });
  }

  function connectRounds(
    parentPositions: NodePos[],
    childPositions: NodePos[],
  ) {
    for (let i = 0; i < childPositions.length; i++) {
      const p1 = parentPositions[i * 2];
      const p2 = parentPositions[i * 2 + 1];
      const child = childPositions[i];
      if (!p1 || !p2 || !child) continue;

      // Line from parent1 right-center to child left-center
      connectors.push({
        x1: p1.x + NODE_W,
        y1: p1.y + NODE_H / 2,
        x2: child.x,
        y2: child.y + NODE_H / 4,
      });
      connectors.push({
        x1: p2.x + NODE_W,
        y1: p2.y + NODE_H / 2,
        x2: child.x,
        y2: child.y + (NODE_H * 3) / 4,
      });
    }
  }

  // Layout top half
  for (const { matches, round } of allRounds) {
    layoutRound(matches, round);
  }
  // Layout bottom half
  for (const { matches, round } of allRoundsBottom) {
    layoutRound(matches, round);
  }

  // Connect top half rounds
  const topR32Pos = nodePositions.filter((n) => topR32.includes(n.match));
  const r16Pos = nodePositions.filter((n) => r16Matches.includes(n.match));
  const qfPos = nodePositions.filter((n) => qfMatches.includes(n.match));
  const sfPos = nodePositions.filter((n) => sfMatches.includes(n.match));
  const finalPos = nodePositions.filter((n) => finalMatch.includes(n.match));

  connectRounds(topR32Pos, r16Pos);
  connectRounds(r16Pos, qfPos);
  connectRounds(qfPos, sfPos);

  // Connect bottom half
  const bottomR32Pos = nodePositions.filter((n) => bottomR32.includes(n.match));
  const r16BottomPos = nodePositions.filter((n) => r16Bottom.includes(n.match));
  const qfBottomPos = nodePositions.filter((n) => qfBottom.includes(n.match));
  const sfBottomPos = nodePositions.filter((n) => sfBottom.includes(n.match));

  connectRounds(bottomR32Pos, r16BottomPos);
  connectRounds(r16BottomPos, qfBottomPos);
  connectRounds(qfBottomPos, sfBottomPos);

  // Connect semis to final
  if (sfPos[0] && sfBottomPos[0] && finalPos[0]) {
    connectors.push({
      x1: sfPos[0].x + NODE_W,
      y1: sfPos[0].y + NODE_H / 2,
      x2: finalPos[0].x,
      y2: finalPos[0].y + NODE_H / 4,
    });
    connectors.push({
      x1: sfBottomPos[0].x + NODE_W,
      y1: sfBottomPos[0].y + NODE_H / 2,
      x2: finalPos[0].x,
      y2: finalPos[0].y + (NODE_H * 3) / 4,
    });
  }

  const svgW = 20 + (NODE_W + H_GAP) * 5 + 40;
  const svgH = bottomOffset + (NODE_H + V_GAP) * 8 + 80;

  // Animate connector lines on scroll
  useEffect(() => {
    if (!linesRef.current) return;
    const paths = linesRef.current.querySelectorAll('path');

    for (const path of paths) {
      const length = path.getTotalLength();
      gsap.set(path, {
        strokeDasharray: length,
        strokeDashoffset: length,
      });
    }

    const trigger = ScrollTrigger.create({
      trigger: svgRef.current,
      start: 'top 80%',
      once: true,
      onEnter: () => {
        gsap.to(linesRef.current?.querySelectorAll('path') ?? [], {
          strokeDashoffset: 0,
          duration: 1.2,
          stagger: 0.04,
          ease: 'power2.inOut',
        });
      },
    });

    return () => trigger.kill();
  }, []);

  function isTeamInMatch(match: BracketMatch, team: string): boolean {
    return match.team1 === team || match.team2 === team;
  }

  return (
    <div className="overflow-x-auto pb-4">
      <svg
        ref={svgRef}
        viewBox={`0 0 ${svgW} ${svgH}`}
        className="w-full min-w-[900px] h-auto"
        onMouseLeave={() => setHighlightedTeam(null)}
      >
        {/* Connector lines */}
        <g ref={linesRef}>
          {connectors.map((c, i) => {
            const midX = (c.x1 + c.x2) / 2;
            return (
              <path
                key={`conn-${i}-${c.x1}-${c.y1}`}
                d={`M${c.x1},${c.y1} C${midX},${c.y1} ${midX},${c.y2} ${c.x2},${c.y2}`}
                fill="none"
                stroke="rgba(0, 255, 156, 0.15)"
                strokeWidth={1.5}
              />
            );
          })}
        </g>

        {/* Match nodes */}
        {nodePositions.map((np) => (
          <g
            key={np.match.id}
            onMouseEnter={() => {
              if (np.match.team1) setHighlightedTeam(np.match.team1);
            }}
          >
            <MatchNode
              match={np.match}
              x={np.x}
              y={np.y}
              width={NODE_W}
              height={NODE_H}
              highlighted={
                highlightedTeam !== null &&
                isTeamInMatch(np.match, highlightedTeam)
              }
            />
          </g>
        ))}

        {/* "CHAMPION" label above final */}
        {finalPos[0] && (
          <text
            x={finalPos[0].x + NODE_W / 2}
            y={finalPos[0].y - 20}
            textAnchor="middle"
            fill="rgba(0, 255, 156, 0.6)"
            fontSize={14}
            fontWeight={700}
            letterSpacing={4}
            fontFamily="'Bebas Neue', sans-serif"
          >
            CAMPEAO
          </text>
        )}
      </svg>
    </div>
  );
}
