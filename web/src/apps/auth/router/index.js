import {createRouter, createMemoryHistory, createWebHistory} from 'vue-router';
import {config} from "../../../services/utils.js";
import Authorize from "../pages/Authorize.vue";
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
      path: '/:pathMatch(.*)*',
      component: PageNotFound
    }
  ]
})

export default router
