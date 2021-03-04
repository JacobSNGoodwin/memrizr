<template>
  <div class="mx-8 my-8">
    <Loader
      v-if="isLoading"
      :height="512"
      class="animate-spin stroke-current text-blue-500 mx-auto"
    />
    <router-view v-else></router-view>
  </div>
</template>

<script>
import { defineComponent, onMounted } from 'vue';
import Loader from './components/ui/Loader.vue';
import { useAuth } from './store/auth';

// Wrapping exported object in define component
// gives us typing help! Woot!
export default defineComponent({
  name: 'App',
  components: {
    Loader,
  },
  setup() {
    const { initializeUser, isLoading } = useAuth();

    onMounted(() => {
      initializeUser();
    });

    return {
      isLoading,
    };
  },
});
</script>