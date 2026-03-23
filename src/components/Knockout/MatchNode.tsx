import { getTeamByCode } from '../../data/teams';

export interface BracketMatch {
  id: string;
  team1: string | null;
  team2: string | null;
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
}

export function MatchNode({
  match,
  x,
  y,
  width = 160,
  height = 56,
  highlighted = false,
}: MatchNodeProps) {
  const team1 = match.team1 ? getTeamByCode(match.team1) : null;
  const team2 = match.team2 ? getTeamByCode(match.team2) : null;

  return (
    <g
      transform={`translate(${x}, ${y})`}
      className="cursor-pointer"
    >
      {/* Background */}
      <rect
        width={width}
        height={height}
        rx={8}
        fill={highlighted ? 'rgba(212, 168, 67, 0.1)' : 'rgba(255, 255, 255, 0.03)'}
        stroke={highlighted ? 'rgba(212, 168, 67, 0.4)' : 'rgba(255, 255, 255, 0.08)'}
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
      <g transform="translate(0, 0)">
        <rect
          width={width}
          height={height / 2}
          rx={8}
          fill="transparent"
        />
        {match.team1 ? (
          <>
            <image
              href={`/images/flags/48/${match.team1}.png`}
              x={10}
              y={6}
              width={20}
              height={14}
            />
            <text
              x={36}
              y={17}
              fill={highlighted ? '#d4a843' : 'rgba(255, 255, 255, 0.8)'}
              fontSize={11}
              fontWeight={600}
            >
              {team1?.name ?? match.team1}
            </text>
          </>
        ) : (
          <text
            x={12}
            y={17}
            fill="rgba(255, 255, 255, 0.2)"
            fontSize={10}
          >
            A definir...
          </text>
        )}
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
        {match.team2 ? (
          <>
            <image
              href={`/images/flags/48/${match.team2}.png`}
              x={10}
              y={6}
              width={20}
              height={14}
            />
            <text
              x={36}
              y={17}
              fill={highlighted ? '#d4a843' : 'rgba(255, 255, 255, 0.8)'}
              fontSize={11}
              fontWeight={600}
            >
              {team2?.name ?? match.team2}
            </text>
          </>
        ) : (
          <text
            x={12}
            y={17}
            fill="rgba(255, 255, 255, 0.2)"
            fontSize={10}
          >
            A definir
          </text>
        )}
      </g>
    </g>
  );
}
