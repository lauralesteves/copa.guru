import { useEffect } from 'react';
import type { Translations } from './translations';

export function useMetaTags(t: Translations) {
  useEffect(() => {
    document.documentElement.lang = t.lang;
    document.title = t.siteTitle;

    const tags: Record<string, string> = {
      'meta[name="description"]': t.metaDescription,
      'meta[property="og:title"]': t.siteTitle,
      'meta[property="og:description"]': t.metaDescription,
      'meta[name="twitter:title"]': t.siteTitle,
      'meta[name="twitter:description"]': t.metaDescription,
    };

    for (const [selector, content] of Object.entries(tags)) {
      const el = document.querySelector(selector);
      if (el) el.setAttribute('content', content);
    }
  }, [t]);
}
