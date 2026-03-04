import { defineStore } from "pinia";

export const useServerStore = defineStore("server", {
  state: () => ({
    // 当前选中的远端节点（由列表组件维护）
    selectedServer: null,
  }),
});
