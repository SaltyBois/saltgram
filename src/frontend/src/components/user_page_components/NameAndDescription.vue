<template>
  <v-layout column
            style="height: 60%; margin: 20px; max-height: 500px;">
    <h3><b>{{this.name}}</b></h3>
    <a :href="webSite">AAA {{this.webSite}}</a>
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
    console.log('this.name: ' + this.name)
    console.log('this.description: ' + this.description)
    console.log('this.webSite: ' + this.webSite)
  }
}
</script>

<style scoped>

</style>