// router/index.js
import Vue from 'vue';
import VueRouter from 'vue-router';
import ItemList from '../components/ItemList.vue';
import ItemDetail from '../components/ItemDetail.vue';

Vue.use(VueRouter);

const routes = [
  { path: '/', component: ItemList },
  { path: '/item/:id', component: ItemDetail, props: true },
];

const router = new VueRouter({
  routes,
});

export default router;
