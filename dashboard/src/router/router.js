import { createRouter, createWebHashHistory } from "vue-router";
import { useAppStore } from "@/store/app";

const routes = [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/Login.vue"),
    meta: {
      requiresAuth: false,
    },
  },
  {
    path: "/",
    name: "Root",
    redirect: { path: "/dashboard/devices" },
    meta: {
      requiresAuth: true,
    },
  },
  {
    path: "/dashboard",
    name: "Dashboard",
    component: () => import("@/views/Dashboard.vue"),
    meta: {
      requiresAuth: true,
    },
    children: [
      {
        path: "devices",
        name: "Devices",
        component: () => import("@/views/dashboard/Devices.vue"),
      },
      {
        path: "accounts",
        name: "Accounts",
        component: () => import("@/views/dashboard/Accounts.vue"),
      },
    ],
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

router.beforeEach(async (to, from) => {
  const store = useAppStore();
  if (to.matched.some((el) => el.meta.requiresAuth) && !store.logged_in) {
    return { name: "Login" };
  }
});

export default router;
