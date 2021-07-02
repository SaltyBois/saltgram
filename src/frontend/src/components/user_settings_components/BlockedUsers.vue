<template>
  <div style="width: 100%; height: 100%; display: flex; flex-direction: column">
    <div style="display: inline-flex; height: 90%">
      <v-layout justify-center
                align-center
                style="height: 100%"
                column>
        <h3>Blocked Users</h3>
        <v-layout column
                  class="scroll-div">
          <BlockedProfile v-for="(item, index) in this.blockedProfiles" 
          :key="index"
          :username-prop="item.username"
          :picture-prop="item.profilePictureURL"
          style="width: 100%" 
          />
        </v-layout>
      </v-layout>
    </div>
    <v-btn class="primary mx-10 mb-3">Confirm changes</v-btn>
  </div>
</template>

<script>
import BlockedProfile from "@/components/user_settings_components/BlockedProfile";


export default {
  name: "BlockedUsers",
  components: { BlockedProfile },

  data: function() {
    return {
      blockedProfiles: [],
    }
  },
  methods: {
    getBlockedProfiles: function() {
      this.axios.get('users/get/blocked', {headers: this.getAHeader()})
      .then(r => {
        this.blockedProfiles = r.data
      })
      .catch(r => console.log(r));
    }
  },
  mounted() {
    getBlockedProfiles();
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