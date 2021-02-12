import { createRouter, createWebHistory } from 'vue-router';

import Auth from './views/Auth.vue';
import Details from './views/Details.vue';
import NotFound from './views/NotFound.vue';

const routes = [
  {
    path: '/authenticate',
    name: 'Auth',
    component: Auth,
  },
  {
    path: '/',
    name: 'Details',
    component: Details,
    // beforeEnter: requireAuth,
  },
  {
    path: '/:catchAll(.*)*',
    name: 'NotFound',
    component: NotFound,
  },
];

const router = createRouter({
  // env variable provided base "base" key of vite config
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

export default router;
