import { createRouter, createWebHistory } from 'vue-router'
// import HomeView from '../views/HomeView.vue'
import Index from '@/views/index.vue'
import Login from '@/views/login.vue'
import About from '@/views/about.vue'
import NotFound from '@/views/404.vue'
import Register from '@/views/register.vue'
import Verify from '@/views/verify.vue'

const routes = [
  {
    path: '/',
    component: Index
  },
  {
    path: '/about',
    component: About
  },
  {
    path: '/login',
    component: Login
  },
  {
    path: '/register',
    component: Register
  },
  {
    path: '/user/verify',
    component: Verify
  },
  {
    path: '/:pathMatch(.*)*',
    name: NotFound,
    component: NotFound
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
