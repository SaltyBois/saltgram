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
          <CloseFriendsProfile v-for="index in 14" :key="index"/>
        </v-layout>
<!--        <v-divider class="mt-5 mb-5"/>-->
      </v-layout>
      <v-layout justify-center
                align-center
                style="height: 100%"
                column>
        <h3>List of Following Users</h3>
        <v-layout column
                  class="scroll-div">
          <FollowingProfile v-for="index in 14" :key="index"/>
        </v-layout>
<!--        <v-divider class="mt-5 mb-5"/>-->
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
  components: { CloseFriendsProfile, FollowingProfile }
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