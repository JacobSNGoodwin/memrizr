<template>
  <h1 class="text-4xl font-bol text-center">{{ currentUser }}</h1>
</template>

<script>
import { defineComponent } from 'vue';
import { useAuth } from '../store/auth';

// Wrapping exported object in define component
// gives us typing help! Woot!
export default defineComponent({
  name: 'Details',
  setup() {
    const { currentUser } = useAuth();

    return {
      currentUser,
    };
  },
});
</script>