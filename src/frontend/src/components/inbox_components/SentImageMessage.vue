<template>
  <div class="inbox-message-div-right">
    <div class="inbox-message-sent">
      <v-img src="https://metro.co.uk/wp-content/uploads/2020/08/IMG_8718-ec27.jpg?quality=90&strip=all" class="image-class" @click="toggle" alt="image"/>
    </div>
    <div class="message-time-sent">
      15.4.2015. 17:46
    </div>
    <ImageMessage v-if="imageVisible"
                  image-src="https://metro.co.uk/wp-content/uploads/2020/08/IMG_8718-ec27.jpg?quality=90&strip=all"
                  @toggle-image-message="toggle"/>
  </div>
</template>

<script>
import ImageMessage from "@/components/inbox_components/ImageMessage";


export default {
  name: "SentImageMessage",
  components: {ImageMessage},
  data: function () {
    return {
      imageVisible: false,
    }
  },
  methods: {
    toggle() {
      this.imageVisible = !this.imageVisible
    },
  }
}
</script>

<style scoped>

  .inbox-message-div-right {
    display: inline-block;
    height: auto;
    width: 100%;
    align-content: center;
    margin: 3px 0;
  }

  .inbox-message-sent {
    background-color: #016ddb;
    color: white;
    border: #323232 2px solid;
    float: right;
    max-width: 60%;
    max-height: 410px;
    overflow-wrap: anywhere;
    padding: 5px;
    margin-right: 5px;

    border-start-end-radius: 10px 10px;
    border-end-end-radius: 10px 10px;
    border-start-start-radius: 10px 10px;
    border-end-start-radius: 10px 10px;
  }

  .message-time-sent {
    float: right;
    text-align: right;
    font: 12px bold;
    padding-top: 10px;
    padding-right: 5px;
  }

  .image-class {
    max-height: 100%;
    max-width: 100%;
    cursor: pointer
  }

</style>