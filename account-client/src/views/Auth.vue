<template>
  <div class="max-w-xl mx-auto px-4">
    <div class="rounded-lg shadow-lg p-4">
      <ul class="mx-8 mb-2 flex justify-center">
        <li
          @click="setIsLogin(true)"
          class="mx-2 cursor-pointer text-center hover:opacity-75 transition-opacity"
          :class="{ 'border-b-2 border-blue-400': isLogin }"
        >
          Login
        </li>
        <li
          @click="setIsLogin(false)"
          class="mx-2 cursor-pointer text-center hover:opacity-75 transition-opacity"
          :class="{ 'border-b-2 border-blue-400': !isLogin }"
        >
          Sign Up
        </li>
      </ul>
      <LoginForm :isLogin="isLogin" @submitAuth="authSubmitted" />
      <!-- Submission errors (when we connect to the API) -->
    </div>
  </div>
</template>

<script>
import { defineComponent, ref } from 'vue';
import LoginForm from '../components/LoginForm.vue';
export default defineComponent({
  name: 'Auth',
  components: {
    LoginForm,
  },
  setup() {
    const isLogin = ref(true);

    const setIsLogin = (nextVal) => {
      isLogin.value = nextVal;
    };

    const authSubmitted = ({ email, password }) => {
      console.log(isLogin, email, password);
    };

    return {
      isLogin,
      setIsLogin,
      authSubmitted,
    };
  },
});
</script>