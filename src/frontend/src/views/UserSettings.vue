<template>
  <div style="background-color: #efeeee; height: 100%">
    <TopBar style="position: sticky; z-index: 2"/>
    <div id="main-div">
      <div class="menu-div">
        <v-btn class="primary my-2"
               @click="option = 0"
               v-bind:class="option === 0 ? 'primary' : 'accent'">Edit profile</v-btn>
        <v-btn class="primary my-2"
               @click="option = 1"
               v-bind:class="option === 1 ? 'primary' : 'accent'">Change password</v-btn>
        <v-btn class="primary my-2"
               @click="option = 2"
               v-bind:class="option === 2 ? 'primary' : 'accent'">Switch to professional Account</v-btn>
        <v-btn class="my-2"
               @click="option = 3"
               v-bind:color="option === 3 ? 'success' : 'close-friends-btn'">Close friends</v-btn>
        <v-btn class="my-2"
               @click="option = 4"
               v-bind:color="option === 4 ? 'error-btn' : 'error'">Muted users</v-btn>
        <v-btn class="my-2"
               @click="option = 5"
               v-bind:color="option === 5 ? 'error-btn' : 'error'">Blocked users</v-btn>

      </div>
      <div class="components-div">

        <EditProfile v-if="option === 0"/>

        <ChangePassword v-if="option === 1"/>

        <ProfessionalAccount v-if="option === 2"/>

        <CloseFriends v-if="option === 3"/>

        <MutedUsers v-if="option === 4" />

        <BlockedUsers v-if="option === 5"/>

      </div>
    </div>
  </div>
</template>

<script>
import TopBar from "@/components/TopBar";
import EditProfile from "@/components/user_settings_components/EditProfile";
import ChangePassword from "@/components/user_settings_components/ChangePassword";
import ProfessionalAccount from "@/components/user_settings_components/ProfessionalAccount";
import CloseFriends from "@/components/user_settings_components/CloseFriends";
import MutedUsers from "@/components/user_settings_components/MutedUsers";
import BlockedUsers from "@/components/user_settings_components/BlockedUsers";

export default {
  name: "UserSettings",
  components: {TopBar, EditProfile, ChangePassword, ProfessionalAccount, CloseFriends, MutedUsers, BlockedUsers},
  data: function() {
    return {
      option: 0,
    }
}

}
</script>

<style scoped>

#main-div {
  background-color: white;
  display: flex;
  margin: 10px 10%;
  height: 80vh;
  flex-direction: row;

  border: black 1px solid ;

  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;
}

.menu-div {
  display: flex;
  padding: 5px;
  height: 100%;
  width: 33%;
  flex-direction: column;
  overflow-y: auto;
  overflow-x: hidden;
  border-right: black 1px solid ;

}

.components-div {
  display: flex;
  padding: 5px;
  height: 100%;
  width: 67%;
  flex-direction: column;
  overflow-y: auto;
  overflow-x: hidden;

}

.close-friends-btn {
  color: white;
  background-color: #74cf5a !important;
}

.error-btn {
  color: white;
  background-color: #ff1d1d !important;
}

</style>