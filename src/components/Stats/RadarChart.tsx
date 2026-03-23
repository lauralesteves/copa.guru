import { useEffect, useRef } from 'react';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';

gsap.registerPlugin(ScrollTrigger);

interface RadarChartProps {
  labels: string[];
  datasets: {
    label: string;
    values: number[];
    color: string;
  }[];
  size?: number;
}

export function RadarChart({ labels, datasets, size = 220 }: RadarChartProps) {
  const svgRef = useRef<SVGSVGElement>(null);
  const polygonRefs = useRef<(SVGPolygonElement | null)[]>([]);

  const cx = size / 2;
  const cy = size / 2;
  const maxR = size / 2 - 30;
  const n = labels.length;
  const levels = 5;

  function angleFor(i: number) {
    return (Math.PI * 2 * i) / n - Math.PI / 2;
  }

  function pointAt(i: number, value: number): string {
    const angle = angleFor(i);
    const r = (value / 100) * maxR;
    return `${cx + r * Math.cos(angle)},${cy + r * Math.sin(angle)}`;
  }

  // Entrance animation
  useEffect(() => {
    const svg = svgRef.current;
    if (!svg) return;

    const trigger = ScrollTrigger.create({
      trigger: svg,
      start: 'top 85%',
      once: true,
      onEnter: () => {
        for (const poly of polygonRefs.current) {
          if (!poly) continue;
          const target = poly.getAttribute('data-points') ?? '';
          // Start from center
          const centerPoints = Array.from({ length: n }, () => `${cx},${cy}`).join(' ');
          poly.setAttribute('points', centerPoints);
          gsap.to(poly, {
            attr: { points: target },
            duration: 1,
            ease: 'power3.out',
            delay: 0.2,
          });
        }
      },
    });

    return () => trigger.kill();
  }, [cx, cy, n]);

  // Grid lines
  const gridLines = [];
  for (let l = 1; l <= levels; l++) {
    const r = (l / levels) * maxR;
    const points = Array.from({ length: n }, (_, i) => {
      const angle = angleFor(i);
      return `${cx + r * Math.cos(angle)},${cy + r * Math.sin(angle)}`;
    }).join(' ');
    gridLines.push(
      <polygon
        key={`grid-${l}`}
        points={points}
        fill="none"
        stroke="rgba(255,255,255,0.06)"
        strokeWidth={0.5}
      />,
    );
  }

  // Axis lines
  const axisLines = Array.from({ length: n }, (_, i) => {
    const angle = angleFor(i);
    return (
      <line
        key={`axis-${i}`}
        x1={cx}
        y1={cy}
        x2={cx + maxR * Math.cos(angle)}
        y2={cy + maxR * Math.sin(angle)}
        stroke="rgba(255,255,255,0.08)"
        strokeWidth={0.5}
      />
    );
  });

  // Labels
  const labelEls = labels.map((label, i) => {
    const angle = angleFor(i);
    const lr = maxR + 16;
    const x = cx + lr * Math.cos(angle);
    const y = cy + lr * Math.sin(angle);
    return (
      <text
        key={`label-${label}`}
        x={x}
        y={y}
        textAnchor="middle"
        dominantBaseline="central"
        fill="rgba(255,255,255,0.5)"
        fontSize={8}
        fontWeight={600}
      >
        {label}
      </text>
    );
  });

  return (
    <svg ref={svgRef} viewBox={`0 0 ${size} ${size}`} className="w-full max-w-[280px] h-auto mx-auto">
      {gridLines}
      {axisLines}

      {/* Data polygons */}
      {datasets.map((ds, di) => {
        const points = ds.values.map((v, i) => pointAt(i, v)).join(' ');
        return (
          <polygon
            key={ds.label}
            ref={(el) => { polygonRefs.current[di] = el; }}
            points={points}
            data-points={points}
            fill={ds.color}
            fillOpacity={0.15}
            stroke={ds.color}
            strokeWidth={1.5}
            strokeOpacity={0.7}
          />
        );
      })}

      {/* Data dots */}
      {datasets.map((ds) =>
        ds.values.map((v, i) => {
          const angle = angleFor(i);
          const r = (v / 100) * maxR;
          return (
            <circle
              key={`${ds.label}-dot-${labels[i]}`}
              cx={cx + r * Math.cos(angle)}
              cy={cy + r * Math.sin(angle)}
              r={2.5}
              fill={ds.color}
              opacity={0.8}
            />
          );
        }),
      )}

      {labelEls}

      {/* Legend */}
      {datasets.map((ds, i) => (
        <g key={`legend-${ds.label}`} transform={`translate(${size - 70}, ${10 + i * 14})`}>
          <rect width={8} height={8} rx={2} fill={ds.color} opacity={0.7} />
          <text x={12} y={7} fill="rgba(255,255,255,0.6)" fontSize={7} fontWeight={600}>
            {ds.label}
          </text>
        </g>
      ))}
    </svg>
  );
}
