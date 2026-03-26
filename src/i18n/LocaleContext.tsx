import { createContext, useContext } from 'react';
import type { Locale, Translations } from './translations';
import { translations } from './translations';

interface LocaleContextValue {
  locale: Locale;
  t: Translations;
  basePath: string;
}

const LocaleContext = createContext<LocaleContextValue>({
  locale: 'pt',
  t: translations.pt,
  basePath: '',
});

export const useLocale = () => useContext(LocaleContext);
export { LocaleContext };

export function getLocaleFromPath(pathname: string): Locale {
  const segment = pathname.split('/')[1];
  if (segment === 'en' || segment === 'es') return segment;
  return 'pt';
}

export function getPagePath(pathname: string): string {
  const locale = getLocaleFromPath(pathname);
  if (locale === 'pt') return pathname;
  return pathname.replace(`/${locale}`, '') || '/';
}

export function buildLocalePath(
  pagePath: string,
  targetLocale: Locale,
): string {
  if (targetLocale === 'pt') return pagePath || '/';
  return `/${targetLocale}${pagePath === '/' ? '' : pagePath}`;
}
