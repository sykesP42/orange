import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/',
    component: () => import('../views/Layout.vue'),
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('../views/Home.vue')
      },
      {
        path: 'my',
        name: 'My',
        component: () => import('../views/My.vue')
      },
      {
        path: 'add',
        name: 'Add',
        component: () => import('../views/Add.vue')
      },
      {
        path: 'statistics',
        name: 'Statistics',
        component: () => import('../views/Statistics.vue')
      },
      {
        path: 'team',
        name: 'Team',
        component: () => import('../views/Team.vue')
      },
      {
        path: 'admin',
        name: 'Admin',
        component: () => import('../views/Admin.vue'),
        meta: { requiresAdmin: true }
      },
      {
        path: 'team/:teamId',
        name: 'TeamHome',
        component: () => import('../views/TeamHome.vue')
      },
      {
        path: 'forums',
        name: 'Forums',
        component: () => import('../views/Forums.vue')
      },
      {
        path: 'log',
        name: 'Log',
        component: () => import('../views/Log.vue')
      },
      {
        path: 'blogger/:id',
        name: 'BloggerDetail',
        component: () => import('../views/BloggerDetail.vue')
      },
      {
        path: 'user/:username',
        name: 'UserDetail',
        component: () => import('../views/UserDetail.vue')
      },
      {
        path: 'invalid-bloggers',
        name: 'InvalidBloggers',
        component: () => import('../views/InvalidBloggers.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0, behavior: 'smooth' }
    }
  }
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  const token = localStorage.getItem('token')

  if (token && !userStore.token) {
    userStore.token = token
    userStore.username = localStorage.getItem('username') || ''
    userStore.realName = localStorage.getItem('realName') || ''
    userStore.role = localStorage.getItem('role') || 'user'
    userStore.team_id = localStorage.getItem('team_id') || null
  }

  if (to.path !== '/login' && !token) {
    next('/login')
  } else if (to.meta.requiresAdmin && userStore.role !== 'admin') {
    next('/')
  } else {
    next()
  }
})

export default router