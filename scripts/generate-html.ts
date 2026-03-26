import { mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { translations } from '../src/i18n/translations';
import type { Locale } from '../src/i18n/translations';

const DIST = join(import.meta.dirname, '..', 'dist');
const BASE_URL = 'https://copa.guru';

const template = readFileSync(join(DIST, 'index.html'), 'utf-8');

interface Route {
  path: string;
  locale: Locale;
}

const routes: Route[] = [
  { path: 'index.html', locale: 'pt' },
  { path: 'privacidade/index.html', locale: 'pt' },
  { path: 'termos/index.html', locale: 'pt' },
  { path: 'en/index.html', locale: 'en' },
  { path: 'en/privacy/index.html', locale: 'en' },
  { path: 'en/terms/index.html', locale: 'en' },
  { path: 'es/index.html', locale: 'es' },
  { path: 'es/privacidad/index.html', locale: 'es' },
  { path: 'es/terminos/index.html', locale: 'es' },
];

for (const route of routes) {
  const t = translations[route.locale];
  const slug = route.path.replace(/\/?index\.html$/, '');
  const url = slug ? `${BASE_URL}/${slug}` : BASE_URL;

  const html = template
    .replace(/%LANG%/g, t.lang)
    .replace(/%TITLE%/g, t.siteTitle)
    .replace(/%DESCRIPTION%/g, t.metaDescription)
    .replace(/%URL%/g, url);

  const filePath = join(DIST, route.path);
  mkdirSync(dirname(filePath), { recursive: true });
  writeFileSync(filePath, html);
  console.log(`✓ ${route.path} (${route.locale})`);
}
