import Vue from "vue";
import Vuetify from "vuetify";
import VueResource from "vue-resource";
import "material-design-icons-iconfont/dist/material-design-icons.css";
import "@mdi/font/css/materialdesignicons.css";
import "./plugins/vuetify";

import App from "./App.vue";
import router from "./router";
import store from "./store";

Vue.config.productionTip = false;

Vue.use(Vuetify, {
  iconfont: "mdi"
});

Vue.use(VueResource);

// Vue.http.options.root = "https://localhost:8081";


new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
