<template>
  <div class="text-4xl font-bol text-center my-2">App Component</div>
  <router-view></router-view>
</template>

<script>
import { defineComponent } from 'vue';
import { provideAuth } from './store/auth';

// Wrapping exported object in define component
// gives us typing help! Woot!
export default defineComponent({
  name: 'App',
  setup() {
    provideAuth();
  },
});
</script>