<template>
  <div style="width: 100%; height: 100%; display: flex; flex-direction: column">
    <div style="display: inline-flex; height: 90%">
      <v-layout justify-center
                align-center
                style="height: 100%"
                column>
        <h3>Muted Users</h3>
        <v-layout column
                  class="scroll-div">
          <MutedProfile v-for="(item, index) in this.mutedProfiles"
           :key="index"
           :username-prop="item.username"
           @get-muted="getMutedProfiles"
           style="width: 100%" />
        </v-layout>
      </v-layout>
    </div>
    <v-btn class="primary mx-10 mb-3">Confirm changes</v-btn>
  </div>
</template>

<script>
import MutedProfile from "@/components/user_settings_components/MutedProfile";
export default {
  name: "MutedUsers",
  components: { MutedProfile },

  data: function() {
    return {
      mutedProfiles: [],
    }
  },
  methods: {
    getMutedProfiles: function() {
      this.axios.get('users/get/muted', {headers: this.getAHeader()})
      .then(r => {
        this.mutedProfiles = r.data
      })
      .catch(r => console.log(r));
    }
  },
  mounted() {
    this.getMutedProfiles();
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