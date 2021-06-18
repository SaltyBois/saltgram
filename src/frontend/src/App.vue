<template>
  <v-app>
    <v-main>
      <router-view :key="$route.fullPath"/>
    </v-main>
  </v-app>
</template>

<script>


export default {
  name: 'App',
  data: () => ({
  }),
  watch: {
    $route: {
      immediate: true,
      handler(to) {
        document.title = to.meta.title || 'Saltgram';
      }
    },
  },
  created() {
    window.addEventListener('beforeunload', () => {
      localStorage.setItem('vuexstore', JSON.stringify(this.$store.state));
    });
    localStorage.getItem('vuexstore' && this.$store.replaceState(Object.assign(this.$store.state, JSON.parse(localStorage.getItem('vuexstore')))));
  },
};
</script>

<style scoped>

</style>

