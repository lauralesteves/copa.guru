import { useLocale } from '../../i18n/LocaleContext';

interface NavLinksProps {
  vertical?: boolean;
  onNavigate?: () => void;
}

export function NavLinks({ vertical = false, onNavigate }: NavLinksProps) {
  const { t } = useLocale();

  const sections = [
    { label: t.nav.groups, href: '#grupos' },
    { label: t.nav.knockout, href: '#mata-mata' },
    { label: t.nav.worldMap, href: '#mapa-mundial' },
    { label: t.nav.teams, href: '#selecoes' },
    { label: t.nav.compare, href: '#comparar' },
    { label: t.nav.predictions, href: '#palpites' },
  ];

  return (
    <div className={`flex ${vertical ? 'flex-col gap-4' : 'items-center'}`}>
      {sections.map((section, i) => (
        <span key={section.href} className={vertical ? '' : 'flex items-center'}>
          <a
            href={section.href}
            onClick={onNavigate}
            className="text-white/80 hover:text-copa-gold text-sm font-semibold tracking-wide uppercase transition-colors"
          >
            {section.label}
          </a>
          {!vertical && i < sections.length - 1 && (
            <span className="text-white/20 mx-3">·</span>
          )}
        </span>
      ))}
    </div>
  );
}
