import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/store/auth'

// 布局组件
const DefaultLayout = () => import('@/layouts/DefaultLayout.vue')
const AuthLayout = () => import('@/layouts/AuthLayout.vue')
const AdminLayout = () => import('@/layouts/AdminLayout.vue')

// 页面组件懒加载
const AdminLogin = () => import('@/views/admin/AdminLoginPage.vue')
const Menu = () => import('@/views/MenuPage.vue')
const Cart = () => import('@/views/CartPage.vue')
const CheckoutSuccess = () => import('@/views/CheckoutSuccessPage.vue')
const AdminDashboard = () => import('@/views/admin/DashboardPage.vue')
const AdminMenu = () => import('@/views/admin/MenuManagePage.vue')
const AdminOrders = () => import('@/views/admin/OrderManagePage.vue')
const AdminProfile = () => import('@/views/admin/AdminProfilePage.vue')
const AdminSettings = () => import('@/views/admin/AdminSettingsPage.vue')
const NotFound = () => import('@/views/NotFoundPage.vue')

// 用户认证页面
const UserLogin = () => import('@/views/auth/LoginPage.vue')
const UserRegister = () => import('@/views/auth/RegisterPage.vue')

// 用户中心页面
const ProfilePage = () => import('@/views/user/ProfilePage.vue')
const MemberCenterPage = () => import('@/views/user/MemberCenterPage.vue')
const OrdersPage = () => import('@/views/user/OrdersPage.vue')

const routes = [
  {
    path: '/',
    component: DefaultLayout,
    children: [
      // 根路径重定向
      {
        path: '',
        redirect: '/menu'
      },
      // 用户界面路由（无需登录）
      {
        path: 'menu',
        name: 'Menu',
        component: Menu,
        meta: {
          title: '菜单',
          requiresAuth: false
        }
      },
      {
        path: 'cart',
        name: 'Cart',
        component: Cart,
        meta: {
          title: '购物车',
          requiresAuth: false
        }
      },
      {
        path: 'checkout-success',
        name: 'CheckoutSuccess',
        component: CheckoutSuccess,
        meta: {
          title: '结算成功',
          requiresAuth: false
        }
      }
    ]
  },
  // 用户认证路由
  {
    path: '/auth',
    component: AuthLayout,
    children: [
      {
        path: 'login',
        name: 'UserLogin',
        component: UserLogin,
        meta: {
          title: '用户登录',
          requiresAuth: false,
          guestOnly: true
        }
      },
      {
        path: 'register',
        name: 'UserRegister',
        component: UserRegister,
        meta: {
          title: '用户注册',
          requiresAuth: false,
          guestOnly: true
        }
      }
    ]
  },
  // 用户中心路由
  {
    path: '/user',
    component: DefaultLayout,
    children: [
      {
        path: 'profile',
        name: 'Profile',
        component: ProfilePage,
        meta: {
          title: '个人中心',
          requiresAuth: false // 页面内部处理未登录状态
        }
      },
      {
        path: 'member',
        name: 'MemberCenter',
        component: MemberCenterPage,
        meta: {
          title: '会员中心',
          requiresAuth: false // 页面内部处理未登录状态
        }
      },
      {
        path: 'orders',
        name: 'UserOrders',
        component: OrdersPage,
        meta: {
          title: '我的订单',
          requiresAuth: false // 页面内部处理未登录状态
        }
      }
    ]
  },
  // 管理员登录页面（独立路由）
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: AdminLogin,
    meta: {
      title: '管理员登录',
      requiresAuth: false
    }
  },
  // 管理员后台
  {
    path: '/admin',
    component: AdminLayout,
    children: [
      {
        path: '',
        name: 'AdminDashboard',
        component: AdminDashboard,
        meta: {
          title: '管理后台',
          requiresAuth: true,
          requiresAdmin: true
        }
      },
      {
        path: 'menu',
        name: 'AdminMenu',
        component: AdminMenu,
        meta: {
          title: '菜单管理',
          requiresAuth: true,
          requiresAdmin: true
        }
      },
      {
        path: 'orders',
        name: 'AdminOrders',
        component: AdminOrders,
        meta: {
          title: '订单管理',
          requiresAuth: true,
          requiresAdmin: true
        }
      },
      {
        path: 'profile',
        name: 'AdminProfile',
        component: AdminProfile,
        meta: {
          title: '个人信息',
          requiresAuth: true,
          requiresAdmin: true
        }
      },
      {
        path: 'settings',
        name: 'AdminSettings',
        component: AdminSettings,
        meta: {
          title: '系统设置',
          requiresAuth: true,
          requiresAdmin: true
        }
      }
    ]
  },
  // 404页面
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound,
    meta: {
      title: '页面不存在'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// 全局前置守卫
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - ${import.meta.env.VITE_APP_TITLE || '咖啡点餐系统'}`
  }

  // 检查是否需要管理员权限
  if (to.meta.requiresAdmin) {
    // 先尝试从localStorage恢复登录状态
    if (!authStore.isAuthenticated) {
      await authStore.checkAuthStatus()
    }

    // 如果仍未登录，跳转到管理员登录页
    if (!authStore.isAuthenticated) {
      next({
        name: 'AdminLogin',
        query: { redirect: to.fullPath }
      })
      return
    }

    // 如果不是管理员，跳转到首页
    if (!authStore.isAdmin) {
      next({ name: 'Menu' })
      return
    }
  }

  // 如果已登录管理员访问管理员登录页，跳转到管理后台
  if (authStore.isAuthenticated && authStore.isAdmin && to.name === 'AdminLogin') {
    next({ name: 'AdminDashboard' })
    return
  }

  // 检查是否为仅限未登录用户访问的页面（用户登录/注册）
  if (to.meta.guestOnly && authStore.isAuthenticated && !authStore.isAdmin) {
    next({ name: 'Menu' })
    return
  }

  next()
})

// 全局后置钩子
router.afterEach((to, from) => {
  // 这里可以添加页面访问统计等逻辑
  console.log(`从 ${from.path} 导航到 ${to.path}`)
})

export default router
