import { useCallback, useEffect, useRef, useState } from 'react';

const KONAMI = [
  'ArrowUp', 'ArrowUp',
  'ArrowDown', 'ArrowDown',
  'ArrowLeft', 'ArrowRight',
  'ArrowLeft', 'ArrowRight',
  'b', 'a',
];

const MOBILE_TAPS = 5;
const MOBILE_TAP_WINDOW = 1500; // 5 taps in 1.5s

export function useKonamiCode() {
  const [activated, setActivated] = useState(false);
  const indexRef = useRef(0);
  const tapTimesRef = useRef<number[]>([]);

  // Desktop: Konami Code
  useEffect(() => {
    function handleKeyDown(e: KeyboardEvent) {
      if (e.key === KONAMI[indexRef.current]) {
        indexRef.current++;
        if (indexRef.current === KONAMI.length) {
          setActivated(true);
          indexRef.current = 0;
        }
      } else {
        indexRef.current = 0;
      }
    }

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  // Mobile: 5 rapid taps on logo (call registerTap from logo onClick)
  const registerTap = useCallback(() => {
    const now = Date.now();
    tapTimesRef.current.push(now);

    // Keep only taps within window
    tapTimesRef.current = tapTimesRef.current.filter(
      (t) => now - t < MOBILE_TAP_WINDOW,
    );

    if (tapTimesRef.current.length >= MOBILE_TAPS) {
      setActivated(true);
      tapTimesRef.current = [];
    }
  }, []);

  const dismiss = () => setActivated(false);

  return { activated, dismiss, registerTap };
}
