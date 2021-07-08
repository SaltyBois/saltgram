<template>
  <v-app>
    <v-main>
      <hexpage-loader />
      <router-view/>
    </v-main>
  </v-app>
</template>

<script>

export default {
  name: 'App',
  created() {
    window.addEventListener('beforeunload', () => {
      localStorage.setItem('agentvuexstore', JSON.stringify(this.$store.state));
    });
    localStorage.getItem('agentvuexstore' && this.$store.replaceState(Object.assign(this.$store.state, JSON.parse(localStorage.getItem('agentvuexstore')))));
  },
};
</script>

<style>
.v-input {
  font-family: monospace !important;
}

.v-input__icon--append .v-icon {
  color: var(--v-secondary-lighten4) !important;
}
.theme--dark.v-label.v-label--active {
  color: var(--v-info-base) !important;
}
</style>
