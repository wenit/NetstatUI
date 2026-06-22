import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import wails from "@wailsio/runtime/plugins/vite";

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
  plugins: [vue(), wails("./bindings")],
  build: {
    // @vueuse/core 14.x emits /* #__PURE__ */ annotations at positions
    // Rolldown's stricter parser rejects; harmless (annotations are just
    // tree-shaking hints that fall back to "include"). Silence the noise.
    rollupOptions: {
      onwarn(warning, defaultHandler) {
        if (warning.code === "INVALID_ANNOTATION" && warning.id?.includes("@vueuse/core")) {
          return
        }
        defaultHandler(warning)
      },
    },
  },
});
