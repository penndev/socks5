export function debounce(fn, delay = 1000) {
  let timer = null

  return (...args) => {
    clearTimeout(timer)

    timer = setTimeout(() => {
      fn(...args)
    }, delay)
  }
}

/** 为列表项生成稳定 __id，用于选择与延迟映射 */
export function withRuntimeIds(serverList) {
  const list = Array.isArray(serverList) ? serverList : [];
  return list.map((server) => {
    const __id = `${server.host}|${server.protocol}|${server.username}|${server.password}`;
    return { ...server, __id };
  });
}