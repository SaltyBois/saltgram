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
  }
};
</script>

<style scoped>

</style>

