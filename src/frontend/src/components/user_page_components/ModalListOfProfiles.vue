<template>
  <div>
    <transition name="fade" appear>
      <div class="modal-overlay" v-if="show" @click="show=!show;"></div>
    </transition>
    <transition name="slide" appear >
      <v-layout class="modal"
                justify-center
                align-center
                v-if="show"
                column>
        <h3>{{title}}</h3>
        <v-layout column
                  class="scroll-div">
          <ProfileInList v-for="(item, index) in this.profiles" :key="index" 
          :username-prop="item.username" 
          :following-prop="item.following" 
          :pending-prop="item.pending" 
          :picture-prop="item.profilePictureURL"
          :user-prop="userProp"/>
        </v-layout>
        <v-divider class="mt-5 mb-5"/>
        <v-btn @click="show=!show;" class="accent">
          Exit
        </v-btn>
      </v-layout>
    </transition>
  </div>
</template>

<script>
import ProfileInList from "@/components/user_page_components/ProfileInList";

export default {
  name: "ModalListOfProfiles",
  components: {ProfileInList},
  props: {
    title: {
      required: true,
      type: String
    },
    userProp: {
      type: String,
      required: true
    }
  },
  data: function () {
    return {
      show: false,
      profiles: [],
    }
  },
  methods: {

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

  width: auto;
  max-width: 600px;
  max-height: 800px;
  background-color: #FFF;
  border-radius: 16px;

  padding: 25px;
}

.scroll-div {
  overflow-x: hidden;
  overflow-y: auto;

  max-height: 400px;
  width: auto;
  max-width: 600px;
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

</style>