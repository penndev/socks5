import { ref, computed } from "vue";
import zhCN from "./locales/zh-CN";
import en from "./locales/en";

const locales = {
  "zh-CN": zhCN,
  en: en,
};

var defaultLocale = "zh-CN";

if (navigator && navigator.language) {
  const language = navigator.languages[0] || navigator.language;
  const lower = language.toLowerCase();
  if (lower.startsWith("zh")) defaultLocale = "zh-CN";
  if (lower.startsWith("en")) defaultLocale = "en";
}

// 当前语言
var currentLocale = ref(defaultLocale);

/**
 * 设置语言
 */
function setLocale(locale) {
  currentLocale.value = locale;
  return locale;
}

const locale = computed(() => currentLocale.value);

// 简单 key-value 翻译：缺失 key 时回退原 key
const t = (key) => {
  const dict = locales[locale.value];
  return dict[key] ?? key;
};

export { defaultLocale, t, setLocale };
