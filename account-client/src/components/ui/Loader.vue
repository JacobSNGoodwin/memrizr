<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    :height="height"
    viewBox="0 0 512 512"
    stroke="black"
  >
    <g transform="translate(256, 256)">
      <circle r="226" stroke-width="60" fill="none" style="opacity: 0.3" />
      <path
        d="M -159.806 -159.806 A 226 226 0 0 1 159.806 -159.806"
        stroke-width="60"
        fill="none"
        stroke-linecap="round"
      />
    </g>
  </svg>
</template>

<script>
import { defineComponent } from 'vue';
export default defineComponent({
  name: 'Loader',
  props: {
    height: {
      type: Number,
      default: 256,
    },
  },
});
</script>