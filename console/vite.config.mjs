import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import path from "path";
import svgLoader from "vite-svg-loader";

// https://vitejs.dev/config/
export default defineConfig(({ command }) => {
  let conf = {
    plugins: [vue(), svgLoader()],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    },
    define: {
      STATE_MAX_ROWS: "10",
      __INFINIMESH_VERSION_TAG__: !process.env.INFINIMESH_VERSION_TAG ? "'development'" :`'${process.env.INFINIMESH_VERSION_TAG}'`,
    },
    build: {
      chunkSizeWarningLimit: 2000,
    },
  };

  console.log(`Using version tag: ${conf.define.__INFINIMESH_VERSION_TAG__}`);

  return conf;
});
