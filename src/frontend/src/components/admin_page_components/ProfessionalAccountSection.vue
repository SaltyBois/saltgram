<template>
  <div style="height: 100%">
    <div class="title">
      <h2>Professional Account Applications Section</h2>
    </div>
    <div style="display: inline-flex; height: 90%; width: 100%">
      <div class="dropdown-list">
        <div class="inner-layout">
          <div v-for="(item, index) in pendingVerificationRequests" :key="index"
            @click="selectedApplication(item)">
            <ProfessionalAccountApplication :username-prop="item.username"
                          :profile-picture-address-prop="item.profilePicture"
                          />
            <v-divider/>
          </div>
        </div>
      </div>
      <div class="application-view">
        <ProfessionalAccountDetails v-if="appdata"
                                    :application-data="appdata"/>
      </div>
    </div>
  </div>
</template>

<script>
import ProfessionalAccountApplication from "@/components/admin_page_components/ProfessionalAccountApplication";
import ProfessionalAccountDetails from "@/components/admin_page_components/ProfessionalAccountDetails";

export default {
  name: "ProfessionalAccountSection",
  components: { ProfessionalAccountApplication, ProfessionalAccountDetails },
  data: function () {
    return {
      pendingVerificationRequests: [],
      appdata: false,
    }
  },
  methods: {
    /*applicationData: {
        username: '{{USERNAME}}',
        profilePictureAddress: 'https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg',
        documentMedia: 'https://www.ozonpress.net/wp-content/uploads/2016/09/licna-karata.jpg',
        fullName: 'Imen Prezimenovic',
        accountType: 'Business'
      },*/
    selectedApplication(appData) {
      this.appdata = appData;
      // console.log(appData);
      // this.$refs.professionalAccountDetails.$data.applicationData.profilePictureAddress = appData.profilePicture;
      // this.$refs.professionalAccountDetails.$data.applicationData.documentMedia = appData.url;
      // this.$refs.professionalAccountDetails.$data.applicationData.fullname = appData.fullname;
      // this.$refs.professionalAccountDetails.$data.applicationData.accountType = appData.category;
      // this.$refs.professionalAccountDetails.$data.applicationData.userId = appData.userId;
      // this.$refs.professionalAccountDetails.$data.applicationData.requestId = appData.id;
    },
    getPendingVerificationRequests() {
      this.axios.get("admin/verificationrequest")
        .then(r => {
          console.log(r);
          this.pendingVerificationRequests = r.data;
        } 
        ).catch( err => {
          console.log("Failed to get pending requests.", err);
        })
    },
  },
  mounted() {
    this.getPendingVerificationRequests();
  }
}
</script>

<style scoped>

.title {
  height: 10%;
  text-align: center;
}

.dropdown-list {
  height: 100%;
  width: 40%;
  overflow-y: hidden;
}

.inner-layout {
  margin: 5px;
  height: 97%;

  border: 1px solid black;
  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;

  overflow-y: auto;
}

.application-view {
  height: 100%;
  width: 60%;
}

</style>