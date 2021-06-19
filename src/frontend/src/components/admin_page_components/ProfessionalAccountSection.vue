<template>
  <div style="height: 100%">
    <div class="title">
      <h2>Professional Account Applications Section</h2>
    </div>
    <div style="display: inline-flex; height: 90%; width: 100%">
      <div class="dropdown-list">
        <div class="inner-layout">
          <div v-for="index in 15" :key="index">
            <ProfessionalAccountApplication :username-prop="index.username"
                          :profile-picture-address-prop="index.profilePictureAddress"
                          @click="selectedApplication(index)"/>
            <v-divider/>
          </div>
        </div>
      </div>
      <div class="application-view">
        <ProfessionalAccountDetails ref="professionalAccountDetails"/>
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

    }
  },
  methods: {
    selectedApplication(applicationData) {
      this.$refs.professionalAccountDetails.$data.applicationData = applicationData
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

.application-view {
  height: 100%;
  width: 60%;
}

</style>