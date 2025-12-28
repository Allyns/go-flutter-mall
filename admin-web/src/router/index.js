import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/Login.vue')
    },
    {
      path: '/',
      name: 'home',
      component: () => import('../views/Home.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'dashboard',
          component: () => import('../views/Dashboard.vue')
        },
        {
          path: 'products',
          name: 'products',
          component: () => import('../views/products/ProductList.vue')
        },
        {
          path: 'orders',
          name: 'orders',
          component: () => import('../views/orders/OrderList.vue')
        },
        {
          path: 'chat',
          name: 'chat',
          component: () => import('../views/chat/ChatWindow.vue')
        },
        {
          path: 'notifications',
          name: 'notifications',
          component: () => import('../views/notifications/NotificationList.vue')
        }
      ]
    }
  ]
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else {
    next()
  }
})

export default router
