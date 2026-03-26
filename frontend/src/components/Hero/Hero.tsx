import { Suspense, lazy, useEffect, useRef } from 'react';
import { useLocale } from '../../i18n/LocaleContext';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';
import { ParticleField } from '../3d/ParticleField';

gsap.registerPlugin(ScrollTrigger);

const Trophy = lazy(() =>
  import('../3d/Trophy').then((m) => ({ default: m.Trophy })),
);

export function Hero() {
  const sectionRef = useRef<HTMLElement>(null);
  const bgRef = useRef<HTMLDivElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);
  const trophyRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const section = sectionRef.current;
    const bg = bgRef.current;
    const content = contentRef.current;
    const trophy = trophyRef.current;

    if (!section || !bg || !content) return;

    const triggers: ScrollTrigger[] = [];

    // Parallax on scroll
    triggers.push(
      ScrollTrigger.create({
        trigger: section,
        start: 'top top',
        end: 'bottom top',
        scrub: true,
        onUpdate: (self) => {
          const p = self.progress;
          gsap.set(bg, { y: p * 150 });
          gsap.set(content, { y: p * -60, opacity: 1 - p * 1.5 });
          if (trophy) gsap.set(trophy, { y: p * -30, opacity: 1 - p * 1.2 });
        },
      }),
    );

    // Entrance: animate text content children (not the trophy)
    const tl = gsap.timeline();
    tl.fromTo(
      content.children,
      { y: 50, opacity: 0 },
      { y: 0, opacity: 1, duration: 0.8, stagger: 0.2, ease: 'power3.out' },
    );

    // Trophy entrance: separate, slightly delayed
    if (trophy) {
      tl.fromTo(
        trophy,
        { y: 30, opacity: 0, scale: 0.9 },
        { y: 0, opacity: 1, scale: 1, duration: 1, ease: 'power3.out' },
        0.1, // start slightly after content
      );
    }

    return () => {
      for (const t of triggers) t.kill();
      tl.kill();
    };
  }, []);

  return (
    <section
      ref={sectionRef}
      className="relative min-h-dvh flex flex-col items-center justify-center text-center px-6 overflow-hidden"
    >
      {/* Background gradient */}
      <div
        ref={bgRef}
        className="absolute inset-0 bg-gradient-to-b from-copa-purple via-copa-dark to-copa-dark will-change-transform"
      />

      {/* Particle field background */}
      <div className="absolute inset-0 z-[1]">
        <ParticleField />
      </div>

      {/* Trophy 3D — separate from content so it doesn't get stagger'd away */}
      <div
        ref={trophyRef}
        className="relative z-10 w-44 h-56 sm:w-56 sm:h-64 md:w-64 md:h-72 mb-2 will-change-transform"
      >
        <Suspense fallback={<div className="w-full h-full" />}>
          <Trophy />
        </Suspense>
      </div>

      {/* Content */}
      <div ref={contentRef} className="relative z-10 will-change-transform">
        <div className="flex items-center justify-center gap-4 md:gap-6">
          <img src="/images/logo.webp" alt="" className="w-16 h-16 sm:w-20 sm:h-20 md:w-28 md:h-28 lg:w-32 lg:h-32" />
          <h1 className="font-display text-6xl sm:text-8xl md:text-[120px] lg:text-[150px] text-copa-gold leading-none tracking-wider">
            COPA.GURU
          </h1>
        </div>

        <HeroText />

        <div className="flex items-center justify-center gap-3 mt-8">
          {['USA', 'CAN', 'MEX'].map((code) => (
            <div
              key={code}
              className="flex items-center gap-2 bg-white/5 border border-white/10 rounded-full px-4 py-2 backdrop-blur-sm"
            >
              <img
                src={`/images/flags/48/${code}.png`}
                alt={code}
                className="w-6 h-4 object-cover rounded-sm"
              />
              <span className="text-sm text-white/60 font-medium">{code}</span>
            </div>
          ))}
        </div>

        <HeroButtons />
      </div>

      {/* Scroll indicator */}
      <div className="absolute bottom-8 left-1/2 -translate-x-1/2 z-10 animate-bounce">
        <svg
          className="w-6 h-6 text-white/30"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          aria-hidden="true"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M19 14l-7 7m0 0l-7-7m7 7V3"
          />
        </svg>
      </div>
    </section>
  );
}

function HeroText() {
  const { t } = useLocale();
  return (
    <p className="font-heading text-lg sm:text-xl md:text-2xl text-white/70 mt-4 tracking-widest uppercase">
      {t.hero.subtitle}
    </p>
  );
}

function HeroButtons() {
  const { t } = useLocale();
  return (
    <div className="mt-12 flex flex-col sm:flex-row gap-4 justify-center">
      <a href="#grupos" className="bg-copa-gold text-copa-dark font-semibold px-8 py-3 rounded-lg hover:bg-copa-gold-light transition-colors text-sm uppercase tracking-wider">
        {t.hero.viewGroups}
      </a>
      <a href="#palpites" className="border border-white/20 text-white/80 font-semibold px-8 py-3 rounded-lg hover:border-copa-gold hover:text-copa-gold transition-colors text-sm uppercase tracking-wider">
        {t.hero.makePredictions}
      </a>
    </div>
  );
}
