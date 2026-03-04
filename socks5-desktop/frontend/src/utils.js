export function debounce(fn, delay = 1000) {
  let timer = null;

  // 仅保留最后一次触发，降低频繁写入/渲染开销
  return function debounced(...args) {
    clearTimeout(timer);

    timer = setTimeout(() => {
      fn.apply(this, args);
    }, delay);
  };
}