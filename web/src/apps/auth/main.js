import { createApp } from 'vue';
import { createPinia } from "pinia";
import App from './App.vue';
import router from './router';
import PrimeVue from 'primevue/config';

import ToastService from 'primevue/toastservice';

import Button from 'primevue/button';
import Checkbox from 'primevue/checkbox'
import InputText from 'primevue/inputtext';
import Message from "primevue/message";
import InlineMessage from 'primevue/inlinemessage';
import Password from 'primevue/password';
import Toast from "primevue/toast";

import './../../style/auth/main.scss'

try {
    window.initData = JSON.parse(document.head.querySelector("meta[name='init-data']").getAttribute("content"))
    document.head.querySelector("meta[name='init-data']").remove()
} catch (e) {
    window.initData = {}
}

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);
app.use(PrimeVue);
app.use(ToastService);

app.component('Button', Button)
app.component('Checkbox', Checkbox)
app.component('InputText', InputText)
app.component('Message', Message)
app.component('InlineMessage', InlineMessage)
app.component('Password', Password)
app.component('Toast', Toast)

app.mount('#app');