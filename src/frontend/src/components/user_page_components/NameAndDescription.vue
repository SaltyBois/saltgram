<template>
  <v-layout column
            style="height: 60%; margin: 20px; max-height: 500px;">
    <h3><b>{{this.name}}</b></h3>
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
    }
  }
}
</script>

<style scoped>

</style>