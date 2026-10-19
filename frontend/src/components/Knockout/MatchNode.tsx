import { getTeamByCode } from '../../data/teams';

export interface BracketMatch {
  id: string;
  team1: string | null;
  team2: string | null;
  goals1: number | null;
  goals2: number | null;
  penalties1: number | null;
  penalties2: number | null;
  winner: string | null;
  label: string;
  stage: string;
}

interface MatchNodeProps {
  match: BracketMatch;
  x: number;
  y: number;
  width?: number;
  height?: number;
  highlighted?: boolean;
  tbdLabel?: string;
}

interface TeamRowProps {
  code: string | null;
  goals: number | null;
  penalties: number | null;
  isWinner: boolean;
  hasWinner: boolean;
  highlighted: boolean;
  width: number;
  tbdLabel: string;
}

function TeamRow({
  code,
  goals,
  penalties,
  isWinner,
  hasWinner,
  highlighted,
  width,
  tbdLabel,
}: TeamRowProps) {
  if (!code) {
    return (
      <text x={12} y={17} fill="rgba(255, 255, 255, 0.2)" fontSize={10}>
        {tbdLabel}
      </text>
    );
  }

  const team = getTeamByCode(code);
  const dimmed = hasWinner && !isWinner;
  let fill = 'rgba(255, 255, 255, 0.8)';
  if (highlighted) fill = '#00ff9c';
  else if (isWinner) fill = '#ffffff';
  else if (dimmed) fill = 'rgba(255, 255, 255, 0.35)';

  return (
    <>
      <image
        href={`/images/flags/48/${code}.png`}
        x={10}
        y={6}
        width={20}
        height={14}
        opacity={dimmed ? 0.45 : 1}
      />
      <text
        x={36}
        y={17}
        fill={fill}
        fontSize={11}
        fontWeight={isWinner ? 700 : 600}
      >
        {team?.name ?? code}
      </text>
      {goals !== null && (
        <text
          x={width - 10}
          y={17}
          textAnchor="end"
          fill={isWinner ? '#00ff9c' : fill}
          fontSize={12}
          fontWeight={700}
        >
          {penalties !== null && (
            <tspan fontSize={8} fill="rgba(255, 255, 255, 0.4)">
              ({penalties}){' '}
            </tspan>
          )}
          {goals}
        </text>
      )}
    </>
  );
}

export function MatchNode({
  match,
  x,
  y,
  width = 160,
  height = 56,
  highlighted = false,
  tbdLabel = 'A definir',
}: MatchNodeProps) {
  const hasWinner = match.winner !== null;
  const isFinal = match.stage === 'final';

  return (
    <g transform={`translate(${x}, ${y})`} className="cursor-pointer">
      {/* Background */}
      <rect
        width={width}
        height={height}
        rx={8}
        fill={
          highlighted ? 'rgba(0, 255, 156, 0.1)' : 'rgba(255, 255, 255, 0.03)'
        }
        stroke={
          highlighted
            ? 'rgba(0, 255, 156, 0.4)'
            : isFinal
              ? 'rgba(0, 255, 156, 0.3)'
              : 'rgba(255, 255, 255, 0.08)'
        }
        strokeWidth={1}
        className="transition-all duration-300"
      />

      {/* Stage label */}
      <text
        x={width / 2}
        y={-6}
        textAnchor="middle"
        fill="rgba(255, 255, 255, 0.25)"
        fontSize={8}
        fontWeight={600}
        letterSpacing={1}
      >
        {match.label}
      </text>

      {/* Team 1 row */}
      <g>
        <TeamRow
          code={match.team1}
          goals={match.goals1}
          penalties={match.penalties1}
          isWinner={hasWinner && match.winner === match.team1}
          hasWinner={hasWinner}
          highlighted={highlighted}
          width={width}
          tbdLabel={tbdLabel}
        />
      </g>

      {/* Divider */}
      <line
        x1={8}
        y1={height / 2}
        x2={width - 8}
        y2={height / 2}
        stroke="rgba(255, 255, 255, 0.06)"
        strokeWidth={1}
      />

      {/* Team 2 row */}
      <g transform={`translate(0, ${height / 2})`}>
        <TeamRow
          code={match.team2}
          goals={match.goals2}
          penalties={match.penalties2}
          isWinner={hasWinner && match.winner === match.team2}
          hasWinner={hasWinner}
          highlighted={highlighted}
          width={width}
          tbdLabel={tbdLabel}
        />
      </g>
    </g>
  );
}
