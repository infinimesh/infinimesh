import { createRouter, createWebHashHistory } from "vue-router"

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
    name: "Home",
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

let isLoggedIn = false

router.beforeEach(async (to, from) => {
  if (to.matched.some((el) => el.meta.requiresAuth) && !isLoggedIn) {
    return { name: "Login" };
  }
});

export default router;
