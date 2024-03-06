import Vue from 'vue'
import App from './App.vue'
import VueRouter from 'vue-router'
import ItemList from './components/ItemList.vue'
import ItemDetail from './components/ItemDetail.vue'

Vue.use(VueRouter)

const routes = [
  { path: '/', component: ItemList },
  { path: '/item/:id', component: ItemDetail }
]

const router = new VueRouter({
  routes
})

new Vue({
  render: h => h(App),
  router
}).$mount('#app')
