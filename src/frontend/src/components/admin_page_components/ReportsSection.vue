<template>
  <div style="height: 100%">
    <div class="title">
      <h2>Reports Section</h2>
    </div>
    <div style="display: inline-flex; height: 90%; width: 100%">
      <div class="dropdown-list">
        <div class="inner-layout">
          <div v-for="index in 15" :key="index">
            <ReportedUser :username-prop="index.username"
                          :profile-picture-address-prop="index.profilePictureAddress"
                          @click="selectedReport(index)"/>
            <v-divider/>
          </div>
        </div>
      </div>
      <div class="report-details">
        <ReportDetails ref="reportDetails"/>
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

    }
  },
  methods: {
    selectedReport(reportData) {
      this.$refs.reportDetails.$data.reportData = reportData
    }
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

.report-details {
  height: 100%;
  width: 60%;
}

</style>