/*
*/

import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

import enTranslation from './locales/en.json';
import frTranslation from './locales/fr.json';
import zhTranslation from './locales/zh.json';
import ruTranslation from './locales/ru.json';
import jaTranslation from './locales/ja.json';
import viTranslation from './locales/vi.json';

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    load: 'languageOnly',
    resources: {
      en: enTranslation,
      zh: zhTranslation,
      fr: frTranslation,
      ru: ruTranslation,
      ja: jaTranslation,
      vi: viTranslation,
    },
    fallbackLng: 'zh',
    interpolation: {
      escapeValue: false,
    },
  });

export default i18n;
