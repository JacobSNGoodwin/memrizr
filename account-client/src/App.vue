<template>
  <h1>Scaffolded App Works Well!</h1>
  <h3 v-if="errorCode">Error code: {{ errorCode }}</h3>
  <h3 v-if="errorMessage">{{ errorMessage }}</h3>
</template>

<script>
import { defineComponent, ref, onMounted } from 'vue';

// Wrapping exported object in define component
// gives us typing help! Woot!
export default defineComponent({
  name: 'App',
  setup() {
    const errorCode = ref(null);
    const errorMessage = ref(null);

    onMounted(async () => {
      const response = await fetch('/api/account/me', {
        method: 'GET',
      });

      const body = await response.json();

      errorCode.value = response.status;
      errorMessage.value = body.error.message;
    });

    return {
      errorCode,
      errorMessage,
    };
  },
});
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>