import { createApp } from 'vue';
import App from './App.vue';
import { createAuthStore } from '../src/store/auth';
import router from './routes';
import './validators';
import './index.css';

const authStore = createAuthStore({
  onAuthRoute: '/',
  requireAuthRoute: '/authenticate',
});

createApp(App).use(authStore).use(router).mount('#app');
