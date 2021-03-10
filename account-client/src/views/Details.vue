<template>
  <UserForm v-if="meData && !meLoading" :user="meData.user" />
  <Loader
    v-if="meLoading"
    :height="256"
    class="animate-spin stroke-current text-blue-500 mx-auto"
  />
  <p v-if="meError" class="text-center text-red-500">Error fetching user</p>
</template>

<script>
import { defineComponent } from 'vue';
import { useAuth } from '../store/auth';
import { useRequest } from '../util';
import UserForm from '../components/UserForm.vue';
import Loader from '../components/ui/Loader.vue';

export default defineComponent({
  name: 'Details',
  components: {
    UserForm,
    Loader,
  },
  setup() {
    const { idToken } = useAuth();

    const { data: meData, error: meError, loading: meLoading } = useRequest(
      {
        url: '/api/account/me',
        method: 'get',
        headers: {
          Authorization: `Bearer ${idToken.value}`,
        },
      },
      {
        execOnMounted: true,
      }
    );

    return {
      meData,
      meError,
      meLoading,
    };
  },
});
</script>