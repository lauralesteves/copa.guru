import { useEffect, useRef } from 'react';
import gsap from 'gsap';
import { Confetti } from './Predictions/Confetti';

interface EasterEggProps {
  active: boolean;
  onDismiss: () => void;
}

export function EasterEgg({ active, onDismiss }: EasterEggProps) {
  const overlayRef = useRef<HTMLDivElement>(null);
  const textRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!active || !overlayRef.current || !textRef.current) return;

    const tl = gsap.timeline({
      onComplete: () => {
        setTimeout(onDismiss, 2000);
      },
    });

    tl.fromTo(
      overlayRef.current,
      { opacity: 0 },
      { opacity: 1, duration: 0.3 },
    )
      .fromTo(
        textRef.current.children,
        { y: 60, opacity: 0, scale: 0.8 },
        {
          y: 0,
          opacity: 1,
          scale: 1,
          duration: 0.6,
          stagger: 0.15,
          ease: 'back.out(1.7)',
        },
        0.2,
      )
      .to(textRef.current, {
        y: -20,
        duration: 0.8,
        ease: 'power2.inOut',
        yoyo: true,
        repeat: 1,
      }, '+=0.5');

    return () => {
      tl.kill();
    };
  }, [active, onDismiss]);

  if (!active) return null;

  return (
    <>
      <div
        ref={overlayRef}
        className="fixed inset-0 z-[300] flex items-center justify-center bg-copa-dark/90 backdrop-blur-md cursor-pointer"
        onClick={onDismiss}
        onKeyDown={(e) => {
          if (e.key === 'Escape') onDismiss();
        }}
        role="button"
        tabIndex={0}
      >
        <div ref={textRef} className="text-center">
          <p className="font-display text-6xl sm:text-8xl text-copa-gold tracking-wider">
            GOOOL!
          </p>
          <p className="font-display text-2xl sm:text-4xl text-white/60 tracking-widest mt-2">
            COPA 2026
          </p>
          <p className="text-white/30 text-sm mt-6">
            Voce descobriu o easter egg!
          </p>
          <p className="text-white/20 text-xs mt-1">
            &#x2191;&#x2191;&#x2193;&#x2193;&#x2190;&#x2192;&#x2190;&#x2192;BA
          </p>
        </div>
      </div>
      <Confetti active={active} duration={4000} />
    </>
  );
}
