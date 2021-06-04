<template>
  <div class="inner-inbox-content-footer">
    <i class="fa fa-photo picture-button" @click="$refs.file.click()"/>
    <input type="file"
           ref="file"
           style="display: none"
           @change="onSelectedFile"
           accept="image/*">
    <v-text-field label="Enter a message" style="padding-left: 5px; margin-bottom: 5px; min-width: 100%" />

    <div style="float: right; height: available; display: inline-block;">
      <v-btn class="follow-button" style="margin: 10px 5px 5px;width: 75px">
        Send
      </v-btn>
    </div>
  </div>
</template>

<script>

export default {
  name: "ChatInput",
  components: {
  },
  data: function () {
    return {
      search: '',
    }
  },
  methods: {
    onSelectedFile(event) {
      console.log(event)
      this.profilePicture = event.target.files[0]
      console.log(this.profilePicture)
    },
  }
}
</script>

<style scoped>

  .inner-inbox-content-footer {
    height: 60px;
    display: flex;
    width: 75%;
    flex-direction: row;
    justify-content: space-between;
    padding-left: 10px;
  }

  .picture-button {
    transform: scale(1.5);
    margin-top: 24px;
    margin-right: 10px;
    cursor: pointer;
    transition: 0.3s;
  }

  .picture-button:hover {
    transform: scale(1.6);
    transition: 0.3s;
    color: #016ddb;
  }

  .follow-button {
    margin: 10px 0;
    width: 100px;
    height: 50px;
    background-color: transparent;
    color: #016ddb;
    border-color: #016ddb;
    border-style: solid;
    border-width: 1px;
    text-align: -webkit-center;
  }

</style>