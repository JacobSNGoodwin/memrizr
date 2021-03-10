<template>
  <div class="max-w-3xl mx-auto">
    <h1 class="text-3xl text-center">Update User</h1>
    <ImageCropper @imageCropped="$emit('imageSubmitted', $event)" />
    <div class="max-w-xl px-4 mx-auto">
      <div class="my-1">
        <label for="email" class="font-semibold">Email Address</label>
        <input
          type="text"
          name="email"
          placeholder="Email Address"
          :value="emailValue"
          @input="handleEmailChange"
          @blur="handleEmailBlur"
          class="px-4 my-2 min-w-full mx-auto border border-gray-500 rounded-full focus:outline-none focus:ring-1 focus:border-blue-300"
        />
        <div class="h-6">
          <p
            v-show="emailMeta.touched && !emailMeta.valid"
            class="text-center text-red-500"
          >
            {{ emailErrorMessage }}
          </p>
        </div>
      </div>
      <div class="my-1">
        <label for="name" class="font-semibold">Name</label>
        <input
          type="text"
          name="name"
          placeholder="Name"
          :value="nameValue"
          @input="handleNameChange"
          @blur="handleNameBlur"
          class="px-4 my-2 min-w-full mx-auto border border-gray-500 rounded-full focus:outline-none focus:ring-1 focus:border-blue-300"
        />
        <div class="h-6">
          <p
            v-show="nameMeta.touched && !nameMeta.valid"
            class="text-center text-red-500"
          >
            {{ nameErrorMessage }}
          </p>
        </div>
      </div>
      <div class="my-1">
        <label for="website" class="font-semibold">Website URL</label>
        <input
          type="text"
          name="website"
          placeholder="Website URL"
          :value="websiteValue"
          @input="handleWebsiteChange"
          @blur="handleWebsiteBlur"
          class="px-4 my-2 min-w-full mx-auto border border-gray-500 rounded-full focus:outline-none focus:ring-1 focus:border-blue-300"
        />
        <div class="h-6">
          <p
            v-show="websiteMeta.touched && !websiteMeta.valid"
            class="text-center text-red-500"
          >
            {{ websiteErrorMessage }}
          </p>
        </div>
      </div>
      <div class="mt-6">
        <button
          :disabled="!formMeta.valid || !formMeta.dirty"
          class="btn btn-blue w-32 block mx-auto"
          @click="submitDetails"
        >
          Update
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent } from 'vue';
import { useForm, useField } from 'vee-validate';
import ImageCropper from './ui/ImageCropper.vue';
export default defineComponent({
  name: 'UserForm',
  components: {
    ImageCropper,
  },
  props: {
    user: {
      type: Object,
      default: null,
    },
  },
  emits: ['detailsSubmitted', 'imageSubmitted'],
  setup(props, { emit }) {
    const { meta: formMeta, handleSubmit } = useForm();
    const {
      value: emailValue,
      meta: emailMeta,
      errorMessage: emailErrorMessage,
      handleBlur: handleEmailBlur,
      handleChange: handleEmailChange,
    } = useField('email', 'email', {
      initialValue: props.user.email,
      validateOnMount: true,
    });

    const {
      value: nameValue,
      meta: nameMeta,
      errorMessage: nameErrorMessage,
      handleBlur: handleNameBlur,
      handleChange: handleNameChange,
    } = useField('name', 'name', {
      initialValue: props.user.name,
      validateOnMount: true,
    });

    const {
      value: websiteValue,
      meta: websiteMeta,
      errorMessage: websiteErrorMessage,
      handleBlur: handleWebsiteBlur,
      handleChange: handleWebsiteChange,
    } = useField('website', 'url', {
      initialValue: props.user.website,
      validateOnMount: true,
    });

    const submitDetails = handleSubmit((formValues) => {
      emit('detailsSubmitted', {
        email: formValues.email,
        name: formValues.name,
        website: formValues.website,
      });
    });

    return {
      formMeta,
      emailValue,
      emailMeta,
      emailErrorMessage,
      handleEmailBlur,
      handleEmailChange,
      nameValue,
      nameMeta,
      nameErrorMessage,
      handleNameBlur,
      handleNameChange,
      websiteValue,
      websiteMeta,
      websiteErrorMessage,
      handleWebsiteBlur,
      handleWebsiteChange,
      submitDetails,
    };
  },
});
</script>