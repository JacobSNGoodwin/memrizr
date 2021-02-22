<template>
  <div class="px-4">
    <h1 class="text-xl text-center font-semibold mb-2">
      {{ isLogin ? 'Login' : 'Sign Up' }}
    </h1>
    <input
      type="text"
      placeholder="Email Address"
      @input="emailField.handleChange"
      @blur="emailField.handleBlur"
      :value="emailField.value"
      class="px-4 my-2 min-w-full mx-auto border border-gray-500 rounded-full focus:outline-none focus:ring-1 focus:border-blue-300"
    />
    <p
      class="text-center text-red-500"
      :style="{
        visibility:
          emailField.meta.touched && !emailField.meta.valid
            ? 'visible'
            : 'hidden',
      }"
    >
      {{ emailField.errorMessage || 'Field is Required' }}
    </p>
    <input
      type="password"
      placeholder="Password"
      @input="
        passwordField.handleChange($event), confirmPasswordField.validate()
      "
      @blur="passwordField.handleBlur"
      :value="passwordField.value"
      class="px-4 my-2 min-w-full mx-auto border border-gray-500 rounded-full focus:outline-none focus:ring-1 focus:border-blue-300"
    />
    <p
      class="text-center text-red-500"
      :style="{
        visibility:
          passwordField.meta.touched && !passwordField.meta.valid
            ? 'visible'
            : 'hidden',
      }"
    >
      {{ passwordField.errorMessage || 'Field is Required' }}
    </p>

    <template v-if="!isLogin">
      <input
        type="password"
        placeholder="Confirm Password"
        @input="confirmPasswordField.handleChange"
        @blur="confirmPasswordField.handleBlur"
        :value="confirmPasswordField.value"
        class="px-4 my-2 min-w-full mx-auto border border-gray-500 rounded-full focus:outline-none focus:ring-1 focus:border-blue-300"
      />
      <p
        :style="{
          visibility:
            confirmPasswordField.meta.touched &&
            !confirmPasswordField.meta.valid
              ? 'visible'
              : 'hidden',
        }"
        class="text-center text-red-500"
      >
        {{ confirmPasswordField.errorMessage || 'Field is Required' }}
      </p>
    </template>

    <div class="flex justify-center mt-4">
      <button
        class="btn btn-blue mx-1"
        :disabled="!formMeta.valid"
        @click="submitForm"
      >
        {{ isLogin ? 'Login' : 'Sign Up' }}
      </button>
    </div>
  </div>
</template>

<script>
import { defineComponent, reactive, computed, watch } from 'vue';
import { useField, useForm } from 'vee-validate';

export default defineComponent({
  name: 'LoginForm',
  props: {
    isLogin: {
      type: Boolean,
      default: true,
    },
    isSubmitting: {
      type: Boolean,
      default: false,
    },
  },
  emits: {
    submitAuth: null, // null means we will not validate event
  },
  setup(props, { emit }) {
    const { meta: formMeta, handleSubmit } = useForm();
    const emailField = reactive(useField('email', 'email'));
    const passwordField = reactive(useField('password', 'password'));

    const confirmPasswordValidator = computed(() => {
      return !props.isLogin ? 'confirmPassword:password' : () => true;
    });

    const confirmPasswordField = reactive(
      useField('confirmPassword', confirmPasswordValidator)
    );

    watch(
      () => props.isLogin,
      () => {
        confirmPasswordField.validate();
      }
    );

    const submitForm = handleSubmit((formValues) => {
      emit('submitAuth', {
        email: formValues.email,
        password: formValues.password,
      });
    });

    return {
      emailField,
      passwordField,
      confirmPasswordField,
      submitForm,
      formMeta,
    };
  },
});
</script>