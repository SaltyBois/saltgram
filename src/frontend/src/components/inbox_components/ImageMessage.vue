<template>
  <transition name="slide" appear>
    <v-layout class="modal-overlay" @click="toggle">
      <v-layout class="modal"
                justify-center
                align-center
                column>
        <v-img class="image" :src="imageSrc" />
      </v-layout>
    </v-layout>
  </transition>
</template>

<script>
export default {
  name: "ImageMessage",
  data: function() {
    return {

    }
  },
  props: {
    imageSrc: {
      type: String,
      required: true
    }
  },
  methods: {
    toggle() {
      this.$emit('toggle-image-message')
    }
  }
}
</script>

<style scoped>

.modal {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 99;

  max-width: 80%;
  max-height: 90vh;
  background-color: #FFF;
  border-radius: 16px;

  padding: 25px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 98;
  background-color: rgba(0, 0, 0, 0.3);
}

.image {
  max-width: available;
  max-height: available;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform .5s;
}

.slide-enter,
.slide-leave-to {
  transform: translateY(-50%) translateX(100vw);
}

</style>