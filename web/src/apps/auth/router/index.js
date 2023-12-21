import { createRouter, createMemoryHistory } from 'vue-router';
import Signin from "../views/Signin.vue";
import ForgotPassword from "../views/ForgotPassword.vue";

export default createRouter({
    history: createMemoryHistory(),
    routes: [
        {
            path: '/',
            name: 'signin',
            component: Signin,
        }, {
            path: '/forgot-password',
            name: 'forgot-password',
            component: ForgotPassword,
        }
    ]
});
