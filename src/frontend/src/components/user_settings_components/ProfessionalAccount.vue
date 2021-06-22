<template>
  <v-form v-model="isFormValid" class="main-div">
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px;">Full Name</h3>
        </div>
        <div style="width: 50%;">
          <v-text-field
          outlined
          value="Name and Lastname"
          v-model="fullname"
          style="width: 400px"
          :rules="[required]"/>
        </div>
      </div>
    </div>
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px; ">Category</h3>
        </div>
        <div style="width: 50%;">
          <v-select outlined :items="roles" v-model="category" style="width: 400px" :rules="[required]"/>
        </div>
      </div>
    </div>
    <div class="item-container " style="height: 400px">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px;">Official document</h3>
          <v-file-input
          style="display: none"
          v-model="document"
          ref="fileInput"
          @change="refreshPreview"
          :rules="[required]"
          />
        </div>
        <div style="width: 50%;">
          <!-- <ImageMessage v-if="showContent" :image-src="this.item.image" @toggle-image-message="showContent = false"/> -->
          <div >
            <!-- <i class="fa fa-image no-content mt-10" v-if="!isUploadedContent"/>
            <v-img class="document-shape" v-else :src="item.image" style="cursor:pointer;" @click="showContent = true"/> -->
            <div v-if="!document" id="add-thumbnail" @click="selectDocument">
              +
            </div>
            <div v-else id="thumbnail">
              <v-img :src="documentUrl" width="128px" height="128px"/>
              <v-btn
              fab
              x-small
              depressed
              absolute
              id="remove-btn" @click="removeDocument">X</v-btn>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="item-container ">
      <v-btn class="primary my-5" @click="sendRequest()" :disabled="!isFormValid">Send request</v-btn>
    </div>
  </v-form>
</template>

<script>
import ImageMessage from "@/components/inbox_components/ImageMessage";
export default {
  name: "ProfessionalAccount",
  components: {ImageMessage},
  data: function () {
    return {
      required: v => !!v || 'Required',
      isFormValid: false,
      documentUrl: '',
      document: null,
      //
      roles : [ 'INFLUENCER', 'SPORTS', 'MEDIA', 'BUSINESS', 'BRAND', 'ORGANIZATION'],
      isUploadedContent: false,
      item: {
        image: ''
      },
      showContent: false,
      fullname: '',
      category: '',
    }
  },
  methods: {

    removeDocument: function() {
      this.document = null;
      this.documentUrl = '';
    },

    refreshPreview: function(file) {
      this.documentUrl = URL.createObjectURL(file);
    },

    selectDocument: function() {
      this.$refs.fileInput.$refs.input.click();
    },

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
    sendRequest() {
      this.refreshToken(this.getAHeader())
          .then(rr => {
            this.$store.state.jws = rr.data;
            let data = new FormData();
            data.append('fullName', this.fullname);
            data.append('category', this.category);
            data.append('document', this.document);
            let config = {
              headers: {
                'Content-Type': 'multipart/form-data',
                'Authorization': 'Bearer ' + this.$store.state.jws,
              },
            };
            this.axios.post("admin/verificationrequest", data, config)
                .then(r =>{
                  console.log(r);
                })
                .catch(r => console.log(r));
          }).catch(r => console.log(r));
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

#add-thumbnail {
  display: grid;
  place-items: center;
  text-align: center;
  font-size: 3rem;
  font-weight: 500;
  background: #eee;
  width: 72px;
  height: 72px;
  cursor: pointer;
}

#thumbnail {
  position: relative;
  display: inline-block;
}

#remove-btn {
  position: absolute;
  top: 0;
  right: 0;
}

</style>