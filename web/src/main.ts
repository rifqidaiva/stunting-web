import { createApp } from "vue"
import { createRouter, createWebHistory } from "vue-router"

import "./style.css"
import App from "./App.vue"

import Home from "./components/Home.vue"
import Admin from "./components/admin/Admin.vue"
import Report from "./components/report/Report.vue"

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
    path: "/report",
    component: Report,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

createApp(App).use(router).mount("#app")
