import { useEffect, useRef, useState } from 'react';
import {
  type SimulationNodeDatum,
  forceCenter,
  forceLink,
  forceManyBody,
  forceSimulation,
} from 'd3-force';
import { select } from 'd3-selection';
import { drag as d3drag } from 'd3-drag';
import { scaleSqrt } from 'd3-scale';
import 'd3-transition';
import { getTeamByCode } from '../../data/teams';
import type { GroupName, GroupStanding } from '../../types/worldcup';

interface MoleculeGraphProps {
  allGroupStandings: Map<GroupName, GroupStanding[]>;
}

interface MoleculeNode extends SimulationNodeDatum {
  id: string;
  teamCode: string | null;
  label: string;
  size: number;
  stage: string;
}

interface MoleculeLink {
  source: string;
  target: string;
}

const BASE_SIZE = 7;

function buildNodes(
  allGroupStandings: Map<GroupName, GroupStanding[]>,
): { nodes: MoleculeNode[]; links: MoleculeLink[] } {
  const nodes: MoleculeNode[] = [];
  const links: MoleculeLink[] = [];

  // R32 nodes (16 teams: 1st and 2nd from each group)
  const groupPairs: [GroupName, GroupName][] = [
    ['A', 'B'], ['C', 'D'], ['E', 'F'], ['G', 'H'], ['I', 'J'], ['K', 'L'],
  ];

  let r32Index = 0;
  const r32Ids: string[] = [];

  for (const [g1, g2] of groupPairs) {
    const s1 = allGroupStandings.get(g1) ?? [];
    const s2 = allGroupStandings.get(g2) ?? [];

    const first1 = s1[0]?.points > 0 ? s1[0].teamCode : null;
    const second1 = s1[1]?.points > 0 ? s1[1].teamCode : null;
    const first2 = s2[0]?.points > 0 ? s2[0].teamCode : null;
    const second2 = s2[1]?.points > 0 ? s2[1].teamCode : null;

    // Match: 1st of g1 vs 2nd of g2
    const id1 = `r32-${r32Index++}`;
    nodes.push({ id: id1, teamCode: first1, label: `1${g1}`, size: BASE_SIZE, stage: 'r32' });
    r32Ids.push(id1);

    const id2 = `r32-${r32Index++}`;
    nodes.push({ id: id2, teamCode: second2, label: `2${g2}`, size: BASE_SIZE, stage: 'r32' });
    r32Ids.push(id2);

    // Match: 1st of g2 vs 2nd of g1
    const id3 = `r32-${r32Index++}`;
    nodes.push({ id: id3, teamCode: first2, label: `1${g2}`, size: BASE_SIZE, stage: 'r32' });
    r32Ids.push(id3);

    const id4 = `r32-${r32Index++}`;
    nodes.push({ id: id4, teamCode: second1, label: `2${g1}`, size: BASE_SIZE, stage: 'r32' });
    r32Ids.push(id4);
  }

  // R16 nodes (8)
  const r16Ids: string[] = [];
  for (let i = 0; i < 8; i++) {
    const id = `r16-${i}`;
    nodes.push({ id, teamCode: null, label: 'Oitavas', size: BASE_SIZE * 1.6, stage: 'r16' });
    r16Ids.push(id);
    links.push({ source: r32Ids[i * 2], target: id });
    links.push({ source: r32Ids[i * 2 + 1], target: id });
  }

  // QF nodes (4)
  const qfIds: string[] = [];
  for (let i = 0; i < 4; i++) {
    const id = `qf-${i}`;
    nodes.push({ id, teamCode: null, label: 'Quartas', size: BASE_SIZE * 2.2, stage: 'qf' });
    qfIds.push(id);
    links.push({ source: r16Ids[i * 2], target: id });
    links.push({ source: r16Ids[i * 2 + 1], target: id });
  }

  // SF nodes (2)
  const sfIds: string[] = [];
  for (let i = 0; i < 2; i++) {
    const id = `sf-${i}`;
    nodes.push({ id, teamCode: null, label: 'Semi', size: BASE_SIZE * 3, stage: 'sf' });
    sfIds.push(id);
    links.push({ source: qfIds[i * 2], target: id });
    links.push({ source: qfIds[i * 2 + 1], target: id });
  }

  // Final node
  const finalId = 'final';
  nodes.push({ id: finalId, teamCode: null, label: 'Final', size: BASE_SIZE * 3.8, stage: 'final' });
  links.push({ source: sfIds[0], target: finalId });
  links.push({ source: sfIds[1], target: finalId });

  // Champion node
  const champId = 'champion';
  nodes.push({ id: champId, teamCode: null, label: 'Campeão', size: BASE_SIZE * 5, stage: 'champion' });
  links.push({ source: finalId, target: champId });

  return { nodes, links };
}

