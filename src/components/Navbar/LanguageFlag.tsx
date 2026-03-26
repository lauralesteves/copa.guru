import { Link, useLocation } from 'react-router-dom';
import {
  buildLocalePath,
  getPagePath,
  useLocale,
} from '../../i18n/LocaleContext';

export function LanguageFlag() {
  const { locale, t } = useLocale();
  const location = useLocation();
  const pagePath = getPagePath(location.pathname);

  return (
    <div className="flex items-center gap-2">
      {t.alternates.map((site) => {
        const isActive = site.locale === locale;
        return (
          <Link
            key={site.locale}
            to={buildLocalePath(pagePath, site.locale)}
            className={`transition-opacity ${isActive ? 'opacity-100' : 'opacity-35 hover:opacity-70'}`}
          >
            <img src={site.flag} alt={site.label} className="w-7 h-[18px]" />
          </Link>
        );
      })}
    </div>
  );
}
