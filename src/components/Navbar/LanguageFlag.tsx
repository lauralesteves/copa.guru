import { Link, useLocation } from 'react-router-dom';
import {
  buildLocalePath,
  getPagePath,
  useLocale,
} from '../../i18n/LocaleContext';

export function LanguageFlag() {
  const { t } = useLocale();
  const location = useLocation();
  const pagePath = getPagePath(location.pathname);

  return (
    <div className="flex items-center gap-2">
      {t.alternates.map((site) => (
        <Link
          key={site.locale}
          to={buildLocalePath(pagePath, site.locale)}
          className="hover:opacity-80 transition-opacity"
        >
          <img src={site.flag} alt={site.label} className="w-7 h-[18px]" />
        </Link>
      ))}
    </div>
  );
}
