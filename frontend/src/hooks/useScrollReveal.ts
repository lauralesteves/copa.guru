import { useEffect, useRef } from 'react';
import gsap from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';

gsap.registerPlugin(ScrollTrigger);

interface ScrollRevealOptions {
  y?: number;
  x?: number;
  opacity?: number;
  duration?: number;
  delay?: number;
  stagger?: number;
  once?: boolean;
}

export function useScrollReveal<T extends HTMLElement>(
  options: ScrollRevealOptions = {},
) {
  const ref = useRef<T>(null);

  useEffect(() => {
    const el = ref.current;
    if (!el) return;

    const {
      y = 40,
      x = 0,
      opacity = 0,
      duration = 0.8,
      delay = 0,
      once = true,
    } = options;

    gsap.set(el, { y, x, opacity });

    const trigger = ScrollTrigger.create({
      trigger: el,
      start: 'top 85%',
      once,
      onEnter: () => {
        gsap.to(el, {
          y: 0,
          x: 0,
          opacity: 1,
          duration,
          delay,
          ease: 'power3.out',
        });
      },
    });

    return () => trigger.kill();
  }, [options.y, options.x, options.opacity, options.duration, options.delay, options.once]);

  return ref;
}

export function useStaggerReveal<T extends HTMLElement>(
  options: ScrollRevealOptions = {},
) {
  const containerRef = useRef<T>(null);

  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;

    const {
      y = 30,
      opacity = 0,
      duration = 0.6,
      stagger = 0.1,
      once = true,
    } = options;

    const children = el.children;
    if (children.length === 0) return;

    gsap.set(children, { y, opacity });

    const trigger = ScrollTrigger.create({
      trigger: el,
      start: 'top 85%',
      once,
      onEnter: () => {
        gsap.to(children, {
          y: 0,
          opacity: 1,
          duration,
          stagger,
          ease: 'power3.out',
        });
      },
    });

    return () => trigger.kill();
  }, [options.y, options.opacity, options.duration, options.stagger, options.once]);

  return containerRef;
}
