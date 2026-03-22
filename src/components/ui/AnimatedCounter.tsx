import { useEffect, useRef, useState } from 'react';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';

gsap.registerPlugin(ScrollTrigger);

interface AnimatedCounterProps {
  value: number;
  prefix?: string;
  suffix?: string;
  duration?: number;
  decimals?: number;
  className?: string;
}

export function AnimatedCounter({
  value,
  prefix = '',
  suffix = '',
  duration = 2,
  decimals = 0,
  className = '',
}: AnimatedCounterProps) {
  const ref = useRef<HTMLSpanElement>(null);
  const [display, setDisplay] = useState('0');
  const triggeredRef = useRef(false);

  useEffect(() => {
    const el = ref.current;
    if (!el) return;

    const trigger = ScrollTrigger.create({
      trigger: el,
      start: 'top 90%',
      once: true,
      onEnter: () => {
        if (triggeredRef.current) return;
        triggeredRef.current = true;

        const obj = { val: 0 };
        gsap.to(obj, {
          val: value,
          duration,
          ease: 'power2.out',
          onUpdate: () => {
            setDisplay(
              obj.val.toLocaleString('pt-BR', {
                minimumFractionDigits: decimals,
                maximumFractionDigits: decimals,
              }),
            );
          },
        });
      },
    });

    return () => trigger.kill();
  }, [value, duration, decimals]);

  // Also update if value changes after initial trigger
  useEffect(() => {
    if (triggeredRef.current) {
      setDisplay(
        value.toLocaleString('pt-BR', {
          minimumFractionDigits: decimals,
          maximumFractionDigits: decimals,
        }),
      );
    }
  }, [value, decimals]);

  return (
    <span ref={ref} className={className}>
      {prefix}
      {display}
      {suffix}
    </span>
  );
}
