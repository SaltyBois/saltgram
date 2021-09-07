<template>
  <div style="background-color: #efeeee; height: 100%;   padding-top: 1%; text-align: -webkit-center">
    <div style="text-align: -webkit-center; display: inline-flex;">
      <v-img id="logo-image"
             src="https://image.flaticon.com/icons/png/512/114/114928.png"/>
      <h3 id="home-title" class="mr-3">Saltgram</h3>
      <h1>Admin Page</h1>
    </div>
    <div style="text-align: center">

    </div>
    <div id="main-div">
      <div class="menu-div">
        <div style="height: 90%;" class="sub-menu-div">
          <v-btn class="primary my-2"
                 @click="option = 0"
                 v-bind:class="option === 0 ? 'primary' : 'accent'">Reports</v-btn>
          <v-btn class="primary my-2"
                 @click="option = 1"
                 v-bind:class="option === 1 ? 'primary' : 'accent'">Profesional Account Aplications</v-btn>
          <v-btn class="primary my-2"
                 @click="option = 2"
                 v-bind:class="option === 2 ? 'primary' : 'accent'">Agent requests</v-btn>
        </div>
        <v-divider/>
        <div style="height: 10%;" class="sub-menu-div">
          <v-btn class="error my-2"
                 @click="logout()">Logout</v-btn>
        </div>

      </div>
      <div class="components-div">

        <ReportsSection v-if="option === 0"/>

        <ProfessionalAccountSection v-if="option === 1" />

        <div v-if="option == 2">
          <div id="agentreqs">
            <div v-for="r in agentRequests" :key="r" class="agentreq">
              <b>{{r}}</b>
              <v-spacer></v-spacer>
              <v-btn @click="acceptAgent(r)">Accept</v-btn>
            </div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script>
import ReportsSection from "@/components/admin_page_components/ReportsSection";
import ProfessionalAccountSection from "@/components/admin_page_components/ProfessionalAccountSection";

export default {
  name: "AdminPage",
  components: { ReportsSection, ProfessionalAccountSection },
  data: function () {
    return {
      option: 0,
      agentRequests: [],
    }
  },
  methods: {
    acceptAgent: function(email) {
      this.refreshToken(this.getAHeader())
      .then(rr => {
        this.$store.state.jws = rr.data;
        this.axios.post('admin/agent', email, {headers: this.getAHeader()})
          .then(() => this.getAgentRequests());
      })
    },

    getAgentRequests: function() {
      this.refreshToken(this.getAHeader())
      .then(rr => {
        this.$store.state.jws = rr.data;
        this.axios.get('admin/agent', {headers: this.getAHeader()})
          .then(r => this.agentRequests = r.data);
      })
    },
    logout: function() {
      this.$store.state.jws = "";
      this.$router.push('/');
    },
  },
  mounted() {
    this.getAgentRequests();
  },
}
</script>

<style scoped>

#main-div {
  background-color: white;
  display: flex;
  margin: 0 10%;
  height: 80vh;
  flex-direction: row;

  border: black 1px solid ;

  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;
}

.menu-div {
  display: flex;
  padding: 5px;
  height: 100%;
  width: 33%;
  flex-direction: column;
  overflow-y: auto;
  overflow-x: hidden;
  border-right: black 1px solid ;
}

.sub-menu-div {
  display: flex;
  flex-direction: column;
}

.components-div {
  display: flex;
  padding: 5px;
  height: 100%;
  width: 67%;
  flex-direction: column;
}

#logo-image {
  width: 50px;
  height: 50px;
}

#home-title {
  font-size: 30px;
  font-family: "Lucida Handwriting", cursive;
  text-transform: capitalize;
}

#agentreqs {
  display: flex;
  padding: 10px;
  flex-direction: column;
}

.agentreq {
  display: flex;
  flex-direction: row;
  padding: 10px;
  align-items: center;
}

</style>