import { useEffect, useRef } from 'react';
import { chord as d3chord, ribbon as d3ribbon } from 'd3-chord';
import { arc as d3arc } from 'd3-shape';
import { select } from 'd3-selection';
import 'd3-transition';
import { getTeamByCode } from '../../data/teams';
import type { Match } from '../../types/worldcup';

interface ChordDiagramProps {
  teamCodes: string[];
  matches: Match[];
}

const MAX_POINTS = 9;
const POINT_RINGS = 9;

export function ChordDiagram({ teamCodes, matches }: ChordDiagramProps) {
  const svgRef = useRef<SVGSVGElement>(null);
  const validTeams = teamCodes.filter((t) => !t.startsWith('TBD'));

  useEffect(() => {
    const svg = select(svgRef.current);
    if (!svgRef.current) return;

    const n = validTeams.length;
    if (n === 0) return;

    // Sizing
    const size = 280;
    const outerRadius = 120;
    const innerRadius = 72;
    const ringPadding = (outerRadius - innerRadius) / POINT_RINGS;
    const flagDistance = outerRadius + 14;

    svg.selectAll('*').remove();
    svg.attr('viewBox', `0 0 ${size} ${size}`);

    const g = svg.append('g').attr('transform', `translate(${size / 2},${size / 2})`);

    // Defs: drop shadow + gradients
    const defs = svg.append('defs');
    const filterId = `shadow-${validTeams.join('')}`;
    const filter = defs.append('filter').attr('id', filterId)
      .attr('x', '-30%').attr('y', '-30%').attr('width', '160%').attr('height', '160%');
    filter.append('feGaussianBlur').attr('in', 'SourceAlpha').attr('stdDeviation', '1.5');
    filter.append('feOffset').attr('dx', '1').attr('dy', '1');
    const merge = filter.append('feMerge');
    merge.append('feMergeNode');
    merge.append('feMergeNode').attr('in', 'SourceGraphic');

    // Compute points & goals per team
    const teamPoints = new Array(n).fill(0);
    const goalsMap = new Map<string, { g1: number; g2: number }>();

    for (const match of matches) {
      const si = validTeams.indexOf(match.team1);
      const ti = validTeams.indexOf(match.team2);
      if (si === -1 || ti === -1) continue;

      if (match.goals1 !== null && match.goals2 !== null) {
        if (match.goals1 > match.goals2) teamPoints[si] += 3;
        else if (match.goals1 < match.goals2) teamPoints[ti] += 3;
        else { teamPoints[si] += 1; teamPoints[ti] += 1; }

        goalsMap.set(`${si}-${ti}`, { g1: match.goals1, g2: match.goals2 });
      }
    }

    // Build matrix for d3-chord (1 = match exists)
    const matrix: number[][] = Array.from({ length: n }, () => new Array(n).fill(0));
    for (const match of matches) {
      const si = validTeams.indexOf(match.team1);
      const ti = validTeams.indexOf(match.team2);
      if (si === -1 || ti === -1) continue;
      matrix[si][ti] = 1;
      matrix[ti][si] = 1;
    }

    const chordLayout = d3chord().padAngle(0.12).sortSubgroups(null).sortChords(null);
    const chords = chordLayout(matrix);

    const sliceAngle = (2 * Math.PI) / n;

    // === POINT RINGS (background + filled) ===
    for (let i = 0; i < n; i++) {
      const team = getTeamByCode(validTeams[i]);
      const startAngle = i * sliceAngle + 0.06;
      const endAngle = (i + 1) * sliceAngle - 0.06;

      // Background rings (9 levels)
      for (let r = 0; r < POINT_RINGS; r++) {
        const rInner = innerRadius + r * ringPadding;
        const rOuter = innerRadius + (r + 1) * ringPadding;

        const arcGen = d3arc<unknown>()
          .innerRadius(rInner)
          .outerRadius(rOuter)
          .startAngle(startAngle)
          .endAngle(endAngle);

        g.append('path')
          .attr('d', arcGen(null as never) ?? '')
          .attr('fill', team?.secondaryColor ?? '#333')
          .attr('fill-opacity', 0.08)
          .attr('stroke', 'rgba(255,255,255,0.03)')
          .attr('stroke-width', 0.3);
      }

      // Filled rings (up to team points)
      const pts = Math.min(teamPoints[i], MAX_POINTS);
      for (let r = 0; r < pts; r++) {
        const rInner = innerRadius + r * ringPadding;
        const rOuter = innerRadius + (r + 1) * ringPadding;

        const arcGen = d3arc<unknown>()
          .innerRadius(rInner)
          .outerRadius(rOuter)
          .startAngle(startAngle)
          .endAngle(endAngle);

        g.append('path')
          .attr('class', `ring ring-${i}`)
          .attr('d', arcGen(null as never) ?? '')
          .attr('fill', team?.primaryColor ?? '#d4a843')
          .attr('fill-opacity', 0)
          .transition()
          .duration(750)
          .delay(r * 60)
          .attr('fill-opacity', 0.5 + (r / POINT_RINGS) * 0.3);
      }

      // Points text in middle of arc
      const midAngle = (startAngle + endAngle) / 2;
      const textR = innerRadius + (outerRadius - innerRadius) / 2;
      g.append('text')
        .attr('x', textR * Math.cos(midAngle - Math.PI / 2))
        .attr('y', textR * Math.sin(midAngle - Math.PI / 2))
        .attr('text-anchor', 'middle')
        .attr('dominant-baseline', 'central')
        .attr('fill', 'white')
        .attr('font-size', '8')
        .attr('font-weight', '700')
        .attr('opacity', pts > 0 ? 1 : 0.25)
        .attr('filter', `url(#${filterId})`)
        .text(teamPoints[i]);

      // Flag image
      const fx = flagDistance * Math.cos(midAngle - Math.PI / 2);
      const fy = flagDistance * Math.sin(midAngle - Math.PI / 2);
      g.append('image')
        .attr('x', fx - 12)
        .attr('y', fy - 8)
        .attr('width', 24)
        .attr('height', 16)
        .attr('href', `/images/flags/48/${validTeams[i]}.png`)
        .attr('filter', `url(#${filterId})`)
        .attr('opacity', 0.85)
        .style('cursor', 'pointer')
        .on('mouseenter', function () {
          // Highlight this team's rings
          g.selectAll(`.ring-${i}`).transition().duration(200).attr('fill-opacity', 0.9);
        })
        .on('mouseleave', function () {
          g.selectAll(`.ring-${i}`).transition().duration(300).attr('fill-opacity', function (_, idx) {
            return 0.5 + (idx / POINT_RINGS) * 0.3;
          });
        });
    }

    // === CHORDS (ribbons) ===
    const ribbonGen = d3ribbon<unknown, unknown>().radius(innerRadius - 2);

    for (const ch of chords) {
      const si = ch.source.index;
      const ti = ch.target.index;
      const sTeam = getTeamByCode(validTeams[si]);
      const tTeam = getTeamByCode(validTeams[ti]);
      const goals = goalsMap.get(`${si}-${ti}`) ?? goalsMap.get(`${ti}-${si}`);
      const hasResult = goals !== undefined;

      // Gradient
      const gradId = `cg-${si}-${ti}`;
      const sAngle = (ch.source.startAngle + ch.source.endAngle) / 2 - Math.PI / 2;
      const tAngle = (ch.target.startAngle + ch.target.endAngle) / 2 - Math.PI / 2;
      const grad = defs.append('linearGradient').attr('id', gradId)
        .attr('gradientUnits', 'userSpaceOnUse')
        .attr('x1', innerRadius * Math.cos(sAngle))
        .attr('y1', innerRadius * Math.sin(sAngle))
        .attr('x2', innerRadius * Math.cos(tAngle))
        .attr('y2', innerRadius * Math.sin(tAngle));
      grad.append('stop').attr('offset', '0%').attr('stop-color', sTeam?.primaryColor ?? '#666');
      grad.append('stop').attr('offset', '100%').attr('stop-color', tTeam?.primaryColor ?? '#666');

      const chordPath = g.append('path')
        .attr('class', 'chord')
        .attr('d', ribbonGen(ch as never) ?? '')
        .attr('fill', `url(#${gradId})`)
        .attr('fill-opacity', hasResult ? 0.25 : 0.07)
        .attr('stroke', 'none')
        .style('cursor', 'pointer')
        .style('transition', 'fill-opacity 0.3s ease');

      // Goal labels (hidden by default)
      // biome-ignore lint/suspicious/noExplicitAny: d3 selection typing
      let goalLabel: any = null;
      if (hasResult && goals) {
        const midA = (sAngle + tAngle) / 2;
        const labelR = innerRadius * 0.35;
        const lx = labelR * Math.cos(midA);
        const ly = labelR * Math.sin(midA);

        // Determine correct order
        const key1 = `${si}-${ti}`;
        const g1 = goalsMap.has(key1) ? goalsMap.get(key1)!.g1 : goalsMap.get(`${ti}-${si}`)!.g2;
        const g2 = goalsMap.has(key1) ? goalsMap.get(key1)!.g2 : goalsMap.get(`${ti}-${si}`)!.g1;

        goalLabel = g.append('text')
          .attr('x', lx)
          .attr('y', ly)
          .attr('text-anchor', 'middle')
          .attr('dominant-baseline', 'central')
          .attr('fill', 'white')
          .attr('font-size', '9')
          .attr('font-weight', '800')
          .attr('filter', `url(#${filterId})`)
          .attr('opacity', 0)
          .style('pointer-events', 'none')
          .text(`${g1} x ${g2}`);
      }

      // Hover interactions
      chordPath
        .on('mouseenter', function () {
          select(this).raise();
          select(this).transition().duration(150).attr('fill-opacity', 0.7);
          goalLabel?.transition().duration(150).attr('opacity', 1);
        })
        .on('mouseleave', function () {
          select(this).transition().duration(300).attr('fill-opacity', hasResult ? 0.25 : 0.07);
          goalLabel?.transition().duration(200).attr('opacity', 0);
        });
    }
  }, [validTeams, matches]);

  return (
    <div className="flex justify-center py-2">
      <svg
        ref={svgRef}
        className="w-full max-w-[280px] h-auto"
      />
    </div>
  );
}
