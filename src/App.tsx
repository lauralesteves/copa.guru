import { Suspense, lazy, useEffect, useState } from 'react';
import { EasterEgg } from './components/EasterEgg';
import { Footer } from './components/Footer/Footer';
import { Hero } from './components/Hero/Hero';
import { Navbar } from './components/Navbar/Navbar';
import { SectionSkeleton } from './components/ui/SectionSkeleton';
import { useKonamiCode } from './hooks/useKonamiCode';
import { usePredictions } from './hooks/usePredictions';
import { useWorldCupData } from './hooks/useWorldCupData';

const GroupStage = lazy(() =>
  import('./components/Groups/GroupStage').then((m) => ({
    default: m.GroupStage,
  })),
);
const KnockoutStage = lazy(() =>
  import('./components/Knockout/KnockoutStage').then((m) => ({
    default: m.KnockoutStage,
  })),
);
const WorldGlobe = lazy(() =>
  import('./components/Globe/WorldGlobe').then((m) => ({
    default: m.WorldGlobe,
  })),
);
const TeamGrid = lazy(() =>
  import('./components/Teams/TeamGrid').then((m) => ({
    default: m.TeamGrid,
  })),
);
const StatsSection = lazy(() =>
  import('./components/Stats/StatsSection').then((m) => ({
    default: m.StatsSection,
  })),
);
const PredictionSection = lazy(() =>
  import('./components/Predictions/PredictionSection').then((m) => ({
    default: m.PredictionSection,
  })),
);

function ScrollProgress() {
  const [width, setWidth] = useState(0);

  useEffect(() => {
    function handleScroll() {
      const scrollTop = window.scrollY;
      const docHeight =
        document.documentElement.scrollHeight - window.innerHeight;
      const progress = docHeight > 0 ? (scrollTop / docHeight) * 100 : 0;
      setWidth(progress);
    }

    window.addEventListener('scroll', handleScroll, { passive: true });
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  return <div className="scroll-progress" style={{ width: `${width}%` }} />;
}

function App() {
  const { matches, updateScore, getGroupMatches, allGroupStandings } =
    useWorldCupData();
  const { predictions, champion, setChampion, stats, clearAll } =
    usePredictions(matches);
  const { activated: konamiActive, dismiss: dismissKonami } = useKonamiCode();

  return (
    <div className="min-h-dvh bg-copa-dark">
      <ScrollProgress />
      <Navbar />
      <main>
        <Hero />
        <Suspense fallback={<SectionSkeleton lines={12} />}>
          <GroupStage
            getGroupMatches={getGroupMatches}
            allGroupStandings={allGroupStandings}
            onScoreChange={updateScore}
          />
        </Suspense>
        <Suspense fallback={<SectionSkeleton lines={2} />}>
          <KnockoutStage allGroupStandings={allGroupStandings} />
        </Suspense>
        <Suspense fallback={<SectionSkeleton lines={2} />}>
          <WorldGlobe />
        </Suspense>
        <Suspense fallback={<SectionSkeleton lines={8} />}>
          <TeamGrid />
        </Suspense>
        <Suspense fallback={<SectionSkeleton lines={2} />}>
          <StatsSection />
        </Suspense>
        <Suspense fallback={<SectionSkeleton lines={2} />}>
          <PredictionSection
            predictions={predictions}
            champion={champion}
            setChampion={setChampion}
            stats={stats}
            clearAll={clearAll}
          />
        </Suspense>
      </main>
      <div className="flex justify-center py-12">
        <button
          type="button"
          onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
          className="cursor-pointer"
        >
          <img src="/images/logo.webp" alt="Voltar ao topo" className="w-20 h-20 opacity-40 hover:opacity-70 transition-opacity" />
        </button>
      </div>
      <Footer />
      <EasterEgg active={konamiActive} onDismiss={dismissKonami} />
    </div>
  );
}

export default App;
