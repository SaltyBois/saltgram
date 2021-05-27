<template>
  <v-layout align-center column style="width: 40%;">
    <h2 style="text-align: center; margin-top: 10px">USERNAME</h2>
    <v-img  id="profile-image"
            src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
            alt="Profile picture"
            @click="showProfileImageDialog = true"/>

    <transition name="fade" appear>
      <div class="modal-overlay" v-if="showProfileImageDialog" @click="showProfileImageDialog = false"></div>
    </transition>
    <transition name="slide" appear>

      <v-layout class="modal"
                v-if="showProfileImageDialog"
                justify-center
                column>
        <v-btn class="primary"
               @click="$refs.file.click(); showProfileImageDialog = false">Upload New Profile Photo</v-btn>

        <v-divider class="mt-5 mb-5"/>
        <v-btn @click="showProfileImageDialog = false" class="accent">
          Cancel
        </v-btn>
      </v-layout>
    </transition>

    <input type="file"
           ref="file"
           style="display: none"
           @change="onSelectedFile"
           accept="image/*">
  </v-layout>
</template>

<script>
export default {
  name: "ProfileImage",
  data: function () {
    return {
      showProfileImageDialog: false,
    }
  },
  methods: {
    onSelectedFile(event) {
      console.log(event)
      this.profilePicture = event.target.files[0]
      console.log(this.profilePicture)
    },
  }
}
</script>

<style scoped>

#profile-image {
  width: 300px;
  height: 300px;
  object-fit: cover;
  border-radius: 20%;
  margin: 10px;
  cursor: pointer;

  border-style: solid;
  border-width: 10px;
  border-color: cornflowerblue;
  transition: .1s;
  z-index: 0;
}

#profile-image:hover {
  transition: .1s;
  border-width: 5px;
  border-color: cornflowerblue;
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

</style>