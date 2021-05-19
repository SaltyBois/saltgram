<template>
  <v-app>
    <portal-target name="drop-down-profile-menu" />
    <portal-target name="settings-menu"/>
    <v-main>
      <top-bar style="position: sticky; z-index: 2"/>
      <router-view />
    </v-main>
  </v-app>
</template>

<script>

import TopBar from "./components/TopBar.vue";

export default {
  name: 'App',

  components: {
    'top-bar': TopBar
  },

  data: () => ({
    //
  }),
};
</script>

<style scoped>

</style>

