import { Suspense, lazy, useCallback, useMemo, useState } from 'react';
import { getTeams } from '../../data/teams';
import { teamCoordinates } from '../../data/coordinates';
import { getGroups } from '../../data/groups';
import { getTeamByCode } from '../../data/teams';
import { useScrollReveal } from '../../hooks/useScrollReveal';

const GlobeComponent = lazy(() => import('./GlobeCanvas'));

interface ArcData {
  startLat: number;
  startLng: number;
  endLat: number;
  endLng: number;
  color: string;
}

export function WorldGlobe() {
  const titleRef = useScrollReveal<HTMLHeadingElement>({ y: 30 });
  const [hoveredTeam, setHoveredTeam] = useState<string | null>(null);

  const pointsData = useMemo(() => {
    return getTeams()
      .filter((t) => teamCoordinates[t.code])
      .map((t) => ({
        lat: teamCoordinates[t.code].lat,
        lng: teamCoordinates[t.code].lng,
        code: t.code,
        name: t.name,
        color: t.primaryColor,
        group: t.group,
        size: t.isHost ? 0.6 : 0.35,
        confederation: t.confederation,
      }));
  }, []);

  const arcsData = useMemo(() => {
    const arcs: ArcData[] = [];
    for (const group of getGroups()) {
      const validTeams = group.teams.filter(
        (c) => !c.startsWith('TBD') && teamCoordinates[c],
      );
      for (let i = 0; i < validTeams.length; i++) {
        for (let j = i + 1; j < validTeams.length; j++) {
          const t1 = getTeamByCode(validTeams[i]);
          const c1 = teamCoordinates[validTeams[i]];
          const c2 = teamCoordinates[validTeams[j]];
          if (!c1 || !c2) continue;
          arcs.push({
            startLat: c1.lat,
            startLng: c1.lng,
            endLat: c2.lat,
            endLng: c2.lng,
            color: t1?.primaryColor ?? '#00ff9c',
          });
        }
      }
    }
    return arcs;
  }, []);

  const handlePointHover = useCallback(
    // biome-ignore lint/suspicious/noExplicitAny: globe.gl point typing
    (point: any) => {
      setHoveredTeam(point ? point.code : null);
    },
    [],
  );

  return (
    <section id="mapa-mundial" className="py-20 px-6">
      <div className="max-w-6xl mx-auto">
        <h2
          ref={titleRef}
          className="font-display text-4xl sm:text-5xl md:text-6xl text-copa-gold text-center tracking-wider mb-4"
        >
          MAPA MUNDIAL
        </h2>
        <p className="text-white text-center mb-8 text-sm uppercase tracking-widest">
          Seleções classificadas ao redor do mundo
        </p>

        <div className="relative w-full aspect-square max-w-[600px] mx-auto">
          <Suspense
            fallback={
              <div className="w-full h-full flex items-center justify-center">
                <div className="skeleton w-64 h-64 rounded-full" />
              </div>
            }
          >
            <GlobeComponent
              pointsData={pointsData}
              arcsData={arcsData}
              onPointHover={handlePointHover}
            />
          </Suspense>

          {/* Tooltip */}
          {hoveredTeam && (
            <div className="absolute top-4 left-1/2 -translate-x-1/2 bg-copa-dark/90 border border-white/10 rounded-lg px-4 py-2 flex items-center gap-2 pointer-events-none z-10">
              <img
                src={`/images/flags/48/${hoveredTeam}.png`}
                alt={hoveredTeam}
                className="w-6 h-4 rounded-sm"
              />
              <span className="text-white text-sm font-semibold">
                {getTeamByCode(hoveredTeam)?.name ?? hoveredTeam}
              </span>
              <span className="text-white/60 text-xs">
                Grupo {getTeamByCode(hoveredTeam)?.group}
              </span>
            </div>
          )}
        </div>
      </div>
    </section>
  );
}
