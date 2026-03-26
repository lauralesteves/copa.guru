import { Link } from 'react-router-dom';

const currentYear = new Date().getFullYear();

export function Footer() {
  return (
    <footer className="bg-copa-dark border-t border-white/10 text-white text-center py-6 text-sm font-body">
      <p className="text-white/60 mb-2">
        Se eu fosse você, eu tentaria{' '}
        <span className="hidden md:inline">
          digitar o Konami Code (↑↑↓↓←→←→BA)
        </span>
        <span className="md:hidden">
          clicar 5x rápido no logo do menu
        </span>
      </p>
      <p className="mb-3">
        <Link
          to="/privacidade"
          className="text-white/60 hover:text-copa-gold text-xs transition-colors"
        >
          Política de Privacidade
        </Link>
        <span className="text-white/30 mx-2">·</span>
        <Link
          to="/termos"
          className="text-white/60 hover:text-copa-gold text-xs transition-colors"
        >
          Termos
        </Link>
      </p>
      <p className="font-semibold uppercase tracking-wide text-copa-gold">
        Copa do Mundo 2026 - EUA, Canadá e México
      </p>
      <p className="mt-1 text-white">
        &copy; 2014 - {currentYear} |{' '}
        <a
          href="https://lauraesteves.com"
          target="_blank"
          className="font-semibold text-white/70 hover:text-copa-gold transition-colors"
          rel="noreferrer"
        >
          Laura Esteves
        </a>
      </p>
    </footer>
  );
}
