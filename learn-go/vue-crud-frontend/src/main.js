import { createApp } from 'vue';
import App from './app.vue';
import router from './router/router';
import axios from './axios';

const app = createApp(App);

// Menggunakan Axios dalam aplikasi
app.config.globalProperties.$axios = axios;

// Menggunakan router dengan instance aplikasi
app.use(router);

// Memasang aplikasi ke elemen '#app'
app.mount('#app');
