import Vue from 'vue'
import App from './app.vue'
import router from './router/router'; // Sesuaikan dengan lokasi rute Anda

new Vue({
  render: h => h(App),
  router
}).$mount('#app');