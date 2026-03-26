import { Link } from 'react-router-dom';
import { useLocale } from '../../i18n/LocaleContext';

const currentYear = new Date().getFullYear();

export function Footer() {
  const { locale, t } = useLocale();

  return (
    <footer className="bg-copa-dark border-t border-white/10 text-white text-center py-6 text-sm font-body">
      <p className="text-white/60 mb-2">
        {t.footer.hintPrefix}{' '}
        <span className="hidden md:inline">{t.footer.hintDesktop}</span>
        <span className="md:hidden">{t.footer.hintMobile}</span>
      </p>
      {locale === 'pt' && (
        <p className="mb-3">
          <Link
            to="/privacidade"
            className="text-white/60 hover:text-copa-gold text-xs transition-colors"
          >
            {t.footer.privacy}
          </Link>
          <span className="text-white/30 mx-2">·</span>
          <Link
            to="/termos"
            className="text-white/60 hover:text-copa-gold text-xs transition-colors"
          >
            {t.footer.terms}
          </Link>
        </p>
      )}
      <p className="font-semibold uppercase tracking-wide text-copa-gold">
        {t.footer.tagline}
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
