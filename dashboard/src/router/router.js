import { createRouter, createWebHashHistory } from "vue-router"

const routes = [
  { path: "/", 
    name: "Home",
    redirect: { name: "Dashboard" },
    meta: {
      requireLogin: true,
    }
  },
  {
    path: "/dashboard",
    name: "Dashboard",
    component: () => import("@/views/Dashboard.vue"),
    meta: {
      requireLogin: true,
    }
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router;