export function MoleculeGraph({ allGroupStandings }: MoleculeGraphProps) {
  const svgRef = useRef<SVGSVGElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const [tooltip, setTooltip] = useState<{ x: number; y: number; text: string } | null>(null);

  useEffect(() => {
    const svgEl = svgRef.current;
    const container = containerRef.current;
    if (!svgEl || !container) return;

    const width = Math.min(container.clientWidth, 960);
    const height = 550;

    const svg = select(svgEl);
    svg.selectAll('*').remove();
    svg.attr('viewBox', `0 0 ${width} ${height}`);

    const { nodes, links } = buildNodes(allGroupStandings);

    const radiusScale = scaleSqrt().domain([0, BASE_SIZE * 5]).range([0, 36]);

    // Force simulation
    const simulation = forceSimulation<MoleculeNode>(nodes)
      .force(
        'link',
        forceLink<MoleculeNode, MoleculeLink>(links)
          .id((d) => d.id)
          .distance((d) => {
            // biome-ignore lint/suspicious/noExplicitAny: d3 force link typing
            const s = d.source as any;
            const t = d.target as any;
            const sSize = typeof s === 'object' ? s.size : BASE_SIZE;
            const tSize = typeof t === 'object' ? t.size : BASE_SIZE;
            return radiusScale(sSize) + radiusScale(tSize) + 14;
          }),
      )
      .force('charge', forceManyBody().strength(-180))
      .force('center', forceCenter(width / 2, height / 2));

    // Defs
    const defs = svg.append('defs');
    const filterId = 'mol-shadow';
    const filter = defs.append('filter').attr('id', filterId)
      .attr('x', '-30%').attr('y', '-30%').attr('width', '160%').attr('height', '160%');
    filter.append('feGaussianBlur').attr('in', 'SourceAlpha').attr('stdDeviation', '2');
    filter.append('feOffset').attr('dx', '1').attr('dy', '1');
    const merge = filter.append('feMerge');
    merge.append('feMergeNode');
    merge.append('feMergeNode').attr('in', 'SourceGraphic');

    // Glow filter for champion
    const glowFilter = defs.append('filter').attr('id', 'glow');
    glowFilter.append('feGaussianBlur').attr('stdDeviation', '4').attr('result', 'blur');
    const glowMerge = glowFilter.append('feMerge');
    glowMerge.append('feMergeNode').attr('in', 'blur');
    glowMerge.append('feMergeNode').attr('in', 'SourceGraphic');

    // Links
    const linkGroup = svg.append('g');
    const linkEls = linkGroup
      .selectAll('line')
      .data(links)
      .enter()
      .append('line')
      .attr('stroke', 'rgba(212, 168, 67, 0.15)')
      .attr('stroke-width', 1);

    // Nodes
    const nodeGroup = svg.append('g');
    const nodeEls = nodeGroup
      .selectAll<SVGGElement, MoleculeNode>('g')
      .data(nodes)
      .enter()
      .append('g')
      .attr('class', 'mol-node')
      .style('cursor', 'grab');

    // Circle background
    nodeEls
      .append('circle')
      .attr('r', (d) => radiusScale(d.size))
      .attr('fill', (d) => {
        if (d.stage === 'champion') return 'rgba(212, 168, 67, 0.15)';
        if (d.stage === 'final') return 'rgba(212, 168, 67, 0.08)';
        return 'rgba(255, 255, 255, 0.04)';
      })
      .attr('stroke', (d) => {
        if (d.stage === 'champion') return 'rgba(212, 168, 67, 0.5)';
        if (d.teamCode) {
          const team = getTeamByCode(d.teamCode);
          return team?.primaryColor ?? 'rgba(255,255,255,0.1)';
        }
        return 'rgba(255, 255, 255, 0.08)';
      })
      .attr('stroke-width', (d) => (d.stage === 'champion' ? 2 : 1))
      .attr('filter', (d) => (d.stage === 'champion' ? 'url(#glow)' : ''));

    // Flag images (only for nodes with teams)
    nodeEls
      .filter((d) => d.teamCode !== null)
      .append('image')
      .attr('href', (d) => `/images/flags/96/${d.teamCode}.png`)
      .attr('width', (d) => radiusScale(d.size) * 1.5)
      .attr('height', (d) => radiusScale(d.size) * 1.1)
      .attr('x', (d) => -radiusScale(d.size) * 0.75)
      .attr('y', (d) => -radiusScale(d.size) * 0.55)
      .attr('opacity', 0.8)
      .attr('filter', `url(#${filterId})`);

    // Stage label (for empty nodes)
    nodeEls
      .filter((d) => d.teamCode === null)
      .append('text')
      .attr('text-anchor', 'middle')
      .attr('dominant-baseline', 'central')
      .attr('fill', (d) => (d.stage === 'champion' ? '#d4a843' : 'rgba(255,255,255,0.25)'))
      .attr('font-size', (d) => Math.max(6, radiusScale(d.size) * 0.4))
      .attr('font-weight', '600')
      .text((d) => d.label);

    // Hover
    nodeEls
      .on('mouseenter', function (event, d) {
        select(this).select('circle')
          .transition().duration(150)
          .attr('stroke-width', 3)
          .attr('stroke', '#d4a843');

        const name = d.teamCode ? (getTeamByCode(d.teamCode)?.name ?? d.teamCode) : d.label;
        const rect = svgEl.getBoundingClientRect();
        setTooltip({
          x: event.clientX - rect.left,
          y: event.clientY - rect.top - 30,
          text: name,
        });
      })
      .on('mouseleave', function (_, d) {
        const defaultStroke = d.stage === 'champion'
          ? 'rgba(212, 168, 67, 0.5)'
          : d.teamCode
            ? (getTeamByCode(d.teamCode)?.primaryColor ?? 'rgba(255,255,255,0.1)')
            : 'rgba(255, 255, 255, 0.08)';

        select(this).select('circle')
          .transition().duration(300)
          .attr('stroke-width', d.stage === 'champion' ? 2 : 1)
          .attr('stroke', defaultStroke);

        setTooltip(null);
      });

    // Drag
    function dragStarted(event: { active: number }, d: MoleculeNode) {
      if (!event.active) simulation.alphaTarget(0.3).restart();
      d.fx = d.x;
      d.fy = d.y;
    }
    function dragged(event: { x: number; y: number }, d: MoleculeNode) {
      d.fx = event.x;
      d.fy = event.y;
    }
    function dragEnded(event: { active: number }, d: MoleculeNode) {
      if (!event.active) simulation.alphaTarget(0);
      d.fx = null;
      d.fy = null;
    }

    // biome-ignore lint/suspicious/noExplicitAny: d3 drag typing
    const dragBehavior = d3drag<SVGGElement, MoleculeNode>() as any;
    nodeEls.call(
      dragBehavior
        .on('start', dragStarted)
        .on('drag', dragged)
        .on('end', dragEnded),
    );

    // Tick
    simulation.on('tick', () => {
      linkEls
        // biome-ignore lint/suspicious/noExplicitAny: d3 force resolves string ids to objects
        .attr('x1', (d: any) => (d.source.x ?? 0))
        .attr('y1', (d: any) => (d.source.y ?? 0))
        .attr('x2', (d: any) => (d.target.x ?? 0))
        .attr('y2', (d: any) => (d.target.y ?? 0));

      nodeEls.attr('transform', (d) => `translate(${d.x ?? 0},${d.y ?? 0})`);
    });

    return () => {
      simulation.stop();
    };
  }, [allGroupStandings]);

  return (
    <div ref={containerRef} className="relative w-full overflow-hidden">
      <svg ref={svgRef} className="w-full h-auto" />
      {tooltip && (
        <div
          className="absolute pointer-events-none bg-copa-dark/90 text-white text-xs font-semibold px-3 py-1.5 rounded-md border border-white/10"
          style={{ left: tooltip.x, top: tooltip.y, transform: 'translateX(-50%)' }}
        >
          {tooltip.text}
        </div>
      )}
    </div>
  );
}
