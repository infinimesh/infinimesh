import { nextTick } from "vue";
import { createRouter, createWebHashHistory } from "vue-router";
import { useAppStore } from "@/store/app";

const BASE_TITLE = 'infinimesh';
const routes = [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/Login.vue"),
    meta: {
      title: "Login",
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
      title: "Dashboard",
      requiresAuth: true,
    },
    children: [
      {
        path: "devices",
        name: "Devices",
        component: () => import("@/views/dashboard/Devices.vue"),
        meta: {
          title: "Devices",
        }
      },
      {
        path: "accounts",
        name: "Accounts",
        component: () => import("@/views/dashboard/Accounts.vue"),
        meta: {
          title: "Accounts",
        }
      },
      {
        path: "namespaces",
        name: "Namespaces",
        component: () => import("@/views/dashboard/Namespaces.vue"),
        meta: {
          title: "Namespaces",
        },
      }
    ],
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

router.afterEach((to, from) => {
  nextTick(() => {
    document.title = [BASE_TITLE, to.meta.title].join(" | ");
  });
});

router.beforeEach(async (to, from) => {
  const store = useAppStore();
  if (to.matched.some((el) => el.meta.requiresAuth) && !store.logged_in) {
    return { name: "Login" };
  }
});

export default router;
