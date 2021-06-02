<template>
  <div class="main-div">
    <div class="post-form-header">
      <h2>Post Form</h2>
    </div>
    <div class="post-form-body">
      <div class="pt-3 pl-3 post-form-body-left-side">
        <div>
          <v-btn class="primary"
                 @click="$refs.file.click(); showProfileImageDialog = false">Upload Content</v-btn>
          <v-btn class="error ml-3"
                 v-if="isUploadedContent"
                 @click="removeContent(item)">Remove content</v-btn>
        </div>

        <input type="file"
               ref="file"
               style="display: none"
               @change="onSelectedFile($event)"
               accept="image/*,video/*">

        <div class="content-shape">
          <ImageMessage v-if="showContent && typeContent === 'image'" :image-src="this.item.image" @toggle-image-message="showContent = false"/>
          <v-img  class="content-item"
                  v-if="isUploadedContent && typeContent === 'image'"
                  :src="this.item.image"
                  alt="Profile picture"
                  @click="showContent = true"/>
          <i class="fa fa-image no-content mt-10" v-if="!isUploadedContent"/>
          <Media class="content-item"
                 v-if="isUploadedContent && typeContent === 'video'"
                 :kind="'video'"
                 :autoplay="true"
                 :controls="true"
                 :loop="true"
                 :style="{width: '500px'}"
                 :src="[this.item.image]"/>
        </div>
      </div>
      <div class="post-form-body-right-side">
        <div>
        </div>
        <div style="height: 100%; padding: 10px 5px">
          <v-textarea no-resize outlined label="Add a description" style="width: 100%; min-height: auto; padding: 5px;"/>
          <v-text-field outlined label="Add location" style="width: 100%; min-height: auto; padding: 5px;"/>
          <v-text-field outlined label="Tag people" style="width: 100%; min-height: auto; padding: 5px;"/>
          <v-btn class="primary" :disabled="!isUploadedContent">Post</v-btn>
        </div>
      </div>
    </div>
<!--    <div class="post-form-footer">-->

<!--    </div>-->
  </div>
</template>

<script>
import ImageMessage from "@/components/inbox_components/ImageMessage";
import Media from "@dongido/vue-viaudio"

export default {
  name: "PostForm",
  components: {ImageMessage, Media},
  data: function () {
    return {
      isUploadedContent: false,
      item: {
        image: ''
      },
      showProfileImageDialog: false,
      showContent: false,
      typeContent: ''
    }
  },
  methods: {
    onSelectedFile(event) {
      var files = event.target.files || event.dataTransfer.files;
      if (!files.length)
        return;
      console.log(files.length)
      console.log(files[0])
      this.item.image = URL.createObjectURL(files[0])
      console.log(this.item.image)
      if (files[0]['type'].includes('image')) this.typeContent = 'image';
      else this.typeContent = 'video';
      console.log(this.typeContent)
      this.isUploadedContent = true;
    },
    removeContent(item) {
      item = {};
      console.log(item)
      this.isUploadedContent = false;
    },
  }
}
</script>

<style scoped>

.main-div {
  width: 100%;
  height: 100%;
}

.post-form-header {
  height: 10%;
  width: 100%;
  text-align: center;
}

.post-form-body {
  display: flex;
  height: 90%;
  width: 100%;
}

.post-form-body-left-side {
  text-align: -webkit-center;
  width: 60%;
  height: 100%;
}

.post-form-body-right-side {
  text-align: -webkit-center;
  width: 40%;
  height: 100%;
}

.content-shape {
  /*display: flex;*/
  width: min-content;
  height: min-content;
  min-width: 100px;
  min-height: 100px;
  max-width: 60vh;
  max-height: 60vh;

  object-fit: contain;

  text-align: -webkit-center;
  justify-content: center;

  margin-top: 10px;
  border: 1px black solid;

  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;
}

.content-item {
  width: 100%;
  height: 100%;

  max-width: 60vh;
  max-height: 60vh;
  display: block;
  text-align: -webkit-center;
  justify-content: center;

  object-fit: scale-down;


  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;

  cursor: pointer;
}

.no-content {
  position: relative;
  top: 10%;
  left: 0;
  transform: scale(2.5);
}

</style>