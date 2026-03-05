import { ref, computed } from "vue";
import zhCN from "./locales/zh-CN";
import en from "./locales/en";

const messages = {
  "zh-CN": zhCN,
  en,
};

const DEFAULT_LOCALE = "en";

/**
 * 语言归一化
 */
function normalizeLocale(lang) {
  if (!lang) return DEFAULT_LOCALE;

  const lower = lang.toLowerCase();

  if (lower.startsWith("zh")) return "zh-CN";
  if (lower.startsWith("en")) return "en";

  return DEFAULT_LOCALE;
}

function hasLocale(locale) {
  return Object.prototype.hasOwnProperty.call(messages, locale);
}

function getNavigatorLanguage() {
  if (typeof navigator === "undefined") {
    return "";
  }

  return (
    navigator.language ||
    (Array.isArray(navigator.languages) && navigator.languages[0]) ||
    ""
  );
}

/**
 * 获取系统语言
 */
export function detectSystemLocale() {
  const navigatorLanguage = getNavigatorLanguage();
  return normalizeLocale(navigatorLanguage);
}

// 当前语言
const currentLocale = ref(detectSystemLocale());

/**
 * 设置语言
 */
export function setLocale(locale) {
  const finalLocale = hasLocale(locale)
    ? locale
    : normalizeLocale(locale || detectSystemLocale());

  currentLocale.value = finalLocale;
  return finalLocale;
}

/**
 * i18n hook
 */
export function useI18n() {
  const locale = computed(() => currentLocale.value);

  // 简单 key-value 翻译：缺失 key 时回退原 key
  const t = (key) => {
    const dict = messages[locale.value] || messages[DEFAULT_LOCALE];
    return dict[key] ?? key;
  };

  return {
    t,
    locale,
  };
};
