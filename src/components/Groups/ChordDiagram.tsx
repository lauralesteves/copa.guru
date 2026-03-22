import { useEffect, useRef } from 'react';
import { getTeamByCode } from '../../data/teams';
import type { Match } from '../../types/worldcup';

interface ChordDiagramProps {
  teamCodes: string[];
  matches: Match[];
}

interface ChordData {
  source: number;
  target: number;
  sourceGoals: number;
  targetGoals: number;
}

export function ChordDiagram({ teamCodes, matches }: ChordDiagramProps) {
  const svgRef = useRef<SVGSVGElement>(null);
  const validTeams = teamCodes.filter((t) => !t.startsWith('TBD'));

  useEffect(() => {
    const svg = svgRef.current;
    if (!svg) return;

    const size = 200;
    const cx = size / 2;
    const cy = size / 2;
    const outerR = 85;
    const innerR = 65;
    const flagR = outerR + 8;

    // Clear previous render
    while (svg.firstChild) svg.removeChild(svg.firstChild);

    const defs = document.createElementNS('http://www.w3.org/2000/svg', 'defs');

    // Drop shadow filter
    const filter = document.createElementNS('http://www.w3.org/2000/svg', 'filter');
    filter.setAttribute('id', `shadow-${validTeams.join('')}`);
    filter.setAttribute('x', '-20%');
    filter.setAttribute('y', '-20%');
    filter.setAttribute('width', '140%');
    filter.setAttribute('height', '140%');
    const feBlur = document.createElementNS('http://www.w3.org/2000/svg', 'feGaussianBlur');
    feBlur.setAttribute('in', 'SourceAlpha');
    feBlur.setAttribute('stdDeviation', '2');
    const feMerge = document.createElementNS('http://www.w3.org/2000/svg', 'feMerge');
    const feMergeNode1 = document.createElementNS('http://www.w3.org/2000/svg', 'feMergeNode');
    const feMergeNode2 = document.createElementNS('http://www.w3.org/2000/svg', 'feMergeNode');
    feMergeNode2.setAttribute('in', 'SourceGraphic');
    feMerge.appendChild(feMergeNode1);
    feMerge.appendChild(feMergeNode2);
    filter.appendChild(feBlur);
    filter.appendChild(feMerge);
    defs.appendChild(filter);
    svg.appendChild(defs);

    const n = validTeams.length;
    if (n === 0) return;

    const sliceAngle = (2 * Math.PI) / n;
    const padAngle = 0.08;

    // Collect chord data
    const chords: ChordData[] = [];
    for (const match of matches) {
      const si = validTeams.indexOf(match.team1);
      const ti = validTeams.indexOf(match.team2);
      if (si === -1 || ti === -1) continue;

      chords.push({
        source: si,
        target: ti,
        sourceGoals: match.goals1 ?? 0,
        targetGoals: match.goals2 ?? 0,
      });
    }

    // Calculate points per team (for arc sizing)
    const teamPoints = new Array(n).fill(0);
    for (const match of matches) {
      if (match.goals1 === null || match.goals2 === null) continue;
      const si = validTeams.indexOf(match.team1);
      const ti = validTeams.indexOf(match.team2);
      if (si === -1 || ti === -1) continue;

      if (match.goals1 > match.goals2) teamPoints[si] += 3;
      else if (match.goals1 < match.goals2) teamPoints[ti] += 3;
      else {
        teamPoints[si] += 1;
        teamPoints[ti] += 1;
      }
    }

    const maxPoints = Math.max(9, ...teamPoints);

    // Draw outer arcs (points indicator)
    for (let i = 0; i < n; i++) {
      const team = getTeamByCode(validTeams[i]);
      const startAngle = i * sliceAngle + padAngle - Math.PI / 2;
      const fullEnd = (i + 1) * sliceAngle - padAngle - Math.PI / 2;
      const pointsFraction = teamPoints[i] / maxPoints;
      const endAngle = startAngle + (fullEnd - startAngle) * pointsFraction;

      // Background arc (track)
      const bgArc = describeArc(cx, cy, outerR, innerR, startAngle, fullEnd);
      const bgPath = document.createElementNS('http://www.w3.org/2000/svg', 'path');
      bgPath.setAttribute('d', bgArc);
      bgPath.setAttribute('fill', 'rgba(255,255,255,0.05)');
      svg.appendChild(bgPath);

      // Points arc
      if (teamPoints[i] > 0) {
        const pArc = describeArc(cx, cy, outerR, innerR, startAngle, endAngle);
        const pPath = document.createElementNS('http://www.w3.org/2000/svg', 'path');
        pPath.setAttribute('d', pArc);
        pPath.setAttribute('fill', team?.primaryColor ?? '#d4a843');
        pPath.setAttribute('opacity', '0.7');
        pPath.style.transition = 'd 0.75s ease, opacity 0.75s ease';
        svg.appendChild(pPath);
      }

      // Points text
      const midAngle = (startAngle + fullEnd) / 2;
      const textR = (outerR + innerR) / 2;
      const tx = cx + textR * Math.cos(midAngle);
      const ty = cy + textR * Math.sin(midAngle);
      const text = document.createElementNS('http://www.w3.org/2000/svg', 'text');
      text.setAttribute('x', String(tx));
      text.setAttribute('y', String(ty));
      text.setAttribute('text-anchor', 'middle');
      text.setAttribute('dominant-baseline', 'central');
      text.setAttribute('fill', 'white');
      text.setAttribute('font-size', '9');
      text.setAttribute('font-weight', '600');
      text.setAttribute('opacity', teamPoints[i] > 0 ? '1' : '0.3');
      text.textContent = String(teamPoints[i]);
      svg.appendChild(text);

      // Flag image
      const fx = cx + flagR * Math.cos(midAngle) - 10;
      const fy = cy + flagR * Math.sin(midAngle) - 7;
      const img = document.createElementNS('http://www.w3.org/2000/svg', 'image');
      img.setAttribute('x', String(fx));
      img.setAttribute('y', String(fy));
      img.setAttribute('width', '20');
      img.setAttribute('height', '14');
      img.setAttribute('href', `/images/flags/48/${validTeams[i]}.png`);
      img.setAttribute('filter', `url(#shadow-${validTeams.join('')})`);
      svg.appendChild(img);
    }

    // Draw chords (connections between teams)
    for (const chord of chords) {
      const hasResult =
        matches.some(
          (m) =>
            m.goals1 !== null &&
            m.goals2 !== null &&
            ((validTeams.indexOf(m.team1) === chord.source &&
              validTeams.indexOf(m.team2) === chord.target) ||
              (validTeams.indexOf(m.team2) === chord.source &&
                validTeams.indexOf(m.team1) === chord.target)),
        );

      const sTeam = getTeamByCode(validTeams[chord.source]);
      const tTeam = getTeamByCode(validTeams[chord.target]);

      const sAngle =
        (chord.source + 0.5) * sliceAngle - Math.PI / 2;
      const tAngle =
        (chord.target + 0.5) * sliceAngle - Math.PI / 2;

      const sx = cx + innerR * Math.cos(sAngle);
      const sy = cy + innerR * Math.sin(sAngle);
      const tx = cx + innerR * Math.cos(tAngle);
      const ty = cy + innerR * Math.sin(tAngle);

      // Create gradient
      const gradId = `grad-${chord.source}-${chord.target}`;
      const grad = document.createElementNS('http://www.w3.org/2000/svg', 'linearGradient');
      grad.setAttribute('id', gradId);
      grad.setAttribute('x1', String(sx / size));
      grad.setAttribute('y1', String(sy / size));
      grad.setAttribute('x2', String(tx / size));
      grad.setAttribute('y2', String(ty / size));
      const stop1 = document.createElementNS('http://www.w3.org/2000/svg', 'stop');
      stop1.setAttribute('offset', '0%');
      stop1.setAttribute('stop-color', sTeam?.primaryColor ?? '#666');
      const stop2 = document.createElementNS('http://www.w3.org/2000/svg', 'stop');
      stop2.setAttribute('offset', '100%');
      stop2.setAttribute('stop-color', tTeam?.primaryColor ?? '#666');
      grad.appendChild(stop1);
      grad.appendChild(stop2);
      defs.appendChild(grad);

      // Draw quadratic bezier through center
      const path = document.createElementNS('http://www.w3.org/2000/svg', 'path');
      path.setAttribute('d', `M${sx},${sy} Q${cx},${cy} ${tx},${ty}`);
      path.setAttribute('fill', 'none');
      path.setAttribute('stroke', `url(#${gradId})`);
      path.setAttribute('stroke-width', hasResult ? '2.5' : '1');
      path.setAttribute('opacity', hasResult ? '0.6' : '0.15');
      path.style.transition = 'stroke-width 0.5s ease, opacity 0.5s ease';

      // Hover interaction
      path.addEventListener('mouseenter', () => {
        path.setAttribute('opacity', '0.9');
        path.setAttribute('stroke-width', '3.5');
      });
      path.addEventListener('mouseleave', () => {
        path.setAttribute('opacity', hasResult ? '0.6' : '0.15');
        path.setAttribute('stroke-width', hasResult ? '2.5' : '1');
      });

      svg.appendChild(path);

      // Goal labels on chord
      if (hasResult) {
        const midX = (sx + tx) / 2 + (cx - (sx + tx) / 2) * 0.5;
        const midY = (sy + ty) / 2 + (cy - (sy + ty) / 2) * 0.5;
        const goalText = document.createElementNS('http://www.w3.org/2000/svg', 'text');
        goalText.setAttribute('x', String(midX));
        goalText.setAttribute('y', String(midY));
        goalText.setAttribute('text-anchor', 'middle');
        goalText.setAttribute('dominant-baseline', 'central');
        goalText.setAttribute('fill', 'rgba(255,255,255,0.8)');
        goalText.setAttribute('font-size', '8');
        goalText.setAttribute('font-weight', '700');
        goalText.textContent = `${chord.sourceGoals}x${chord.targetGoals}`;
        goalText.style.pointerEvents = 'none';
        svg.appendChild(goalText);
      }
    }
  }, [validTeams, matches]);

  return (
    <div className="flex justify-center py-2">
      <svg
        ref={svgRef}
        viewBox="0 0 200 200"
        className="w-full max-w-[200px] h-auto"
      />
    </div>
  );
}

function describeArc(
  cx: number,
  cy: number,
  outerR: number,
  innerR: number,
  startAngle: number,
  endAngle: number,
): string {
  const outerStart = {
    x: cx + outerR * Math.cos(startAngle),
    y: cy + outerR * Math.sin(startAngle),
  };
  const outerEnd = {
    x: cx + outerR * Math.cos(endAngle),
    y: cy + outerR * Math.sin(endAngle),
  };
  const innerStart = {
    x: cx + innerR * Math.cos(endAngle),
    y: cy + innerR * Math.sin(endAngle),
  };
  const innerEnd = {
    x: cx + innerR * Math.cos(startAngle),
    y: cy + innerR * Math.sin(startAngle),
  };

  const largeArc = endAngle - startAngle > Math.PI ? 1 : 0;

  return [
    `M ${outerStart.x} ${outerStart.y}`,
    `A ${outerR} ${outerR} 0 ${largeArc} 1 ${outerEnd.x} ${outerEnd.y}`,
    `L ${innerStart.x} ${innerStart.y}`,
    `A ${innerR} ${innerR} 0 ${largeArc} 0 ${innerEnd.x} ${innerEnd.y}`,
    'Z',
  ].join(' ');
}
