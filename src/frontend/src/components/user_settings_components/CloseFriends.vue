<template>
  <div style="width: 100%; height: 100%; display: flex; flex-direction: column">
    <div style="display: inline-flex; height: 90%">
      <v-layout justify-center
                align-center
                style="height: 100%"
                column>
        <h3>Close Friends</h3>
        <v-layout column
                  class="scroll-div">
          <CloseFriendsProfile v-for="(item, index) in this.closeFriends" 
          :key="index"
          :username-prop="item.username"
          :picture-prop="item.profilePictureURL"
          @refresh="refreshData"
          />
        </v-layout>
      </v-layout>
      <v-layout justify-center
                align-center
                style="height: 100%"
                column>
        <h3>List of Following Users</h3>
        <v-layout column
                  class="scroll-div">
          <FollowingProfile v-for="(item, index) in this.following"
          :key="index"
          :username-prop="item.username"
          :picture-prop="item.profilePictureURL"
          @refresh="refreshData"/>
        </v-layout>
      </v-layout>
    </div>
    <v-btn class="primary mx-10 mb-3">Confirm changes</v-btn>
  </div>
</template>

<script>
import CloseFriendsProfile from "@/components/user_settings_components/CloseFriendsProfile";
import FollowingProfile from "@/components/user_settings_components/FollowingProfile";
export default {
  name: "CloseFriends",
  components: { CloseFriendsProfile, FollowingProfile },

  data: function() {
    return {
      closeFriends: [],
      following: [],
    }
  },
  methods: {
    refreshData: function() {
      this.getCloseFriends();
      this.getFollowing();
    },
    getCloseFriends: function() {
      this.axios.get('users/get/closefriend', {headers: this.getAHeader()})
      .then(r => {
        this.closeFriends = r.data
      })
      .catch(r => console.log(r));
    },
    getFollowing: function() {
      this.axios.get('users/get/closefriend/following', {headers: this.getAHeader()})
      .then(r => {
        this.following = r.data
      })
      .catch(r => console.log(r));
    },
  },
  mounted() {
    this.getCloseFriends();
    this.getFollowing();
  }

}
</script>

<style scoped>

.scroll-div {
  overflow-x: hidden;
  overflow-y: auto;

  height: 75%;
  max-height: 90%;
  width: 90%;
  max-width: 600px;
  border: 1px solid black;
}

</style>