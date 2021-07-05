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
        <h2>Start new conversation</h2>

        <v-text-field outlined label="to" v-model="messageReceiver"></v-text-field>

        <v-textarea no-resize label="Message Content" outlined v-model="messageText"></v-textarea>

        <v-btn @click="showDialog = false"
               class="mute-button my-2 "
               property="disabled"
               :disabled="!formValid">
          Send a message
        </v-btn>

        <v-divider class="my-5"/>

        <v-btn @click="showDialog = false" class="mute-button">
          Cancel
        </v-btn>
      </v-layout>
    </transition>
  </div>
</template>

<script>
export default {
  name: "ChatNew",
  data: function () {
    return {
      showDialog: false,
      messageReceiver: '',
      messageText: '',
    }
  },
  computed: {
    formValid: function () {
      return this.messageText !== '' && this.messageReceiver !== ''
    }
  },
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