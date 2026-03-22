interface NavLinksProps {
  vertical?: boolean;
  onNavigate?: () => void;
}

const sections = [
  { label: 'Grupos', href: '#grupos' },
  { label: 'Selecoes', href: '#selecoes' },
  { label: 'Mata-Mata', href: '#mata-mata' },
  { label: 'Palpites', href: '#palpites' },
];

export function NavLinks({ vertical = false, onNavigate }: NavLinksProps) {
  return (
    <div
      className={`flex ${vertical ? 'flex-col gap-4' : 'items-center gap-6'}`}
    >
      {sections.map((section) => (
        <a
          key={section.href}
          href={section.href}
          onClick={onNavigate}
          className="text-white/80 hover:text-copa-gold text-sm font-semibold tracking-wide uppercase transition-colors"
        >
          {section.label}
        </a>
      ))}
    </div>
  );
}
