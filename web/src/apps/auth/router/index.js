import {createRouter, createMemoryHistory, createWebHistory} from 'vue-router';
import {config} from "../../../services/utils.js";
import Authorize from "../pages/Authorize.vue";
import ForgotPassword from "../pages/ForgotPassword.vue";
import ResetPassword from "../pages/ResetPassword.vue";
import PageNotFound from "../pages/PageNotFound.vue";

const router = createRouter({
  history: config('PROD') ? createWebHistory() :createMemoryHistory(),
  routes: [
    {
      path: '/',
      component: Authorize,
    },
    {
      path: '/oauth/authorize',
      name: 'authorize',
      component: Authorize,
    }, {
      path: '/oauth/forgot-password',
      name: 'forgot-password',
      component: ForgotPassword,
    }, {
      path: '/oauth/reset-password',
      name: 'reset-password',
      component: ResetPassword,
    }, {
      path: '/:pathMatch(.*)*',
      component: PageNotFound
    }
  ]
})

export default router
