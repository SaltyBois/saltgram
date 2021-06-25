<template>
  <div style="margin: 5px; border: 1px solid black; border-radius: 10px; padding: 3px; height: 97%; overflow-y: auto; overflow-x: hidden">
    <div style="height: 97%">
      <v-img  class="chat-head"
              v-if="reportData.profilePicture"
              :src="reportData.profilePicture"
              @click="$router.push('/user/' + reportData.username)"
              alt="Profile picture"/>
      <v-img  class="chat-head"
              v-else
              :src="require('@/assets/profile_placeholder.png')"
              @click="$router.push('/user/' + reportData.username)"
              alt="Profile picture"/>
      <h3>{{reportData.username}}</h3>
      <ImageMessage v-if="showContent && this.reportData.typeMedia === 'image'" :image-src="this.reportData.url" @toggle-image-message="showContent = false"/>
      <v-img  class="content-item my-2"
              style="border: 1px solid black; border-radius: 10px"
              v-if="this.reportData.url && this.reportData.typeMedia === 'image'"
              :src="this.reportData.url"
              alt="Profile picture"
              @click="showContent = true"/>
      <video class="content-item my-2"
             v-if="this.reportData.url && this.reportData.typeMedia === 'video'"
             :kind="'video'"
             :autoplay="true"
             :controls="true"
             :loop="true"
             :style="{width: '300px'}"
             :src="[this.reportData.url]"/>
      <v-textarea no-resize height="150px" outlined style="width: 90%" readonly v-model="reportData.description"/>
      <div style="display: inline-flex; text-align: -webkit-center; width: 100%;">
        <v-btn class="sanction-button mx-2" width="30%" style="font-size: 12px; letter-spacing: 0" @click="removeContent()">Remove content</v-btn>
        <v-btn class="sanction-button" width="30%" style="font-size: 12px; letter-spacing: 0" @click="removeUser()">Remove user</v-btn>
        <v-btn class="reject-button mx-2" width="30%" style="font-size: 12px; letter-spacing: 0" @click="rejectReport()">Reject report</v-btn>
      </div>
    </div>

  </div>
</template>

<script>
import ImageMessage from "@/components/inbox_components/ImageMessage";

export default {
  name: "ReportDetails",
  components: { ImageMessage, },
  props: {
    reportData: { type: Object, required: true }
  },
  data: function () {
    return {
      // reportData: {
      //   username: '',
      //   profilePictureAddress: '',
      //   reportedMedia: '',
      //   typeMedia: 'image',
      //   description: ''
      // },
      showContent: false,
    }
  },
  methods: {
    removeContent() {
      let remove = {id: this.reportData.id, sharedMediaId: this.reportData.sharedMediaId}
      this.axios.put("admin/removeinappropriatecontent", remove)
        .then(r => {
          console.log(r);
          this.$router.go(0);
        } 
        ).catch( err => {
          console.log("Failed to remove content", err);
        })
    },
    removeUser() {
      let remove = {id: this.reportData.id, sharedMediaId: this.reportData.sharedMediaId}
      this.axios.put("admin/removeprofile", remove)
        .then(r => {
          console.log(r);
          this.$router.go(0);
        } 
        ).catch( err => {
          console.log("Failed to remove profile", err);
        })
    },
    rejectReport() {
      let reject = {id: this.reportData.id}
      this.axios.put("admin/rejectinappropriatecontent", reject)
        .then(r => {
          console.log(r);
          this.$router.go(0);
        } 
        ).catch( err => {
          console.log("Failed to reject report.", err);
        })
    },
  },
  mounted() {
    console.log(this.reportData)
  }
}
</script>

<style scoped>

.chat-head {
  width: 60px;
  height: 60px;
  object-fit: cover;
  border-radius: 20%;
  cursor: pointer;

  border-style: solid;
  border-width: 2px;
  border-color: #323232;

  transition: .3s;
  z-index: 0;
}

.reject-button, .sanction-button  {
  width: 100px;
  height: 50px;
  background-color: transparent;
  color: #016ddb;
  border-color: #016ddb;
  border-style: solid;
  border-width: 1px;
  text-align: -webkit-center;
}

.sanction-button {
  color: #ff2626;
  border-color: #ff2626;
}

.content-item {
  width: 100%;
  height: 100%;

  max-width: 25vh;
  max-height: 25vh;
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

</style>