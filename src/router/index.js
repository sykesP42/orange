import { createRouter, createWebHistory } from 'vue-router'
import { getActivePinia } from 'pinia'
import { useUserStore } from '../stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Login.vue'),
    meta: { title: '注册' }
  },
  {
    path: '/',
    component: () => import('../views/Layout.vue'),
    children: [
      {
        path: '',
        name: 'Home',
        component: () => import('../views/Home.vue'),
        meta: { title: '首页' }
      },
      {
        path: 'my',
        name: 'My',
        component: () => import('../views/My.vue'),
        meta: { title: '个人中心' }
      },
      {
        path: 'add',
        name: 'Add',
        component: () => import('../views/Add.vue'),
        meta: { title: '录入博主' }
      },
      {
        path: 'my-bloggers',
        name: 'MyBloggers',
        component: () => import('../views/MyBloggers.vue'),
        meta: { title: '我的博主' }
      },
      {
        path: 'statistics',
        name: 'Statistics',
        component: () => import('../views/Statistics.vue'),
        meta: { title: '数据统计' }
      },
      {
        path: 'calendar',
        name: 'Calendar',
        component: () => import('../views/CalendarPage.vue'),
        meta: { title: '提醒日历' }
      },
      {
        path: 'team',
        name: 'Team',
        component: () => import('../views/Team.vue'),
        meta: { title: '团队中心' }
      },
      {
        path: 'admin',
        name: 'Admin',
        component: () => import('../views/Admin.vue'),
        meta: { requiresAdmin: true, title: '管理后台' }
      },
      {
        path: 'team/:teamId',
        name: 'TeamHome',
        component: () => import('../views/TeamHome.vue'),
        meta: { title: '小组主页' }
      },
      {
        path: 'team/:teamId/post/:postId',
        name: 'SharedPost',
        component: () => import('../views/TeamHome.vue'),
        meta: { title: '帖子详情' }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('../views/PublicUsers.vue'),
        meta: { title: '用户' }
      },
      {
        path: 'log',
        name: 'Log',
        component: () => import('../views/Log.vue'),
        meta: { title: '操作日志', requiresAdmin: true }
      },
      {
        path: 'analytics',
        name: 'Analytics',
        component: () => import('../views/AnalyticsDashboard.vue'),
        meta: { title: '数据分析' }
      },
      {
        path: 'blogger/:id',
        name: 'BloggerDetail',
        component: () => import('../views/BloggerDetail.vue'),
        meta: { title: '博主详情' }
      },
      {
        path: 'user/:username',
        name: 'UserDetail',
        component: () => import('../views/UserDetail.vue'),
        meta: { title: '用户详情' }
      },
      {
        path: 'user-settings',
        name: 'UserSettings',
        component: () => import('../views/UserSettings.vue'),
        meta: { title: '个人设置' }
      },
      {
        path: 'invalid-bloggers',
        name: 'InvalidBloggers',
        component: () => import('../views/InvalidBloggers.vue'),
        meta: { title: '失效博主库' }
      },
      {
        path: 'chat',
        name: 'Chat',
        component: () => import('../views/Chat.vue'),
        meta: { title: '私信' }
      },
      {
        path: 'workflow',
        name: 'Workflow',
        component: () => import('../views/Workflow.vue'),
        meta: { title: '工作流自动化' }
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
    }
    return { top: 0, behavior: 'smooth' }
  }
})

router.beforeEach(async (to, from, next) => {
  if (!getActivePinia()) {
    next()
    return
  }
  const userStore = useUserStore()
  const token = localStorage.getItem('token')

  if (token && !userStore.token) {
    userStore.token = token
    userStore.username = localStorage.getItem('username') || ''
    userStore.realName = localStorage.getItem('realName') || ''
    userStore.role = localStorage.getItem('role') || 'user'
    userStore.id = localStorage.getItem('id') ? parseInt(localStorage.getItem('id')) : null
    userStore.team_id = localStorage.getItem('team_id') ? parseInt(localStorage.getItem('team_id')) : null
    userStore.team_name = localStorage.getItem('team_name') || null
    userStore.team_color = localStorage.getItem('team_color') || null
    userStore.isLoggedIn = true
    userStore.isAdmin = localStorage.getItem('role') === 'admin'
  }

  if (to.path !== '/login' && !token) {
    next('/login')
  } else if (to.meta.requiresAdmin) {
    if (userStore.role !== 'admin') {
      next('/')
    } else {
      if (!userStore._roleVerified) {
        try {
          const api = (await import('../api')).default
          const res = await api.get('/user/profile')
          if (res.code === 200 && res.data?.role === 'admin') {
            userStore._roleVerified = true
            next()
          } else {
            userStore.role = res.data?.role || 'user'
            userStore.isAdmin = false
            localStorage.setItem('role', userStore.role)
            next('/')
          }
        } catch {
          next('/')
        }
      } else {
        next()
      }
    }
  } else {
    next()
  }
})

router.afterEach((to) => {
  const title = to.meta?.title
  document.title = title ? `${title} - Orange` : 'Orange'
})

export default router
