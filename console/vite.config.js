import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import path from "path";
import svgLoader from 'vite-svg-loader';

// https://vitejs.dev/config/
export default defineConfig(({ command }) => {
  let conf = {
    plugins: [vue(), svgLoader()],
    resolve: {
      alias: {
        "@": path.resolve(__dirname, "./src"),
      },
    }
  }

  if (command == 'build') {
    conf.define = {
      INFINIMESH_VERSION_TAG: process.env.INFINIMESH_VERSION_TAG
    }
    if (!process.env.INFINIMESH_VERSION_TAG) conf.define.INFINIMESH_VERSION_TAG = "development"
    
    console.log(`Using version tag: ${conf.define.INFINIMESH_VERSION_TAG}`)
  }

  return conf
});
