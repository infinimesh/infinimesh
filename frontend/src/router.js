import Vue from "vue";
import Router from "vue-router";
import Home from "./views/Home.vue";
import RegisterDevice from "./views/RegisterDevice.vue";
import DeleteDevice from "./views/DeleteDevice.vue";

Vue.use(Router);

export default new Router({
  mode: "history",
  base: process.env.BASE_URL,
  routes: [
    {
      path: "/devices",
      name: "View your devices",
      component: Home
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      // component: () =>
      //   import(/* webpackChunkName: "Home" */ "./views/Home.vue")
    },
    {
      path: "/devices/register",
      name: "Register device",
      component: RegisterDevice
    },
    {
      path: "/devices/delete",
      name: "Delete device",
      component: DeleteDevice
    }
  ]
});
