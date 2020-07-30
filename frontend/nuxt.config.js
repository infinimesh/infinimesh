module.exports = {
  /*
   ** Headers of the page
   */
  head: {
    title: "console.infinimesh.app",
    meta: [
      { charset: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      {
        hid: "description",
        name: "description",
        content: "Console Infinimesh UI"
      }
    ],
    link: [{ rel: "icon", type: "image/x-icon", href: "/favicon.ico" }]
  },
  plugins: [
    "@/plugins/ant-design-vue",
    "@/plugins/typeface-exo"
    // "@/plugins/axios"
  ],
  /*
   ** Customize the progress bar color
   */
  loading: { color: "#3B8070" },
  /*
   ** Build configuration
   */
  build: {
    /*
     ** Run ESLint on save
     */
    extend(config, { isDev, isClient }) {
      if (isDev && isClient) {
        config.module.rules.push({
          enforce: "pre",
          test: /\.(js|vue)$/,
          loader: "eslint-loader",
          exclude: /(node_modules)/
        });
      }
    }
  },
  buildModules: ["@nuxt/typescript-build"],
  modules: ["@nuxtjs/axios", "@nuxtjs/auth"],
  axios: {
    baseURL:
      // process.env.NODE_ENV == "development"
      "https://api.infinimesh.app/"
    // : "http://localhost:8001"
  },
  auth: {
    strategies: {
      local: {
        endpoints: {
          login: {
            url: "account/token",
            method: "post",
            propertyName: "token"
          },
          user: { url: "account", method: "get", propertyName: false }
        },
        tokenType: "bearer"
      }
    },
    redirect: {
      logout: "/login"
    }
  },
  router: {
    middleware: ["auth"]
  },
  build: {
    loaders: {
      less: {
        javascriptEnabled: true,
        modifyVars: require("./assets/styles/antThemeMod.js")
      }
    }
  },
  css: [{ lang: "less", src: "@/assets/styles/themes.less" }],
  server: {
    host: process.env.NODE_ENV == "production" ? "0.0.0.0" : "localhost",
    port: process.env.NODE_ENV == "production" ? 80 : 3000
  }
};
