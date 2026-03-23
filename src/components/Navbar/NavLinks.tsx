interface NavLinksProps {
  vertical?: boolean;
  onNavigate?: () => void;
}

const sections = [
  { label: 'Grupos', href: '#grupos' },
  { label: 'Mata-Mata', href: '#mata-mata' },
  { label: 'Mapa Mundial', href: '#mapa-mundial' },
  { label: 'Seleções', href: '#selecoes' },
  { label: 'Comparar', href: '#comparar' },
  { label: 'Palpites', href: '#palpites' },
];

export function NavLinks({ vertical = false, onNavigate }: NavLinksProps) {
  return (
    <div
      className={`flex ${vertical ? 'flex-col gap-4' : 'items-center'}`}
    >
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
