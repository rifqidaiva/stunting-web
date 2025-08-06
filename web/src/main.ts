import { createApp } from "vue"
import { createRouter, createWebHistory } from "vue-router"

import "./style.css"
import App from "./App.vue"

import Home from "./components/Home.vue"
import Admin from "./components/admin/Admin.vue"
import Community from "./components/community/Community.vue"

const routes = [
  {
    path: "/",
    component: Home,
  },
  {
    path: "/admin",
    component: Admin,
  },
  {
    path: "/community",
    component: Community,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

createApp(App).use(router).mount("#app")
