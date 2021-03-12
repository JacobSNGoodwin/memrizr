<template>
  <h1 class="text-3xl text-center my-8">User Account</h1>
  <Loader
    v-if="loading"
    :height="256"
    class="animate-spin stroke-current text-blue-500 mx-auto"
  />
  <div v-else class="mb-4">
    <div
      v-if="!imageSrc"
      class="bg-gray-200 w-32 h-32 rounded-full mx-auto"
    ></div>
    <img
      v-else
      :src="imageSrc"
      class="w-32 h-32 rounded-full mx-auto"
      alt="User profile image"
    />
    <button
      type="button"
      class="btn btn-red w-32 block mx-auto my-2"
      @click="deleteImage"
    >
      Delete
    </button>
    <UserForm
      v-if="user"
      :user="user"
      @imageSubmitted="handleImageSubmitted"
      @detailsSubmitted="handleDetailsSubmitted"
    />
    <button
      type="button"
      class="btn btn-red w-32 block mx-auto my-2"
      @click="signout"
    >
      Signout
    </button>
    <p v-if="meError" class="text-center text-red-500">Error fetching user</p>
    <p v-if="postImageError" class="text-center text-red-500">
      Failed to update image
    </p>
    <p v-if="deleteImageError" class="text-center text-red-500">
      Failed to delete image
    </p>
    <p v-if="updateDetailsError" class="text-center text-red-500">
      Failed to update user details
    </p>
  </div>
</template>

<script>
import { defineComponent, computed, ref, watch } from 'vue';
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
    const { idToken, signout } = useAuth();
    const imageSrc = ref(null);

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

    const {
      exec: postImage,
      data: postImageData,
      error: postImageError,
      loading: postImageLoading,
    } = useRequest({
      url: '/api/account/image',
      method: 'post',
      headers: {
        Authorization: `Bearer ${idToken.value}`,
        'Content-Type': 'multipart/form-data',
      },
    });

    const {
      exec: deleteImage,
      data: deleteImageData,
      error: deleteImageError,
      loading: deleteImageLoading,
    } = useRequest({
      url: '/api/account/image',
      method: 'delete',
      headers: {
        Authorization: `Bearer ${idToken.value}`,
      },
    });

    const {
      exec: updateDetails,
      data: updateDetailsData,
      error: updateDetailsError,
      loading: updateDetailsLoading,
    } = useRequest({
      url: '/api/account/details',
      method: 'put',
      headers: {
        Authorization: `Bearer ${idToken.value}`,
      },
    });

    const handleDetailsSubmitted = (userDetails) => {
      updateDetails(userDetails);
    };

    const handleImageSubmitted = async (imageBlob) => {
      const formData = new FormData();
      formData.append('imageFile', imageBlob);
      postImage(formData);
    };

    watch(meData, (newData) => {
      if (!newData) return;

      if (!newData?.user?.imageUrl.length) {
        // empty imageUrl comes as empty string from API
        imageSrc.value = null;
      } else {
        imageSrc.value = newData.user.imageUrl + `?_${Date.now()}`;
      }
    });

    watch(postImageData, (newData) => {
      if (!newData) return;

      imageSrc.value = newData?.imageUrl + `?_${Date.now()}`;
    });

    watch(deleteImageData, (newData) => {
      if (!newData) return;

      if (newData.message === 'success') {
        imageSrc.value = null;
      }
    });

    const loading = computed(() => {
      return (
        meLoading.value ||
        postImageLoading.value ||
        deleteImageLoading.value ||
        updateDetailsLoading.value
      );
    });

    const user = computed(() => {
      return updateDetailsData?.value?.user || meData?.value?.user;
    });

    return {
      meData,
      meError,
      postImage,
      postImageData,
      postImageError,
      handleImageSubmitted,
      loading,
      imageSrc,
      deleteImage,
      deleteImageError,
      handleDetailsSubmitted,
      updateDetailsError,
      user,
      signout,
    };
  },
});
</script>