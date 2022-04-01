import { computed } from 'vue'
import { createRouter, createWebHashHistory } from "vue-router"
import { useAppStore } from "@/store/app";

const routes = [
  {
    path: "/login",
    name: "Login",
    component: () => import("@/views/Login.vue"),
    meta: {
      requiresAuth: false
    }
  },
  { path: "/", 
    name: "Root",
    redirect: { name: "Dashboard" },
    meta: {
      requiresAuth: true,
    }
  },
  {
    path: "/dashboard",
    name: "Dashboard",
    component: () => import("@/views/Dashboard.vue"),
    meta: {
      requiresAuth: true,
    }
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

router.beforeEach(async (to, from) => {
  const store = useAppStore()
  if (to.matched.some((el) => el.meta.requiresAuth) && !store.logged_in) {
    return { name: "Login" };
  }
});

export default router;
