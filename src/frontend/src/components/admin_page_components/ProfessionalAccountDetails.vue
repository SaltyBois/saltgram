<template>
  <div style="margin: 5px; border: 1px solid black; border-radius: 10px; padding: 3px; height: 97%; overflow-y: auto; overflow-x: hidden">
    <div style="height: 97%">
      <v-img  class="chat-head"
              :src="applicationData.profilePictureAddress"
              @click="$router.push('/user/' + applicationData.username)"
              alt="Profile picture"/>
      <h3 class="my-1">Username: {{applicationData.username}}</h3>
      <h4 class="my-1">Name and LastName: {{applicationData.fullname}}</h4>
      <h4 class="my-1">Professional Account Type: {{applicationData.accountType}}</h4>
      <ImageMessage v-if="showContent" :image-src="this.applicationData.documentMedia" @toggle-image-message="showContent = false"/>
      <v-img  class="content-item my-2"
              style="border: 1px solid black; border-radius: 10px"
              v-if="this.applicationData.documentMedia"
              :src="this.applicationData.documentMedia"
              alt="Profile picture"
              @click="showContent = true"/>
      <div style="display: inline-flex; text-align: -webkit-center; width: 100%;">
        <v-btn class="accept-button mx-2" width="45%" style="font-size: 12px; letter-spacing: 0" @click="acceptApplication">Accept Application</v-btn>
        <v-btn class="reject-button mx-2" width="45%" style="font-size: 12px; letter-spacing: 0" @click="rejectApplication">Reject Application</v-btn>
      </div>
    </div>
  </div>
</template>

<script>
import ImageMessage from "@/components/inbox_components/ImageMessage";
export default {
  name: "ProfessionalAccountDetails",
  components: { ImageMessage },
  data: function () {
    return {
      applicationData: {
        username: '{{USERNAME}}',
        profilePictureAddress: 'https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg',
        documentMedia: 'https://www.ozonpress.net/wp-content/uploads/2016/09/licna-karata.jpg',
        fullname: 'Imen Prezimenovic',
        accountType: 'Business',
        userId: '',
        requestId: '',
      },
      showContent: false,
    }
  },
  methods: {
    acceptApplication() {
      let acceptrequest = {id: this.applicationData.requestId, status: 'ACCEPTED'}
      this.axios.put("admin/verificationrequest", acceptrequest)
        .then(r => {
          console.log(r);
          this.$router.go(0);
        } 
        ).catch( err => {
          console.log("Failed to accept requests.", err);
        })
    },
    rejectApplication() {
      let rejectrequest = {id: this.applicationData.requestId, status: 'REJECTED'}
      this.axios.put("admin/verificationrequest", rejectrequest)
        .then(r => {
          console.log(r);
          this.$router.go(0);
        } 
        ).catch( err => {
          console.log("Failed to reject requests.", err);
        })
    }
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

.accept-button, .reject-button  {
  margin: 3px 0;
  width: 100px;
  height: 50px;
  background-color: transparent;
  color: #13d700;
  border-color: #13d700;
  border-style: solid;
  border-width: 1px;
  text-align: -webkit-center;
}

.reject-button {
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