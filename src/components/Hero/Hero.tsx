import { useEffect, useRef } from 'react';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';

gsap.registerPlugin(ScrollTrigger);

const FLOATING_FLAGS = [
  { code: 'BRA', x: '8%', y: '15%', size: 32, delay: 0 },
  { code: 'ARG', x: '85%', y: '20%', size: 28, delay: 0.5 },
  { code: 'GER', x: '12%', y: '70%', size: 24, delay: 1.2 },
  { code: 'FRA', x: '90%', y: '65%', size: 30, delay: 0.8 },
  { code: 'ESP', x: '20%', y: '40%', size: 20, delay: 1.5 },
  { code: 'ENG', x: '78%', y: '80%', size: 22, delay: 0.3 },
  { code: 'POR', x: '5%', y: '50%', size: 26, delay: 1.0 },
  { code: 'JPN', x: '92%', y: '40%', size: 22, delay: 1.8 },
];

export function Hero() {
  const sectionRef = useRef<HTMLElement>(null);
  const bgRef = useRef<HTMLDivElement>(null);
  const circle1Ref = useRef<HTMLDivElement>(null);
  const circle2Ref = useRef<HTMLDivElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);
  const flagsRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const section = sectionRef.current;
    const bg = bgRef.current;
    const c1 = circle1Ref.current;
    const c2 = circle2Ref.current;
    const content = contentRef.current;

    if (!section || !bg || !c1 || !c2 || !content) return;

    const triggers: ScrollTrigger[] = [];

    // Parallax: background moves slower than scroll
    triggers.push(
      ScrollTrigger.create({
        trigger: section,
        start: 'top top',
        end: 'bottom top',
        scrub: true,
        onUpdate: (self) => {
          const p = self.progress;
          gsap.set(bg, { y: p * 150 });
          gsap.set(c1, { y: p * 80, x: p * -30, scale: 1 + p * 0.3 });
          gsap.set(c2, { y: p * 120, x: p * 40, scale: 1 - p * 0.2 });
          gsap.set(content, { y: p * -60, opacity: 1 - p * 1.5 });
        },
      }),
    );

    // Entrance animation
    const tl = gsap.timeline();
    tl.fromTo(
      content.children,
      { y: 50, opacity: 0 },
      { y: 0, opacity: 1, duration: 0.8, stagger: 0.2, ease: 'power3.out' },
    );

    // Floating flags subtle animation
    if (flagsRef.current) {
      const flags = flagsRef.current.children;
      for (let i = 0; i < flags.length; i++) {
        gsap.to(flags[i], {
          y: 'random(-15, 15)',
          x: 'random(-8, 8)',
          rotation: 'random(-5, 5)',
          duration: 'random(3, 5)',
          repeat: -1,
          yoyo: true,
          ease: 'sine.inOut',
          delay: FLOATING_FLAGS[i]?.delay ?? 0,
        });
      }
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
      {/* Background gradient - parallax layer */}
      <div
        ref={bgRef}
        className="absolute inset-0 bg-gradient-to-b from-copa-purple via-copa-dark to-copa-dark will-change-transform"
      />

      {/* Decorative circles - parallax layers */}
      <div
        ref={circle1Ref}
        className="absolute top-1/4 left-1/4 w-96 h-96 bg-copa-gold/5 rounded-full blur-3xl will-change-transform"
      />
      <div
        ref={circle2Ref}
        className="absolute bottom-1/4 right-1/4 w-80 h-80 bg-copa-blue/5 rounded-full blur-3xl will-change-transform"
      />

      {/* Floating flags */}
      <div ref={flagsRef} className="absolute inset-0 pointer-events-none">
        {FLOATING_FLAGS.map((f) => (
          <img
            key={f.code}
            src={`/images/flags/48/${f.code}.png`}
            alt=""
            className="absolute opacity-[0.07] will-change-transform"
            style={{
              left: f.x,
              top: f.y,
              width: f.size,
            }}
            aria-hidden="true"
          />
        ))}
      </div>

      {/* Content - parallax layer */}
      <div ref={contentRef} className="relative z-10 will-change-transform">
        <h1 className="font-display text-6xl sm:text-8xl md:text-[120px] lg:text-[150px] text-copa-gold leading-none tracking-wider">
          COPA.GURU
        </h1>

        <p className="font-heading text-lg sm:text-xl md:text-2xl text-white/70 mt-4 tracking-widest uppercase">
          Copa do Mundo 2026
        </p>

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

        <div className="mt-12 flex flex-col sm:flex-row gap-4 justify-center">
          <a
            href="#grupos"
            className="bg-copa-gold text-copa-dark font-semibold px-8 py-3 rounded-lg hover:bg-copa-gold-light transition-colors text-sm uppercase tracking-wider"
          >
            Ver Grupos
          </a>
          <a
            href="#palpites"
            className="border border-white/20 text-white/80 font-semibold px-8 py-3 rounded-lg hover:border-copa-gold hover:text-copa-gold transition-colors text-sm uppercase tracking-wider"
          >
            Fazer Palpites
          </a>
        </div>
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
