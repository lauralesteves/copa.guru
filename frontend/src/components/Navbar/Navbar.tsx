import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { useLocale } from '../../i18n/LocaleContext';
import { LanguageFlag } from './LanguageFlag';
import { MobileMenu } from './MobileMenu';
import { NavLinks } from './NavLinks';
import { SocialIcons } from './SocialIcons';

interface NavbarProps {
  onLogoTap?: () => void;
}

export function Navbar({ onLogoTap }: NavbarProps) {
  const [scrolled, setScrolled] = useState(false);
  const { basePath } = useLocale();

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
      <div className="flex items-center px-6 pr-[80px]">
        <div className="flex items-center gap-2">
          <img
            src="/images/logo.webp"
            alt=""
            className="w-8 h-8 md:w-10 md:h-10 cursor-pointer"
            onClick={(e) => {
              e.preventDefault();
              onLogoTap?.();
            }}
          />
          <Link
            to={basePath || '/'}
            className="font-display text-copa-gold text-3xl md:text-4xl tracking-wider hover:text-copa-gold-light transition-colors"
          >
            COPA.GURU
          </Link>
        </div>

        <div className="hidden md:flex items-center ml-auto gap-6">
          <NavLinks />
          <div className="border-l border-white/20 h-5" />
          <SocialIcons />
          <div className="border-l border-white/20 h-5" />
          <LanguageFlag />
        </div>

        <div className="ml-auto md:hidden">
          <MobileMenu />
        </div>
      </div>
    </nav>
  );
}
