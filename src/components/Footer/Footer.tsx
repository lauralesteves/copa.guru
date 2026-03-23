const currentYear = new Date().getFullYear();

export function Footer() {
  return (
    <footer className="bg-copa-dark border-t border-white/10 text-white text-center py-6 text-sm font-body">
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
