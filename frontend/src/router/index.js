import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录', requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { title: '注册', requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/components/Layout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { title: '工作台', icon: 'Odometer' }
      },
      {
        path: 'works',
        name: 'MyWorks',
        component: () => import('@/views/works/MyWorks.vue'),
        meta: { title: '我的作品', icon: 'Headset' }
      },
      {
        path: 'works/register',
        name: 'RegisterWork',
        component: () => import('@/views/works/RegisterWork.vue'),
        meta: { title: '版权存证', icon: 'DocumentAdd' }
      },
      {
        path: 'works/search',
        name: 'SearchWorks',
        component: () => import('@/views/works/SearchWorks.vue'),
        meta: { title: '作品搜索', icon: 'Search' }
      },
      {
        path: 'works/:workID',
        name: 'WorkDetail',
        component: () => import('@/views/works/WorkDetail.vue'),
        meta: { title: '作品详情', icon: 'View', hidden: true }
      },
      {
        path: 'works/:workID/verify-hash',
        name: 'VerifyHash',
        component: () => import('@/views/works/VerifyHash.vue'),
        meta: { title: '哈希验真', icon: 'Check', hidden: true }
      },
      {
        path: 'works/:workID/transfer',
        name: 'Transfer',
        component: () => import('@/views/works/Transfer.vue'),
        meta: { title: '版权转让', icon: 'Share', hidden: true }
      },
      {
        path: 'works/:workID/certificate',
        name: 'Certificate',
        component: () => import('@/views/works/Certificate.vue'),
        meta: { title: '存证证书', icon: 'Medal', hidden: true }
      },
      {
        path: 'works/:workID/dispute',
        name: 'FileDispute',
        component: () => import('@/views/works/FileDispute.vue'),
        meta: { title: '争议存证', icon: 'Warning', hidden: true }
      },
      {
        path: 'licenses/grant',
        name: 'GrantLicense',
        component: () => import('@/views/licenses/GrantLicense.vue'),
        meta: { title: '发放授权', icon: 'Tickets' }
      },
      {
        path: 'licenses/verify',
        name: 'VerifyLicense',
        component: () => import('@/views/licenses/VerifyLicense.vue'),
        meta: { title: '授权核验', icon: 'CircleCheck' }
      },
      {
        path: 'licenses/my',
        name: 'MyLicenses',
        component: () => import('@/views/licenses/MyLicenses.vue'),
        meta: { title: '我的授权', icon: 'Key' }
      },
      {
        path: 'disputes',
        name: 'DisputeList',
        component: () => import('@/views/disputes/DisputeList.vue'),
        meta: { title: '争议记录', icon: 'Warning' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')

  if (to.meta.requiresAuth !== false && !token) {
    next({ path: '/login', query: { redirect: to.fullPath } })
  } else if ((to.path === '/login' || to.path === '/register') && token) {
    next({ path: '/' })
  } else {
    next()
  }
})

router.afterEach((to) => {
  document.title = to.meta.title ? `${to.meta.title} - 音乐版权保护` : '音乐数字版权保护平台'
})

export default router
