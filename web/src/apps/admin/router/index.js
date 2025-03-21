import {createMemoryHistory, createRouter, createWebHistory} from 'vue-router';
import {config} from "../../../services/utils.js";
import Home from "./../pages/Home.vue"
import Clients from "./../pages/Clients.vue"
import CreateClient from "./../pages/CreateClient.vue"
import EditClient from "./../pages/EditClient.vue"
import Users from "./../pages/Users.vue"
import CreateUser from "./../pages/CreateUser.vue"
import EditUser from "./../pages/EditUser.vue"
import Sessions from "./../pages/Sessions.vue"
import Session from "./../pages/Session.vue"
import NotFound from "./../pages/NotFound.vue"

const router = createRouter({
  history: config('PROD') ? createWebHistory() : createMemoryHistory(),
  routes: [
    {
      path: '/admin',
      name: 'home',
      component: Home,
      meta: {sider: 'home'},
    }, {
      path: '/admin/clients',
      name: 'clients',
      component: Clients,
      meta: {sider: 'clients'},
    }, {
      path: '/admin/clients/create',
      name: 'create-client',
      component: CreateClient,
      meta: {sider: 'clients'},
    }, {
      path: '/admin/clients/:id',
      name: 'edit-client',
      component: EditClient,
      props: true,
      meta: {sider: 'clients'},
    }, {
      path: '/admin/users',
      name: 'users',
      component: Users,
      meta: {sider: 'users'},
    }, {
      path: '/admin/users/create',
      name: 'create-user',
      component: CreateUser,
      meta: {sider: 'users'},
    }, {
      path: '/admin/users/:id',
      name: 'edit-user',
      component: EditUser,
      props: true,
      meta: {sider: 'users'},
    }, {
      path: '/admin/sessions',
      name: 'sessions',
      component: Sessions,
      meta: {sider: 'sessions'},
    }, {
      path: '/admin/sessions/:id',
      name: 'session',
      component: Session,
      props: true,
      meta: {sider: 'sessions'},
    }, {
      path: '/:pathMatch(.*)*',
      component: NotFound
    }
  ]
})

export default router
