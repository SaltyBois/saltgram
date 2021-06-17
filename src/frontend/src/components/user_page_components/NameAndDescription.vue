<template>
  <v-layout column
            style="height: 60%; margin: 20px; max-height: 500px;">
    <h3><b>{{this.name}}</b></h3>
    <a style="color: #2b80e0" @click="redirectWebSite">{{this.webSite}}</a>
    <h4>
      {{this.description}}
    </h4>
  </v-layout>
</template>

<script>
export default {
  name: "NameAndDescription",
  props: {
    name: {
      type: String,
      required: true
    },
    description: {
      type: String,
      required: true
    },
    webSite: {
      type: String,
      required: true
    },
  },
  mounted() {
  },
  methods: {
    redirectWebSite() {
      if (this.webSite.includes('https://www.') || this.webSite.includes('http://www.')) window.location.href = this.webSite
      else if (this.webSite.includes('www.')) window.location.href = 'https://' + this.webSite
      else window.location.href = 'https://www.' + this.webSite
    }
  }
}
</script>

<style scoped>

</style>