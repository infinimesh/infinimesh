import Vue from "vue";
import Vuetify from "vuetify";
import VueResource from "vue-resource";
import "material-design-icons-iconfont/dist/material-design-icons.css";
import "@mdi/font/css/materialdesignicons.css";
import "./plugins/vuetify";
import VueDragTree from "vue-drag-tree";
import "vue-drag-tree/dist/vue-drag-tree.min.css";

import App from "./App.vue";
import router from "./router";
import store from "./store";
import VueAnalytics from "vue-analytics";

Vue.config.productionTip = false;

const isProd = process.env.NODE_ENV == "production";

Vue.use(Vuetify, {
  iconfont: "mdi"
});

Vue.use(VueResource);

Vue.use(VueDragTree);

Vue.http.options.root = "$APISERVER_URL";

if (Vue.http.options.root.startsWith("$")) {
  Vue.http.options.root = "http://localhost:8081";
}

console.log(Vue.http.options.root);
if (Vue.http.options.root === "https://api.infinimesh.io") {
  console.log("Enable analytics");
  Vue.use(VueAnalytics, {
    id: "UA-119169777-2",
    router,
    debug: {
      enabled: !isProd,
      sendHitTask: !isProd
      //remove ! for real prod
    }
  });
}

Vue.http.interceptors.push(function(request) {
  if (request.url === "accounts/token") {
    return;
  }
  // modify headers
  request.headers.set("Authorization", `Bearer ${localStorage.token}`);
});

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
