import { ref } from "vue";

const THEME_LIGHT = "light";
const THEME_DARK = "dark";
const THEME_SYSTEM = "system";
const DARK_MEDIA_QUERY = "(prefers-color-scheme: dark)";

let mediaQuery = null;
let systemThemeListener = null;

export const resolvedTheme = ref(THEME_LIGHT);

export function detectSystemTheme() {
  if (typeof window === "undefined" || !window.matchMedia) {
    return THEME_LIGHT;
  }
  return window.matchMedia(DARK_MEDIA_QUERY).matches ? THEME_DARK : THEME_LIGHT;
}

function applyThemeClass(theme) {
  if (typeof document === "undefined") return;

  const root = document.documentElement;
  root.classList.remove("theme-light", "theme-dark");
  root.classList.add(theme === THEME_DARK ? "theme-dark" : "theme-light");
}

function bindSystemThemeListener() {
  if (typeof window === "undefined" || !window.matchMedia) {
    return detectSystemTheme();
  }

  const mq = window.matchMedia(DARK_MEDIA_QUERY);
  mediaQuery = mq;

  const nextTheme = mq.matches ? THEME_DARK : THEME_LIGHT;

  // 跟随系统主题变化，实时同步到根节点 class
  systemThemeListener = (event) => {
    const currentTheme = event.matches ? THEME_DARK : THEME_LIGHT;
    resolvedTheme.value = currentTheme;
    applyThemeClass(currentTheme);
  };
  mq.addEventListener("change", systemThemeListener);

  return nextTheme;
}

function clearSystemListener() {
  if (mediaQuery && systemThemeListener) {
    mediaQuery.removeEventListener("change", systemThemeListener);
  }
  mediaQuery = null;
  systemThemeListener = null;
}

export function applyTheme(mode = "system") {
  if (typeof document === "undefined") return THEME_LIGHT;

  let actualTheme = THEME_LIGHT;

  if (mode === THEME_SYSTEM) {
    clearSystemListener();
    actualTheme = bindSystemThemeListener();
  } else {
    clearSystemListener();
    actualTheme = mode === THEME_DARK ? THEME_DARK : THEME_LIGHT;
  }

  resolvedTheme.value = actualTheme;
  applyThemeClass(actualTheme);
  return actualTheme;
}

