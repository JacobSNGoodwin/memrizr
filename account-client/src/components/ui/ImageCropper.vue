<template>
  <div v-show="imageSrc" class="my-2 w-64 h-64 object-fill mx-auto">
    <img class="block max-w-full" ref="img" :src="imageSrc" />
  </div>
  <div class="flex justify-center content-end mt-2">
    <button
      v-if="!imageSrc"
      class="btn btn-blue w-32 mx-2"
      @click="imageInput.click()"
    >
      Select Image
    </button>
    <button v-else class="btn btn-blue w-32 mx-2" @click="handleImageCropped">
      Crop Image
    </button>
    <button class="btn btn-gray w-32 mx-2" @click="fileCleared">Clear</button>
    <input
      type="file"
      ref="imageInput"
      accept=".jpg,.jpeg,.png"
      @change="fileChanged"
      :style="{ display: 'none' }"
    />
  </div>
  <div class="my-2 align-baseline text-center">
    <span>Selected File: </span>
    <span v-if="selectedFile">{{ selectedFile.name }}</span>
  </div>
</template>

<script>
import {
  defineComponent,
  ref,
  watchEffect,
  onMounted,
  onUnmounted,
  watch,
} from 'vue';
import Cropper from 'cropperjs';

export default defineComponent({
  name: 'ImageCropper',
  events: ['imageCropped'],
  setup(_props, { emit }) {
    const imageInput = ref(null); // template ref for file input
    const selectedFile = ref(null);
    const imageSrc = ref(null);
    const img = ref(null);
    const fileReader = new FileReader();
    let cropper = null;

    fileReader.onload = (event) => {
      imageSrc.value = event.target.result;
    };

    const handleImageCropped = () => {
      cropper
        .getCroppedCanvas({
          width: 256,
          height: 256,
        })
        .toBlob((blob) => {
          console.log(blob);
          emit('imageCropped', blob);
        }, 'image/jpeg');
    };

    const fileChanged = (e) => {
      const files = e.target.files || e.dataTransfer.files;
      if (files.length) {
        selectedFile.value = files[0];
      }
    };

    const fileCleared = (_) => {
      selectedFile.value = null;
    };

    onMounted(() => {
      cropper = new Cropper(img.value, {
        aspectRatio: 1,
        minCropBoxWidth: 256,
        minCropBoxHeight: 256,
        viewMode: 3,
        dragMode: 'move',
        background: false,
        cropBoxMovable: false,
        cropBoxResizable: false,
      });
    });

    onUnmounted(() => {
      cropper.destroy();
    });

    watchEffect(() => {
      if (selectedFile.value) {
        fileReader.readAsDataURL(selectedFile.value);
      } else {
        imageSrc.value = null;
      }
    });

    watch(
      imageSrc,
      () => {
        if (imageSrc.value) {
          cropper.replace(imageSrc.value);
        }
      },
      {
        flush: 'post', // watch runs after component updates
      }
    );

    return {
      imageInput,
      selectedFile,
      fileChanged,
      fileCleared,
      imageSrc,
      img,
      handleImageCropped,
    };
  },
});
</script>