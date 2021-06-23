<template>
  <v-layout row
            style="justify-content: center; margin-top: 80px">
    <v-layout column
              align-center>
      <h4>Posts</h4>
      <h3><b>{{ postsNumber }}</b></h3>
    </v-layout>
    <v-layout column
              align-center
              class="following-follower-div"
              @click="toggleVisibilityFollwing(); title='Following'">
      <h4>Following</h4>
      <h3><b>{{this.followingProp}}</b></h3>
    </v-layout>
    <v-layout column
              align-center
              class="following-follower-div"
              @click="toggleVisibilityFollwers(); title='Followers'">
      <h4>Followers</h4>
      <h3><b>{{this.followersProp}}</b></h3>
    </v-layout>
    <ModalListOfProfiles ref="modalList" :title="title" :user-prop="userProp"/>
  </v-layout>
</template>

<script>
import ModalListOfProfiles from "@/components/user_page_components/ModalListOfProfiles";

export default {
  name: "ProfileHeader",
  components: {ModalListOfProfiles},
  data: function () {
    return {
      title: ''
    }
  },
  props: {
    userProp: {
      type: String,
      required: true
    },
    followingProp: {
      type: Number,
      required: true
    },
    followersProp: {
      type: Number,
      required: true
    },
    postsNumber: {
      type: Number,
      required: true
    }
  },
  methods: {
    toggleVisibilityFollwing: function() {
      this.$refs.modalList.$data.show = !this.$refs.modalList.$data.show
      this.axios.get("users/following/detailed/" + this.$route.params.username, {headers: this.getAHeader()})
      .then(r => {
        console.log(r.data);
        this.$refs.modalList.$data.profiles = r.data
      })
    },
    toggleVisibilityFollwers: function(){
      this.$refs.modalList.$data.show = !this.$refs.modalList.$data.show
      this.axios.get("users/followers/detailed/" + this.$route.params.username, {headers: this.getAHeader()})
      .then(r => {
        console.log(r.data);
        this.$refs.modalList.$data.profiles = r.data
      })
    },
  }
}
</script>

<style scoped>

.following-follower-div {
  background-color: transparent;
  transition: 0.3s;
  cursor: pointer;

  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;
}

.following-follower-div:hover {
  background-color: #858585;
  transition: 0.3s;
}

</style>