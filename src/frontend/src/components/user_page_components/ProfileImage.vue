<template>
  <v-layout align-center column style="width: 40%;">
    <h2 style="text-align: center; margin-top: 10px">
      {{this.username}}
      <i class="fa fa-check-square verified-icon ml-5"/>
    </h2>
    <v-img  id="profile-image"
            src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
            alt="Profile picture"
            @click="showProfileImageDialog = true"/>

    <v-btn class="follow-button" v-if="isFollowBtnVisible && $store.state.jws" @click="emitToggleFollowing()">Follow</v-btn>
    <v-btn style="border-color: black; border-style: solid; border-width: 1px;" v-if="isRequestBtnVisible && $store.state.jws" @click="emitToggleFollowing()">Requested</v-btn>
    <v-btn class="unfollow-button" v-if="isUnfollowBtnVisible && $store.state.jws" @click="emitToggleFollowing()">Unfollow</v-btn>

    <transition name="fade" appear>
      <div class="modal-overlay" v-if="showProfileImageDialog" @click="showProfileImageDialog = false"></div>
    </transition>
    <transition name="slide" appear>

      <v-layout class="modal"
                v-if="showProfileImageDialog"
                justify-center
                column>
        <v-btn class="primary mb-2"
               v-if="isMyProfile && $store.state.jws"
               @click="$refs.file.click(); showProfileImageDialog = false">Upload New Profile Photo</v-btn>
        <v-btn  @click="showDialog = false" class="mute-button my-2">
          Show story
        </v-btn>
        <v-btn v-if="isMutedBtnVisible && $store.state.jws" @click="showDialog = false" class="other-buttons my-2">
          Mute
        </v-btn>
        <v-btn v-if="!isMutedBtnVisible && $store.state.jws" @click="showDialog = false" class="mute-button my-2">
          Unmute
        </v-btn>
        <v-btn class="other-buttons my-2"
               v-if="!isMyProfile && $store.state.jws"
               @click="showDialog = false">Report</v-btn>
        <v-btn class="other-buttons my-2"
               v-if="!isMyProfile && $store.state.jws"
               @click="showDialog = false">Block @{{username}}</v-btn>

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
      following: false,
      muted: false,
      isMyProfile: false,
      profile: '',
      waitingForResponse: false
    }
  },
  mounted() {
  },
  props: {
    followingProp: {
      type: Boolean,
      required: true
    },
    username: {
      type: String,
      required: true
    },
    imageSrc: {
      type: String,
      required: false
    },
    isMyProfileProp: {
      type: Boolean,
      required: true
    }
  },
  computed: {
    isFollowBtnVisible() {
      return !this.isMyProfile && !this.following;
    },
    isRequestBtnVisible() {
      return !this.isMyProfile && !this.following && this.waitingForResponse;
    },
    isUnfollowBtnVisible() {
      return !this.isMyProfile && this.following;
    },
    isMutedBtnVisible() {
      return !this.muted && !this.isMyProfile;
    }
  },
  methods: {
    onSelectedFile(event) {
      console.log(event)
      this.profilePicture = event.target.files[0]
      console.log(this.profilePicture)

      this.refreshToken(this.getAHeader())
        .then(rr => {
          this.$store.state.jws = rr.data

          let data = new FormData();
          data.append('profileImg', this.profilePicture);
          let config = {
            headers: {
              'Content-Type': 'multipart/form-data',
              'Authorization': 'Bearer ' + this.$store.state.jws,
            },
          };
          this.axios.post("profilepicture", data, config)
            .then(() => this.isUploadedContent = true)
            .catch(r => console.log(r));
        }).catch(() => this.$router.push('/'));

    },
    emitToggleFollowing() {
      console.log(this.username)
      this.axios.post("users/create/follow", {profile: this.username}, {headers: this.getAHeader()})
        .then(r => {
          console.log(r);
          this.following = !this.following;
          this.$emit('toggle-following', this.following);
        })
        .catch(r => {
          console.log(r)
        })
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

.verified-icon {
  color: #858585;
  transform: scale(1.5);
}

.verified-icon:hover {
  color: #016ddb;
  transform: scale(1.5);
}

.follow-button, .unfollow-button  {

  background-color: transparent;
  color: #016ddb;
  border-color: #016ddb;
  border-style: solid;
  border-width: 1px;
  text-align: -webkit-center;
}

.unfollow-button {
  color: #ff2626;
  border-color: #ff2626;
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