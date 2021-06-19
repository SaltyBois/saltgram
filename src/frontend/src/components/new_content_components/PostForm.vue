<template>
  <div class="main-div">
    <div class="post-form-body">
      <div class="pt-3 pl-3 post-form-body-left-side">
        <h2>Add posts</h2>
        <div class="thumbnails">
          <div class="thumbnail"
          v-for="img in imageUrls" :key="img">
            <v-img :src="img" height="128px" max-width="128px"/>
            <v-btn
            fab
            absolute
            depressed
            x-small
            @click="removeImg(img)"
            class="remove-btn">X</v-btn>
          </div>
          <v-file-input
          v-model="images"
          hide-input
          multiple
          prepend-icon="fa-plus-circle"
          @change="refreshPreview"
          style="display:none"
          ref="fileInput"
          />  
          <!-- <v-btn @click="uploadFiles">Add</v-btn> -->
          <div class="add-thumbnail" @click="uploadFiles">
            +
          </div>
        </div>
      </div>
      <div class="post-form-body-right-side">
        <h2>Upload info</h2>
        <v-textarea no-resize outlined label="Add a description" dense hide-details style="flex: 0 1 auto"/>
        <v-text-field outlined label="Add location" dense hide-details style="flex: 0 1 auto"/>
        <v-text-field outlined label="Tag people" dense hide-details style="flex: 0 1 auto"/>
        <v-spacer></v-spacer>
        <v-btn color="accent" :disabled="!images.length">Upload</v-btn>
      </div>
    </div>
  </div>
</template>

<script>

export default {
  name: "PostForm",
  data: function () {
    return {
      images: [],
      imageUrls: [],

    }
  },
  methods: {
    uploadFiles: function() {
      this.$refs.fileInput.$refs.input.click();
    },

    removeImg: function(img) {
      let index = this.imageUrls.indexOf(img);
      this.imageUrls.splice(index, 1);
      this.imageUrls = [...this.imageUrls];

      this.images.splice(index, 1);
      this.images = [...this.images];
    },

    refreshPreview: function(files) {
      this.imageUrls = [];
      files = files.slice(0, 10);
      this.images = this.images.slice(0, 10);
      files.forEach(f => {
        this.imageUrls.push(URL.createObjectURL(f));
      });
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
  height: 100%;
  width: 100%;
}

.post-form-body-left-side {
  text-align: -webkit-center;
  width: 60%;
  height: 100%;
}

.post-form-body-right-side {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  width: 40%;
  margin: 10px;
}
.post-form-body-right-side > * {
  padding-top: 10px;
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

.thumbnail {
  position: relative;
  display: inline-block;
  margin: 10px;
}

.add-thumbnail {
  cursor: pointer;
  display: grid;
  place-items: center;
  font-weight: 500;
  font-size: 3rem;
  width: 128px;
  height: 128px;
  background: #eee;
  margin: 10px;
}

.remove-btn {
  top: 0;
  right: 0;
}

.thumbnails {
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  align-items: center;
  flex-wrap: wrap;
}

</style>