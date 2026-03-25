import { useEffect, useRef } from 'react';
import gsap from 'gsap';

interface EasterEggProps {
  active: boolean;
  onDismiss: () => void;
}

export function EasterEgg({ active, onDismiss }: EasterEggProps) {
  const overlayRef = useRef<HTMLDivElement>(null);
  const dribleRef = useRef<HTMLImageElement>(null);
  const golRef = useRef<HTMLImageElement>(null);
  const bicicletaRef = useRef<HTMLImageElement>(null);
  const audio1Ref = useRef<HTMLAudioElement | null>(null);
  const audio2Ref = useRef<HTMLAudioElement | null>(null);
  const tlRef = useRef<gsap.core.Timeline | null>(null);

  useEffect(() => {
    if (!active) return;

    const overlay = overlayRef.current;
    const drible = dribleRef.current;
    const gol = golRef.current;
    const bicicleta = bicicletaRef.current;
    if (!overlay || !drible || !gol || !bicicleta) return;

    // Preload audio
    const audio1 = new Audio('/images/grande-jogada.mp3');
    const audio2 = new Audio('/images/goal-goal.mp3');
    audio1Ref.current = audio1;
    audio2Ref.current = audio2;

    // Reset elements
    gsap.set(drible, { x: '100vw', opacity: 1 });
    gsap.set(gol, { scale: 0, opacity: 0 });
    gsap.set(bicicleta, { x: '100vw', opacity: 1 });

    const tl = gsap.timeline();
    tlRef.current = tl;

    // Phase 1: Background fades in + grande-jogada.mp3 + drible runs left
    tl.fromTo(overlay, { opacity: 0 }, { opacity: 1, duration: 0.3 })
      .call(() => {
        audio1.play().catch(() => {});
      })
      .to(drible, {
        x: '-120vw',
        duration: 2.5,
        ease: 'power1.inOut',
      }, 0.3)

      // Phase 2: After grande-jogada ends, goal-goal.mp3 + gol appears + bicicleta runs
      .call(() => {
        audio2.play().catch(() => {});
      }, [], '+=0.3')
      .to(gol, {
        scale: 1,
        opacity: 1,
        duration: 0.4,
        ease: 'back.out(2)',
      })
      // Shake gol — like jQuery UI: direction up, times 7, distance 10px
      .to(gol, {
        y: '-=10',
        duration: 0.05,
        repeat: 13,
        yoyo: true,
        ease: 'sine.inOut',
      })
      // Bicicleta runs across during shake
      .to(bicicleta, {
        x: '-120vw',
        duration: 2,
        ease: 'power1.inOut',
      }, '<')
      // Fade everything out
      .to([gol, overlay], {
        opacity: 0,
        duration: 0.5,
        delay: 0.3,
        onComplete: () => {
          onDismiss();
        },
      });

    return () => {
      tl.kill();
      audio1.pause();
      audio2.pause();
    };
  }, [active, onDismiss]);

  if (!active) return null;

  return (
    <div
      ref={overlayRef}
      className="fixed inset-0 z-[300] overflow-hidden cursor-pointer"
      onClick={() => {
        tlRef.current?.kill();
        audio1Ref.current?.pause();
        audio2Ref.current?.pause();
        onDismiss();
      }}
      onKeyDown={(e) => {
        if (e.key === 'Escape') {
          tlRef.current?.kill();
          audio1Ref.current?.pause();
          audio2Ref.current?.pause();
          onDismiss();
        }
      }}
      role="button"
      tabIndex={0}
    >
      {/* Background */}
      <div className="absolute inset-0 bg-black/30" />
      <img
        src="/images/allejo-bg.webp"
        alt=""
        className="absolute inset-0 w-full h-full object-cover opacity-70"
        aria-hidden="true"
      />

      {/* Drible - starts right, runs to left */}
      <img
        ref={dribleRef}
        src="/images/allejo-drible.webp"
        alt=""
        className="absolute bottom-4 w-32 h-32 sm:w-40 sm:h-40 object-contain"
        style={{ right: 0 }}
        aria-hidden="true"
      />

      {/* Gol - center, shakes */}
      <img
        ref={golRef}
        src="/images/allejo-gol.webp"
        alt=""
        className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-80 h-80 sm:w-[28rem] sm:h-[28rem] md:w-[34rem] md:h-[34rem] object-contain"
        aria-hidden="true"
      />

      {/* Bicicleta - starts right, runs to left */}
      <img
        ref={bicicletaRef}
        src="/images/allejo-bicicleta.webp"
        alt=""
        className="absolute bottom-4 w-32 h-32 sm:w-40 sm:h-40 object-contain"
        style={{ right: 0 }}
        aria-hidden="true"
      />
    </div>
  );
}
