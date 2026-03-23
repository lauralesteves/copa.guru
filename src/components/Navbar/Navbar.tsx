import { useEffect, useState } from 'react';
import { MobileMenu } from './MobileMenu';
import { NavLinks } from './NavLinks';
import { SocialIcons } from './SocialIcons';

export function Navbar() {
  const [scrolled, setScrolled] = useState(false);

  useEffect(() => {
    function handleScroll() {
      setScrolled(window.scrollY > 50);
    }

    window.addEventListener('scroll', handleScroll, { passive: true });
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  return (
    <nav
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
        scrolled
          ? 'bg-copa-dark/90 backdrop-blur-md py-2 shadow-lg shadow-black/20'
          : 'bg-transparent py-4'
      }`}
    >
      <div className="flex items-center px-6">
        <a
          href="#"
          className="flex items-center gap-2 font-display text-copa-gold text-3xl md:text-4xl tracking-wider hover:text-copa-gold-light transition-colors"
        >
          <img src="/images/logo.webp" alt="" className="w-8 h-8 md:w-10 md:h-10" />
          COPA.GURU
        </a>

        <div className="hidden md:flex items-center ml-auto gap-6">
          <NavLinks />

          <div className="border-l border-white/20 h-5" />

          <SocialIcons />
        </div>

        <div className="ml-auto md:hidden">
          <MobileMenu />
        </div>
      </div>
    </nav>
  );
}
