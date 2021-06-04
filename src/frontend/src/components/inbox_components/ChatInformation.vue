<template>
  <div>
    <transition name="fade" appear>
      <div class="modal-overlay" v-if="showDialog" @click="showDialog = false"></div>
    </transition>
    <transition name="slide" appear>

      <v-layout class="modal"
                v-if="showDialog"
                justify-center
                column>
        <v-btn v-if="muted" @click="showDialog = false" class="other-buttons my-2">
          Mute
        </v-btn>
        <v-btn v-else @click="showDialog = false" class="mute-button my-2">
          Unmute
        </v-btn>
        <v-btn class="other-buttons my-2"
               @click="showDialog = false">Delete chat</v-btn>
        <v-btn class="other-buttons my-2"
               @click="showDialog = false">Block @{{username}}</v-btn>
        <v-btn class="other-buttons my-2"
               @click="showDialog = false">Report @{{username}}</v-btn>

        <v-divider class="mt-5 mb-5"/>
        <v-btn @click="showDialog = false" class="mute-button">
          Cancel
        </v-btn>
      </v-layout>
    </transition>
  </div>
</template>

<script>

export default {
  name: "ChatInformation",
  data: function () {
    return {
      showDialog: false,
      muted: false
    }
  },
  props: {
    username: {
      type: String,
      required: true
    }
  }
}
</script>

<style scoped>

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 98;
  background-color: rgba(0, 0, 0, 0.3);
}

.modal {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 99;

  width: 100%;
  max-width: 400px;
  background-color: #FFF;
  border-radius: 16px;

  padding: 25px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity .5s;
}

.fade-enter,
.fade-leave-to {
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform .5s;
}

.slide-enter,
.slide-leave-to {
  transform: translateY(-50%) translateX(100vw);
}

.mute-button, .other-buttons  {
  background-color: transparent;
  color: #016ddb;
  border-color: #016ddb;
  border-style: solid;
  border-width: 1px;
  text-align: -webkit-center;
}

.other-buttons {
  color: #ff2626;
  border-color: #ff2626;
}

</style>