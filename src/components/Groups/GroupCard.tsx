import { Suspense, lazy, useState } from 'react';
import { getTeamByCode } from '../../data/teams';
import { useLocale } from '../../i18n/LocaleContext';
import type { Group, GroupStanding, Match } from '../../types/worldcup';
import { Flag } from '../ui/Flag';
import { ChordDiagram } from './ChordDiagram';
import { MatchCard } from './MatchCard';
import { StandingsTable } from './StandingsTable';

const GroupGalaxy = lazy(() =>
  import('./GroupGalaxy').then((m) => ({ default: m.GroupGalaxy })),
);

interface GroupCardProps {
  group: Group;
  matches: Match[];
  standings: GroupStanding[];
  onScoreChange: (
    matchId: number,
    goals1: number | null,
    goals2: number | null,
  ) => void;
}

type Tab = 'table' | 'matches' | 'chord' | 'galaxy';

export function GroupCard({
  group,
  matches,
  standings,
  onScoreChange,
}: GroupCardProps) {
  const { t } = useLocale();
  const [tab, setTab] = useState<Tab>('table');

  return (
    <div className="bg-white/5 border border-white/10 rounded-xl overflow-hidden hover:border-copa-gold/30 transition-colors">
      {/* Header */}
      <div className="bg-gradient-to-r from-copa-purple-light to-copa-purple px-5 py-3 flex items-center justify-between">
        <h3 className="font-display text-2xl text-copa-gold tracking-wider">
          GRUPO {group.name}
        </h3>
        <div className="flex -space-x-1">
          {group.teams.map((code) => {
            if (code.startsWith('TBD')) {
              return (
                <div
                  key={code}
                  className="w-7 h-5 bg-white/10 rounded-sm border border-white/20 flex items-center justify-center"
                >
                  <span className="text-[7px] text-white/60">TBD</span>
                </div>
              );
            }
            return (
              <Flag key={code} code={code} size={48} className="rounded-sm" />
            );
          })}
        </div>
      </div>

      {/* Teams list (always visible) */}
      <div className="px-4 pt-3 pb-1">
        <div className="space-y-1">
          {group.teams.map((code, i) => {
            const team = getTeamByCode(code);
            const isTBD = code.startsWith('TBD');
            const standing = standings.find((s) => s.teamCode === code);
            const position = standings.findIndex((s) => s.teamCode === code);
            const isQualified = position >= 0 && position < 2;

            return (
              <div
                key={code}
                className={`flex items-center gap-2.5 py-1 px-2 rounded-md transition-all duration-300 ${
                  isQualified && standing && standing.points > 0
                    ? 'bg-copa-gold/5 border-l-2 border-copa-gold/40'
                    : 'border-l-2 border-transparent'
                }`}
              >
                <span className="text-white text-xs font-mono w-3">
                  {i + 1}
                </span>
                {isTBD ? (
                  <>
                    <div className="w-7 h-5 bg-white/10 rounded-sm" />
                    <span className="text-white/60 text-xs">
                      {t.groups.tbd}
                    </span>
                  </>
                ) : (
                  <>
                    <Flag code={code} size={48} className="rounded-sm" />
                    <span className="text-white text-xs font-medium flex-1 truncate">
                      {team?.name ?? code}
                    </span>
                    {team?.isHost && (
                      <span className="text-[8px] bg-copa-gold/20 text-copa-gold px-1.5 py-0.5 rounded-full">
                        {t.groups.host}
                      </span>
                    )}
                    {team?.isDebutant && (
                      <span className="text-[8px] bg-green-500/20 text-green-400 px-1.5 py-0.5 rounded-full">
                        {t.groups.debutant}
                      </span>
                    )}
                  </>
                )}
              </div>
            );
          })}
        </div>
      </div>

      {/* Tab selector */}
      <div className="flex border-t border-white/5 mx-4 mt-2">
        {([
          { key: 'table', label: t.groups.tabTable },
          { key: 'matches', label: t.groups.tabMatches },
          { key: 'chord', label: t.groups.tabChart },
          { key: 'galaxy', label: t.groups.tab3D },
        ] as { key: Tab; label: string }[]).map(({ key, label }) => (
          <button
            key={key}
            type="button"
            onClick={() => setTab(key)}
            className={`flex-1 py-2 text-[10px] font-semibold uppercase tracking-wider transition-colors ${
              tab === key
                ? 'text-copa-gold border-b-2 border-copa-gold'
                : 'text-white border-b-2 border-transparent hover:text-white/70'
            }`}
          >
            {label}
          </button>
        ))}
      </div>

      {/* Tab content */}
      <div className="px-4 pb-4">
        {tab === 'table' && <StandingsTable standings={standings} />}

        {tab === 'matches' && (
          <div className="mt-3 space-y-1">
            {matches.length > 0 ? (
              matches.map((match) => (
                <MatchCard
                  key={match.id}
                  match={match}
                  onScoreChange={onScoreChange}
                />
              ))
            ) : (
              <p className="text-white text-xs text-center py-4">
                {t.groups.noMatches}
              </p>
            )}
          </div>
        )}

        {tab === 'chord' && (
          <ChordDiagram teamCodes={group.teams} matches={matches} />
        )}

        {tab === 'galaxy' && (
          <Suspense fallback={<div className="h-[200px] skeleton rounded-lg" />}>
            <GroupGalaxy
              teamCodes={group.teams}
              standings={standings}
              groupName={`Grupo ${group.name}`}
            />
          </Suspense>
        )}
      </div>
    </div>
  );
}
