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

<style>

::-webkit-scrollbar {
  width: 13px;
}
::-webkit-scrollbar-track {
  background: #ffffff;
  border-radius: 10px;
  border: 1px solid #aeaeb5
}
::-webkit-scrollbar-thumb {
  background: #dadae3;
  border-radius: 10px;
  border: 1px solid #aeaeb5
}
::-webkit-scrollbar-thumb:hover {
  background: #c4c4cc
}
::-webkit-scrollbar-thumb:active {
  background: #aeaeb5
}

::-webkit-scrollbar-corner {
  border-radius: 10px;
}

</style>

