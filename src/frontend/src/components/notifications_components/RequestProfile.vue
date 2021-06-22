<template>
  <div class="layout-div">
    <div class="post-header-left-side">
      <v-img  class="post-header-profile"
              src="https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg"
              @click="$router.push('/user')"
              alt="Profile picture"/>
      <b @click="$router.push('/user')" style="cursor: pointer">{{this.usernameProp}}</b>
      <p style="margin-left: 5px; margin-right: 5px; margin-top: 15px">has sent You a follow request</p>
      <v-btn class="add-button mx-2" @click="acceptRequest(), $emit('reloadRequests')">
        Accept
      </v-btn>
      <v-btn class="ignore-button mx-2" @click="ignoreRequest()">
        Ignore
      </v-btn>
      <p style="color: #858585; margin-top: 15px">1h ago</p>
    </div>
  </div>
</template>

<script>
export default {
  name: "RequestProfile",
  props: {
    usernameProp: {
      type: String,
      required: true
    }
  },
  methods: {
    acceptRequest: function(){
      let dto = {
        profile: this.usernameProp,
        accepted: true,
      };
      this.axios.post("users/follow/request/", dto,  {headers: this.getAHeader()})
      .then(r =>{
        console.log(r.data);
        this.$emit('reload-requests')
      })
    },
    ignoreRequest: function() {
      let dto = {
        profile: this.usernameProp,
        accepted: false,
      };
      this.axios.post("users/follow/request/", dto,  {headers: this.getAHeader()})
      .then(r =>{
        console.log(r.data);
        this.$emit('reload-requests');
      })

    }
  }
}



</script>

<style scoped>

.layout-div {
  height: 70px;
  width: 100%;

  display: flex;
  flex-direction: row;
  flex-flow: wrap;
  transition: 0.3s;
  cursor: pointer;
}

.layout-div:hover {
  background-color: #e0e0e0;
  transition: 0.3s;
}

.post-header-left-side {
  direction: ltr;
  flex-direction: row;
  text-align: -webkit-center;
  align-items: center;
  float: left;
  display: flex;
  justify-content: center
}

.post-header-profile {
  width: 30px;
  height: 30px;
  object-fit: cover;
  border-radius: 20%;
  margin: 10px;
  cursor: pointer;


  filter: brightness(1);

  transition: .3s;
  z-index: 0;
}

.add-button, .ignore-button  {
  margin: 10px 0;
  width: 100px;
  height: 50px;
  background-color: transparent;
  color: #ff2626;
  border-color: #ff2626;
  border-style: solid;
  border-width: 1px;
  text-align: -webkit-center;
}

.add-button {
  color: #26a900;
  border-color: #26a900;
}

</style>