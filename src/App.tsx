import { Suspense, lazy, useEffect, useState } from 'react';
import { Outlet, Route, Routes, useLocation } from 'react-router-dom';
import { EasterEgg } from './components/EasterEgg';
import { Footer } from './components/Footer/Footer';
import { GithubCorner } from './components/GithubCorner/GithubCorner';
import { Hero } from './components/Hero/Hero';
import { Navbar } from './components/Navbar/Navbar';
import { SectionSkeleton } from './components/ui/SectionSkeleton';
import { useKonamiCode } from './hooks/useKonamiCode';
import { usePredictions } from './hooks/usePredictions';
import { useWorldCupData } from './hooks/useWorldCupData';
import { LocaleContext, getLocaleFromPath } from './i18n/LocaleContext';
import { translations } from './i18n/translations';
import { useMetaTags } from './i18n/useMetaTags';

const GroupStage = lazy(() =>
  import('./components/Groups/GroupStage').then((m) => ({ default: m.GroupStage })),
);
const KnockoutStage = lazy(() =>
  import('./components/Knockout/KnockoutStage').then((m) => ({ default: m.KnockoutStage })),
);
const WorldGlobe = lazy(() =>
  import('./components/Globe/WorldGlobe').then((m) => ({ default: m.WorldGlobe })),
);
const TeamGrid = lazy(() =>
  import('./components/Teams/TeamGrid').then((m) => ({ default: m.TeamGrid })),
);
const StatsSection = lazy(() =>
  import('./components/Stats/StatsSection').then((m) => ({ default: m.StatsSection })),
);
const PredictionSection = lazy(() =>
  import('./components/Predictions/PredictionSection').then((m) => ({ default: m.PredictionSection })),
);
const PrivacyPolicy = lazy(() =>
  import('./components/Legal/PrivacyPolicy').then((m) => ({ default: m.PrivacyPolicy })),
);
const Terms = lazy(() =>
  import('./components/Legal/Terms').then((m) => ({ default: m.Terms })),
);

function ScrollProgress() {
  const [width, setWidth] = useState(0);

  useEffect(() => {
    function handleScroll() {
      const scrollTop = window.scrollY;
      const docHeight = document.documentElement.scrollHeight - window.innerHeight;
      const progress = docHeight > 0 ? (scrollTop / docHeight) * 100 : 0;
      setWidth(progress);
    }
    window.addEventListener('scroll', handleScroll, { passive: true });
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  return <div className="scroll-progress" style={{ width: `${width}%` }} />;
}

function HomePage() {
  const { matches, updateScore, getGroupMatches, allGroupStandings } = useWorldCupData();
  const { predictions, champion, setChampion, stats, clearAll } = usePredictions(matches);

  return (
    <>
      <main>
        <Hero />
        <Suspense fallback={<SectionSkeleton lines={12} />}>
          <GroupStage getGroupMatches={getGroupMatches} allGroupStandings={allGroupStandings} onScoreChange={updateScore} />
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
          <PredictionSection predictions={predictions} champion={champion} setChampion={setChampion} stats={stats} clearAll={clearAll} />
        </Suspense>
      </main>
      <div className="flex justify-center py-12">
        <button type="button" onClick={() => window.scrollTo({ top: 0, behavior: 'smooth' })} className="cursor-pointer">
          <img src="/images/logo.webp" alt="Voltar ao topo" className="w-20 h-20 opacity-40 hover:opacity-70 transition-opacity" />
        </button>
      </div>
    </>
  );
}

function LocaleLayout() {
  const location = useLocation();
  const locale = getLocaleFromPath(location.pathname);
  const t = translations[locale];
  const basePath = locale === 'pt' ? '' : `/${locale}`;
  const { activated: konamiActive, dismiss: dismissKonami, registerTap } = useKonamiCode();

  useMetaTags(t);

  return (
    <LocaleContext.Provider value={{ locale, t, basePath }}>
      <div className="min-h-dvh bg-copa-dark">
        <ScrollProgress />
        <GithubCorner />
        <Navbar onLogoTap={registerTap} />
        <Outlet />
        <Footer />
        <EasterEgg active={konamiActive} onDismiss={dismissKonami} />
      </div>
    </LocaleContext.Provider>
  );
}

function App() {
  return (
    <Routes>
      <Route element={<LocaleLayout />}>
        {/* Portuguese (default) */}
        <Route path="/" element={<HomePage />} />
        <Route path="/privacidade" element={<Suspense fallback={<SectionSkeleton lines={3} />}><PrivacyPolicy /></Suspense>} />
        <Route path="/termos" element={<Suspense fallback={<SectionSkeleton lines={3} />}><Terms /></Suspense>} />
        {/* English */}
        <Route path="/en" element={<HomePage />} />
        <Route path="/en/privacy" element={<Suspense fallback={<SectionSkeleton lines={3} />}><PrivacyPolicy /></Suspense>} />
        <Route path="/en/terms" element={<Suspense fallback={<SectionSkeleton lines={3} />}><Terms /></Suspense>} />
        {/* Spanish */}
        <Route path="/es" element={<HomePage />} />
        <Route path="/es/privacidad" element={<Suspense fallback={<SectionSkeleton lines={3} />}><PrivacyPolicy /></Suspense>} />
        <Route path="/es/terminos" element={<Suspense fallback={<SectionSkeleton lines={3} />}><Terms /></Suspense>} />
        {/* Fallback */}
        <Route path="*" element={<HomePage />} />
      </Route>
    </Routes>
  );
}

export default App;
