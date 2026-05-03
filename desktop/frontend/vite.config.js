import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import wails from "@wailsio/runtime/plugins/vite";
import Components from "unplugin-vue-components/vite";
import { AntDesignVueResolver } from "unplugin-vue-components/resolvers";
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    wails("./bindings"),
    Components({
      resolvers: [AntDesignVueResolver({importStyle: false})],
    }),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@bindings': path.resolve(__dirname, './bindings'),
    }
  },
  // Wails v3 dev proxy dials IPv4 loopback (tcp4) for localhost; Vite's default
  // `localhost` bind can be IPv6-only on macOS, causing persistent proxy 502s.
  server: {
    host: "127.0.0.1",
    strictPort: true,
  },
});
