import { nextTick } from "vue";
import { createRouter, createWebHashHistory } from "vue-router";
import { useAppStore } from "@/store/app";

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
    component: () => import("@/views/index.vue"),
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
        path: "",
        name: "DashboardMain",
        component: () => import("@/views/dashboard/index.vue"),
        meta: {
          title: "Dashboard",
        },
      },
      {
        path: "devices",
        name: "Devices",
        component: () => import("@/views/dashboard/Devices.vue"),
        meta: {
          title: "Devices",
        },
      },
      {
        path: "accounts",
        name: "Accounts",
        component: () => import("@/views/dashboard/Accounts.vue"),
        meta: {
          title: "Accounts",
        },
      },
      {
        path: "namespaces",
        name: "Namespaces",
        component: () => import("@/views/dashboard/Namespaces.vue"),
        meta: {
          title: "Namespaces",
        },
      },
      {
        path: "media",
        name: "Media",
        component: () => import("@/views/dashboard/Media.vue"),
        meta: {
          title: "Media",
        },
      },
      {
        path: "plugins",
        name: "Plugins",
        component: () => import("@/views/dashboard/Plugins.vue"),
        meta: {
          title: "Plugins",
        },
      },
    ],
  },
  {
    path: "/settings",
    name: "Settings",
    component: () => import("@/views/Settings.vue"),
    meta: {
      title: "Settings",
      requiresAuth: true,
    },
    children: [
      {
        path: "",
        name: "SettingsMain",
        component: () => import("@/views/settings/index.vue"),
        meta: {
          title: "Settings",
        },
      },
      {
        path: "profile",
        name: "Profile",
        component: () => import("@/views/settings/Profile.vue"),
        meta: {
          title: "Profile",
        },
      },
      {
        path: "tokens",
        name: "Tokens",
        component: () => import("@/views/settings/Tokens.vue"),
        meta: {
          title: "Tokens",
        },
      },
      {
        path: "credentials",
        name: "Credentials",
        component: () => import("@/views/settings/Credentials.vue"),
        meta: {
          title: "Credentials",
        },
      },
    ]
  },
  {
    path: "/offline",
    name: "Offline",
    component: () => import("@/views/Offline.vue"),
    meta: {
      requiresAuth: true,
    },
  },
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

router.afterEach((to, from) => {
  nextTick(() => {
    document.title = [PLATFORM_NAME, to.meta.title].join(" | ");
  });
});

router.beforeEach(async (to, from) => {
  const store = useAppStore();
  if (to.matched.some((el) => el.meta.requiresAuth) && !store.logged_in) {
    return { name: "Login" };
  }
});

export default router;
