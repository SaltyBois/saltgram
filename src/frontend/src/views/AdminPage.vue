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
        </div>
        <div style="height: 10%;" class="sub-menu-div">
          <v-btn class="error my-2"
                 @click="logout()">Logout</v-btn>
        </div>

      </div>
      <div class="components-div">


      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "AdminPage",
  components: {},
  data: function () {
    return {
      option: 0,
    }
  },
  methods: {
    logout: function() {
      this.$store.state.jws = "";
      this.$router.push('/');
    },
  }
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
  overflow-y: auto;
  overflow-x: hidden;
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

</style>