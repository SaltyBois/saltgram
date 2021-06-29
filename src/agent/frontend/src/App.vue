<template>
  <v-app>
    <v-main>
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
