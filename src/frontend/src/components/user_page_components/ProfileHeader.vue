<template>
  <v-layout row
            style="justify-content: center; margin-top: 80px">
    <v-layout column
              align-center>
      <h4>Posts</h4>
      <h3><b>123</b></h3>
    </v-layout>
    <v-layout column
              align-center
              class="following-follower-div"
              @click="toggleVisibility(); title='Following'">
      <h4>Following</h4>
      <h3><b>{{this.followingProp}}</b></h3>
    </v-layout>
    <v-layout column
              align-center
              class="following-follower-div"
              @click="toggleVisibility(); title='Followers'">
      <h4>Followers</h4>
      <h3><b>{{this.followersProp}}</b></h3>
    </v-layout>
    <ModalListOfProfiles ref="modalList" :title="title"/>
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
    followingProp: {
      type: Number,
      required: true
    },
    followersProp: {
      type: Number,
      required: true
    },
  },
  methods: {
    toggleVisibility() {
      this.$refs.modalList.$data.show = !this.$refs.modalList.$data.show
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