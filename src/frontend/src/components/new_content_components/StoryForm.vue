<template>
  <div class="main-div">
    <div class="post-form-body">
      <div class="pt-3 pl-3 post-form-body-left-side">
        <h2>Add stories</h2>
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
          <div class="add-thumbnail" @click="selectFiles">
            +
          </div>
        </div>
      </div>
      <div class="post-form-body-right-side">
        <h2>Upload info</h2>
        <v-textarea
        v-model="description"
        no-resize
        outlined
        label="Add a description"
        dense
        hide-details
        style="flex: 0 1 auto"/>
        <!--<v-text-field
        
        outlined
        label="TODO Add location"
        dense
        hide-details
        style="flex: 0 1 auto"/>-->
        <geosearch @selected="selectLocation"></geosearch>
        <div v-if="selectedLocation">
          <v-text-field 
          v-model="location.state"
          label="Country"
          disabled/>
          <v-text-field 
          label="City"
          v-model="location.city"
          disabled/>
          <v-text-field 
          v-model="location.street"
          label="Street"
          disabled/>
        </div>
        <v-combobox
        v-model="tags"
        chips
        clearable
        label="Tags"
        multiple
        append-icon=""
        solo>
          <template v-slot:selection="{ attrs, item }">
              <v-chip
              v-bind="attrs"
              close
              @click:close="removeTag(item)">
                {{item}}
              </v-chip>
          </template>
        </v-combobox>
        <v-spacer></v-spacer>
        <v-checkbox
        v-model="closeFriends"
        label="Close friends"/>
        <v-btn color="accent" :disabled="!images.length" @click="uploadFiles" :loading="uploading">Upload</v-btn>
      </div>
    </div>
  </div>
</template>

<script>
import ImageMessage from "@/components/inbox_components/ImageMessage";
import Media from "@dongido/vue-viaudio"

export default {
  name: "StoryForm",
  components: {ImageMessage, Media},
  data: function () {
    return {
      closeFriends: false,
      uploading: false,
      description: "",
      location: {
				country: "",
				state:   "",
				zipCode: "",
        city: "",
				street:  "",
        name: "",
      },
      tags: [],
      images: [],
      imageUrls: [],

      selectedLocation: false,
    }
  },
  methods: {

     selectLocation: function(l) {
      this.selectedLocation = true;
      this.location.country = l.address.country_code;
      this.location.state = l.address.country;
      this.location.zipCode = l.address.postcode;
      this.location.city = l.address.city;
      this.location.street = l.address.road;
      this.location.name = l.display_name.split(',')[0];
    },

    uploadFiles: function() {
      this.uploading = true;
      let data = new FormData();
      this.images.forEach(img => {
        data.append('stories', img);
      });
      this.tags.forEach(tag => {
        data.append('tags', tag);
      });
      data.append('description', this.description)
      data.append('location', JSON.stringify(this.location))
      data.append('closeFriends', this.closeFriends)
      this.refreshToken(this.getAHeader())
        .then(rr => {
          this.$store.state.jws = rr.data;
          let config = {
            headers: {
              'Content-Type': 'multipart/form-data',
              'Authorization': 'Bearer ' + this.$store.state.jws,
            },
          };
          this.axios.post('content/story', data, config)
            .then(() => {
              this.uploading = false
              this.$router.push('/main')
            })
            .catch(r => console.log(r));
        }).catch(() => this.$router.push('/'));
    },

    selectFiles: function() {
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

    select: function(r) {
      this.selected = r;
      this.results = [];
      this.$emit('selected', r);
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