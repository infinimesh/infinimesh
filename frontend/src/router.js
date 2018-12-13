import Vue from "vue";
import Router from "vue-router";
import Devices from "./views/Devices.vue";
import RegisterDevice from "./views/RegisterDevice.vue";
import DeleteDevice from "./views/DeleteDevice.vue";
import Shadow from "./views/Shadow.vue";

Vue.use(Router);

export default new Router({
  mode: "history",
  base: process.env.BASE_URL,
  routes: [
    {
      path: "/devices",
      name: "View devices",
      component: Devices
      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      // component: () =>
      //   import(/* webpackChunkName: "Home" */ "./views/Home.vue"),
    },
    {
      path: "/devices/show/:id",
      component: Shadow,
      name: Shadow
    },
    {
      path: "/devices/register",
      name: "Register device",
      component: RegisterDevice
    },
    {
      path: "/devices/:id/delete",
      name: "Delete device",
      component: DeleteDevice
    },
    {
      path: "*",
      redirect: "/devices"
    }
  ]
});
