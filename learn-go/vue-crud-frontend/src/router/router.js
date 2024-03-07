// router/index.js
import { createRouter, createWebHistory } from 'vue-router';
import ItemList from '../components/ItemList.vue';
import ItemDetail from '../components/ItemDetail.vue';

// Use createRouter here, not VueRouter
const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: ItemList },
    { path: '/item/:id', component: ItemDetail, props: true },
  ],
});

export default router;
