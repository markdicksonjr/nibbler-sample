import Vue from 'vue'
import Router from 'vue-router'
import Resource from 'vue-resource'
import ResetPassword from './views/ResetPassword.vue'
import ForgotPassword from './views/ForgotPassword.vue'
import Home from './views/Home.vue'
import Login from './views/Login.vue'
import Logout from './views/Logout.vue'
import VerifyEmail from './views/VerifyEmail'

Vue.use(Router)
Vue.use(Resource)

const router = new Router({
  mode: 'hash',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/register',
      name: 'register',
      component: () => import(/* webpackChunkName: "register" */ './views/Register')
    },
    {
      path: '/verify-email',
      name: 'verify-email',
      component: VerifyEmail
    },
    {
      path: '/forgot-password',
      name: 'forgot-password',
      component: ForgotPassword
    },
    {
      path: '/reset-password',
      name: 'reset-password',
      component: ResetPassword
    },
    {
      path: '/login',
      name: 'login',
      component: Login
    },
    {
      path: '/logout',
      name: 'logout',
      component: Logout
    },
    {
      path: '/about',
      name: 'about',

      // route level code-splitting
      // this generates a separate chunk (about.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import(/* webpackChunkName: "about" */ './views/About.vue')
    }
  ]
})

router.beforeEach((to, from, next) => {

  // redirect to login page if not logged in and trying to access a restricted page
  const publicPages = ['/login', '/register', '/forgot-password']
  const authRequired = !publicPages.includes(to.path)
  const user = localStorage.getItem('user')

  if(authRequired) {
    if(!user) {
      return next('/login')
    }

    try {
      const userDetails = JSON.parse(user)
      const loginTime = parseInt(localStorage.getItem('lastLoginTime'))
      const maxTime = loginTime + userDetails.sessionAgeSeconds * 1000

      // if the session has expired
      if(isNaN(loginTime) || maxTime < new Date().valueOf()) {
        localStorage.removeItem('user')
        return next('/login')
      }
    } catch(ex) {
      localStorage.removeItem('user')
      return next('/login')
    }
  }

  next();
})

export default router