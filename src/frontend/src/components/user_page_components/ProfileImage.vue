<template>
  <v-layout align-center column style="width: 40%;">
    <h2 style="text-align: center; margin-top: 10px">
      {{this.username}}
      <i v-if="verified" class="fa fa-check-square verified-icon ml-5"/>
    </h2>
    <v-img  class="profile-image"
            v-if="imageSrc"
            :src="imageSrc"
            alt="Profile picture"
            @click="showProfileImageDialog = true"/>
    <v-img  class="profile-image"
            v-else
            :src="require('@/assets/profile_placeholder.png')"
            alt="Profile picture"
            @click="showProfileImageDialog = true"/>
    <StoryView ref="storyView" v-if="userStories" :stories="userStories"/>
    <v-img v-else class="head"
      @click="showContent = true"
      :src="require('@/assets/profile_placeholder.png')"/>

    <v-btn class="follow-button" v-if="isFollowBtnVisible && $store.state.jws" @click="emitToggleFollowing()">Follow</v-btn>
    <v-btn style="border-color: black; border-style: solid; border-width: 1px;" v-if="isRequestBtnVisible && $store.state.jws">Pending</v-btn>
    <v-btn class="unfollow-button" v-if="isUnfollowBtnVisible && $store.state.jws" @click="emitToggleUnfollow()">Unfollow</v-btn>

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
        <v-btn v-if="userStories.length != 0" @click="showProfileImageDialog = false; toggle()" class="mute-button my-2">
          Show story
        </v-btn>
        <v-btn v-if="canMute && isMutedBtnVisible && $store.state.jws" @click="muteProfile()" class="other-buttons my-2">
          Mute
        </v-btn>
        <v-btn v-if="canMute && !isMutedBtnVisible && $store.state.jws" @click="unmuteProfile()" class="mute-button my-2">
          Unmute
        </v-btn>
        <v-btn class="other-buttons my-2"
               v-if="!isMyProfile && $store.state.jws"
               @click="showProfileImageDialog = false">Report</v-btn>
        <v-btn class="other-buttons my-2"
               v-if="isBlocked && $store.state.jws"
               @click="blockProfile()">Block @{{username}}</v-btn>
        <v-btn class="other-buttons my-2"
               v-if="!isBlocked && $store.state.jws"
               @click="unblockProfile()">Unblock @{{username}}</v-btn>

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
import StoryView from "@/components/StoryView";

export default {
  name: "ProfileImage",
  components: {StoryView},
  data: function () {
    return {
      showProfileImageDialog: false,
      following: false,
      muted: false,
      isMyProfile: false,
      profile: '',
      waitingForResponse: false,
      userStories: [],
      storyVisible: false,
      blocked: false,
    }
  },
  mounted() {
    console.log(this.userStories);
  },
  props: {
    username: {
      type: String,
      required: true
    },
    imageSrc: {
      type: String,
      required: true,
    },
    isMyProfileProp: {
      type: Boolean,
      required: true
    },
    verified: {
      type: Boolean,
      required: true
    },
  },
  computed: {
    isFollowBtnVisible() {
      return !this.isMyProfileProp && !this.following && !this.waitingForResponse;
    },
    isRequestBtnVisible() {
      return !this.isMyProfileProp && !this.following && this.waitingForResponse;
    },
    isUnfollowBtnVisible() {
      return !this.isMyProfileProp && this.following;
    },
    isMutedBtnVisible() {
      return !this.muted && !this.isMyProfileProp;
    },
    isBlocked() {
      return !this.isMyProfileProp && !this.blocked;
    },
    canMute() {
      if (this.isMyProfileProp) {
        return false;
      } else if( !this.isMyProfileProp && this.following) {
        return true;
      }
      else {
        return false;
      }
    }
  },
  methods: {
    toggle() {
      //this.$refs.storyView.$data.stories = this.userStories;
      this.$refs.storyView.toggleView();

      //console.log(this.$refs.storyView.$data.stories);
      //console.log(this.userStories);
    },
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
          if(r.data == "PENDING") {
            this.following = false;
            this.waitingForResponse = true;
          } else if(r.data == "Following") {
            this.following = true;
          }
         this.$emit('toggle-following', this.following);
         this.$emit('following-changed');
        })
        .catch(r => {
          console.log(r)
        })
    },
    emitToggleUnfollow() {
      this.axios.post("users/unfollow", {profile: this.username}, {headers: this.getAHeader()})
        .then(r => {
          console.log(r)
          this.following = false;
          this.waitingForResponse = false;
          this.$emit('toggle-following', this.following);
          this.$emit('following-changed');
        })
        .catch(r => {
          console.log(r)
        })
    },
    checkIfMuted: function() {
      if(!this.isMyProfileProp){
        this.axios.get("users/check/muted/" + this.username,{headers: this.getAHeader()})
        .then(r => {
          this.muted = r.data;
        })
        .catch(r => {
          console.log(r);
        })
      }
    },
    muteProfile: function() {
      let dto = {
        profile: this.username,
      }
      this.axios.post('/users/mute/profile',dto,  {headers: this.getAHeader()})
      .then(r => {
        console.log(r);
        this.showProfileImageDialog = false;
        this.muted = true;
      })
      .cathc(r =>{
        console.log(r);
      })
    },
    unmuteProfile: function() {
      let dto = {
        profile: this.username,
      }
      this.axios.post('/users/unmute/profile',dto,  {headers: this.getAHeader()})
      .then(r => {
        console.log(r);
        this.showProfileImageDialog = false;
        this.muted = false;
      })
      .cathc(r =>{
        console.log(r);
      })
    }, 
    blockProfile: function() {
      let dto = {
        profile: this.username,
      }
      this.axios.post('/users/block/profile',dto,  {headers: this.getAHeader()})
      .then(r => {
        console.log(r);
        this.showProfileImageDialog = false;
        this.blocked = true;
        this.$router.go(0);
      })
      .cathc(r =>{
        console.log(r);
      })
    },
    unblockProfile: function() {
      let dto = {
        profile: this.username,
      }
      this.axios.post('/users/unblock/profile',dto,  {headers: this.getAHeader()})
      .then(r => {
        console.log(r);
        this.showProfileImageDialog = false;
        this.blocked = false;
        this.$router.go(0);
      })
      .cathc(r =>{
        console.log(r);
      })
    }
  },
  
}
</script>

<style scoped>

.profile-image {
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