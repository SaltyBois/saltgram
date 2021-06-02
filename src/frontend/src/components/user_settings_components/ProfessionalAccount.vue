<template>
  <div class="main-div">

      <input type="file"
             ref="file"
             style="display: none"
             @change="onSelectedFile($event)"
             accept="image/*">
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px;">Name</h3>
        </div>
        <div style="width: 50%;">
          <v-text-field outlined value="Name and Lastname" style="width: 400px"/>
        </div>
      </div>
    </div>
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px; ">Category</h3>
        </div>
        <div style="width: 50%;">
          <v-select outlined :items="roles"  style="width: 400px"/>
        </div>
      </div>
    </div>
    <div class="item-container " style="height: 400px">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px;">Official document</h3>
          <v-btn class="primary"
                 @click="$refs.file.click(); showProfileImageDialog = false">Upload Content</v-btn>
          <v-btn class="error mt-3"
                 v-if="isUploadedContent"
                 @click="removeContent(item)">Remove content</v-btn>
        </div>
        <div style="width: 50%;">
          <ImageMessage v-if="showContent" :image-src="this.item.image" @toggle-image-message="showContent = false"/>
          <div class="image-shape">
            <i class="fa fa-image no-content mt-10" v-if="!isUploadedContent"/>
            <v-img class="document-shape" v-else :src="item.image" style="cursor:pointer;" @click="showContent = true"/>
          </div>
        </div>
      </div>
    </div>
    <div class="item-container ">
      <v-btn class="primary my-5">Send request</v-btn>
    </div>
  </div>
</template>

<script>
import ImageMessage from "@/components/inbox_components/ImageMessage";
export default {
  name: "ProfessionalAccount",
  components: {ImageMessage},
  data: function () {
    return {
      roles : [ 'Influencer', 'Sports', 'New/Media', 'Business', 'Brand', 'Organization'],
      isUploadedContent: false,
      item: {
        image: ''
      },
      showContent: false
    }
  },
  methods: {
    removeContent(item) {
      item = {};
      console.log(item)
      this.isUploadedContent = false;
    },
    onSelectedFile(event) {
      var files = event.target.files || event.dataTransfer.files;
      if (!files.length)
        return;
      console.log(files.length)
      this.item.image = URL.createObjectURL(files[0])
      this.isUploadedContent = true;
    },
  }
}
</script>

<style scoped>

.main-div {
  display: inline-flex;
  flex-direction: column;
  height: 100%;
  overflow-y: auto;
}

.item-container {
  height: 100px;
  display: inline-flex;
  flex-direction: row;
  justify-content: center;
  align-content: center;
  text-align: -webkit-center;
}

.no-content {
  position: relative;
  top: 10%;
  left: 0;
  transform: scale(2.5);
}

.image-shape {
  border: black 1px solid ;
  min-height: 50%;
  height: auto;

  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;
}

.document-shape {
  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;
}

</style>