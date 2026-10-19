import { drag as d3drag } from 'd3-drag';
import {
  type SimulationNodeDatum,
  forceCenter,
  forceLink,
  forceManyBody,
  forceSimulation,
} from 'd3-force';
import { scaleSqrt } from 'd3-scale';
import { select } from 'd3-selection';
import { useEffect, useRef, useState } from 'react';
import 'd3-transition';
import { getTeamByCode } from '../../data/teams';
import { useLocale } from '../../i18n/LocaleContext';
import { type KnockoutTree, getWinner } from '../../lib/knockout';
import type { Match } from '../../types/worldcup';

interface MoleculeGraphProps {
  tree: KnockoutTree;
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

interface MoleculeLabels {
  roundOf16: string;
  quarters: string;
  semi: string;
  final: string;
  champion: string;
}

// Each level holds the teams that reached that round: 32 teams in the Round
// of 32, then the winners of each round, up to the champion.
function buildNodes(
  tree: KnockoutTree,
  labels: MoleculeLabels,
): { nodes: MoleculeNode[]; links: MoleculeLink[] } {
  const nodes: MoleculeNode[] = [];
  const links: MoleculeLink[] = [];

  const r32Teams: (string | null)[] = [];
  for (let i = 0; i < 16; i++) {
    const m = tree.round32[i];
    r32Teams.push(m?.team1 ?? null, m?.team2 ?? null);
  }

  const levels: {
    stage: string;
    label: string;
    size: number;
    teams: (string | null)[];
  }[] = [
    { stage: 'r32', label: '', size: BASE_SIZE, teams: r32Teams },
    {
      stage: 'r16',
      label: labels.roundOf16,
      size: BASE_SIZE * 1.4,
      teams: Array.from({ length: 16 }, (_, i) => winnerOf(tree.round32[i])),
    },
    {
      stage: 'qf',
      label: labels.quarters,
      size: BASE_SIZE * 1.9,
      teams: Array.from({ length: 8 }, (_, i) => winnerOf(tree.round16[i])),
    },
    {
      stage: 'sf',
      label: labels.semi,
      size: BASE_SIZE * 2.5,
      teams: Array.from({ length: 4 }, (_, i) => winnerOf(tree.quarter[i])),
    },
    {
      stage: 'final',
      label: labels.final,
      size: BASE_SIZE * 3.2,
      teams: Array.from({ length: 2 }, (_, i) => winnerOf(tree.semi[i])),
    },
    {
      stage: 'champion',
      label: labels.champion,
      size: BASE_SIZE * 5,
      teams: [tree.champion],
    },
  ];

  let prevIds: string[] = [];
  for (const level of levels) {
    const ids: string[] = [];
    level.teams.forEach((teamCode, i) => {
      const id = `${level.stage}-${i}`;
      nodes.push({
        id,
        teamCode,
        label: level.label,
        size: level.size,
        stage: level.stage,
      });
      ids.push(id);
      if (prevIds.length > 0) {
        links.push({ source: prevIds[i * 2], target: id });
        links.push({ source: prevIds[i * 2 + 1], target: id });
      }
    });
    prevIds = ids;
  }

  return { nodes, links };
}

function winnerOf(match: Match | undefined): string | null {
  return match ? getWinner(match) : null;
}

export function MoleculeGraph({ tree }: MoleculeGraphProps) {
  const { t } = useLocale();
  const svgRef = useRef<SVGSVGElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const [tooltip, setTooltip] = useState<{
    x: number;
    y: number;
    text: string;
  } | null>(null);

  useEffect(() => {
    const svgEl = svgRef.current;
    const container = containerRef.current;
    if (!svgEl || !container) return;

    const width = Math.min(container.clientWidth, 960);
    const height = 550;

    const svg = select(svgEl);
    svg.selectAll('*').remove();
    svg.attr('viewBox', `0 0 ${width} ${height}`);

    const { nodes, links } = buildNodes(tree, t.knockout);

    const radiusScale = scaleSqrt()
      .domain([0, BASE_SIZE * 5])
      .range([0, 32]);

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
      .force('charge', forceManyBody().strength(-120))
      .force('center', forceCenter(width / 2, height / 2));

    // Defs
    const defs = svg.append('defs');
    const filterId = 'mol-shadow';
    const filter = defs
      .append('filter')
      .attr('id', filterId)
      .attr('x', '-30%')
      .attr('y', '-30%')
      .attr('width', '160%')
      .attr('height', '160%');
    filter
      .append('feGaussianBlur')
      .attr('in', 'SourceAlpha')
      .attr('stdDeviation', '2');
    filter.append('feOffset').attr('dx', '1').attr('dy', '1');
    const merge = filter.append('feMerge');
    merge.append('feMergeNode');
    merge.append('feMergeNode').attr('in', 'SourceGraphic');

    // Glow filter for champion
    const glowFilter = defs.append('filter').attr('id', 'glow');
    glowFilter
      .append('feGaussianBlur')
      .attr('stdDeviation', '4')
      .attr('result', 'blur');
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
      .attr('stroke', 'rgba(0, 255, 156, 0.15)')
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
        if (d.stage === 'champion') return 'rgba(0, 255, 156, 0.15)';
        if (d.stage === 'final') return 'rgba(0, 255, 156, 0.08)';
        return 'rgba(255, 255, 255, 0.04)';
      })
      .attr('stroke', (d) => {
        if (d.stage === 'champion') return 'rgba(0, 255, 156, 0.5)';
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
      .attr('fill', (d) =>
        d.stage === 'champion' ? '#00ff9c' : 'rgba(255,255,255,0.25)',
      )
      .attr('font-size', (d) => Math.max(6, radiusScale(d.size) * 0.4))
      .attr('font-weight', '600')
      .text((d) => d.label);

    // Hover
    nodeEls
      .on('mouseenter', function (event, d) {
        select(this)
          .select('circle')
          .transition()
          .duration(150)
          .attr('stroke-width', 3)
          .attr('stroke', '#00ff9c');

        const teamName = d.teamCode
          ? (getTeamByCode(d.teamCode)?.name ?? d.teamCode)
          : null;
        const name =
          teamName && d.label
            ? `${teamName} · ${d.label}`
            : (teamName ?? d.label);
        const rect = svgEl.getBoundingClientRect();
        setTooltip({
          x: event.clientX - rect.left,
          y: event.clientY - rect.top - 30,
          text: name,
        });
      })
      .on('mouseleave', function (_, d) {
        const defaultStroke =
          d.stage === 'champion'
            ? 'rgba(0, 255, 156, 0.5)'
            : d.teamCode
              ? (getTeamByCode(d.teamCode)?.primaryColor ??
                'rgba(255,255,255,0.1)')
              : 'rgba(255, 255, 255, 0.08)';

        select(this)
          .select('circle')
          .transition()
          .duration(300)
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
        .attr('x1', (d: any) => d.source.x ?? 0)
        .attr('y1', (d: any) => d.source.y ?? 0)
        .attr('x2', (d: any) => d.target.x ?? 0)
        .attr('y2', (d: any) => d.target.y ?? 0);

      nodeEls.attr('transform', (d) => `translate(${d.x ?? 0},${d.y ?? 0})`);
    });

    return () => {
      simulation.stop();
    };
  }, [tree, t]);

  return (
    <div ref={containerRef} className="relative w-full overflow-hidden">
      <svg ref={svgRef} className="w-full h-auto" />
      {tooltip && (
        <div
          className="absolute pointer-events-none bg-copa-dark/90 text-white text-xs font-semibold px-3 py-1.5 rounded-md border border-white/10"
          style={{
            left: tooltip.x,
            top: tooltip.y,
            transform: 'translateX(-50%)',
          }}
        >
          {tooltip.text}
        </div>
      )}
    </div>
  );
}
