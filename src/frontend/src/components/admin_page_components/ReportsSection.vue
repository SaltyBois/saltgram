<template>
  <div style="height: 100%">
    <div class="title">
      <h2>Reports Section</h2>
    </div>
    <div style="display: inline-flex; height: 90%; width: 100%">
      <div class="dropdown-list">
        <div class="inner-layout">
          <div v-for="(item, index) in pendingReports"
               :key="index"
            @click="selectedReport(item)">
            <ReportedUser :username-prop="item.username"
                          :profile-picture-address-prop="item.profilePicture"
                          />
            <v-divider/>
          </div>
        </div>
      </div>
      <div class="report-details">
        <ReportDetails :report-data="repData" v-if="repData" ref="reportDetails"/>
      </div>
    </div>
  </div>
</template>

<script>
import ReportedUser from "@/components/admin_page_components/ReportedUser";
import ReportDetails from "@/components/admin_page_components/ReportDetails";

export default {
  name: "ReportsSection",
  components: { ReportedUser, ReportDetails },
  data: function () {
    return {
      pendingReports: [],
      repData: false
    }
  },
  methods: {
    selectedReport(repData) {
      // console.log(repData)
      this.repData = repData;
     // this.repData.typeMedia = 'image'; // TODO --- THIS INFORMATION ISN'T SENT FROM BACK
      this.$refs.reportDetails.$data.reportData.username = repData.username;
      this.$refs.reportDetails.$data.reportData.reportedMedia = repData.url;
      this.$refs.reportDetails.$data.reportData.profilePictureAddress = repData.profilePicture;
      this.$refs.reportDetails.$data.reportData.id = repData.id;
      this.$refs.reportDetails.$data.reportData.sharedMediaId = repData.sharedMediaId;

    },
    getPendingReports() {
      this.axios.get("admin/inappropriatecontent")
        .then(r => {
          // console.log(r);
          this.pendingReports = r.data;
        } 
        ).catch( err => {
          console.log("Failed to get pending reports.", err);
        })
    },
  },
  mounted() {
    this.getPendingReports();
  },
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

.report-details {
  height: 100%;
  width: 60%;
}

</style>