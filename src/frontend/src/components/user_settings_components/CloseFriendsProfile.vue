<template>
  <div class="profile">
    <div style="width: 70px">
      <v-img  v-if="this.pictureProp"
              class="head"
              :src="this.pictureProp"
              @click="$router.push('/user/' + this.usernameProp)"
              alt="Profile picture"/>
       <v-img v-else class="head"
          @click="$router.push('/user/' + this.usernameProp)"
          :src="require('@/assets/profile_placeholder.png')"/>
    </div>
    <div style="margin: 0 3px; text-align: -webkit-left; width: auto; padding-top: 5px">
      <h3>@{{this.usernameProp}}</h3>
    </div>
    <div style="margin: 0 3px; text-align: -webkit-center">
      <v-btn v-if="!clicked"
             @click="removeCloseFriend()"
             depressed
             class="remove-button">
        remove
      </v-btn>
      <v-btn v-else
             @click="clicked=!clicked"
             depressed
             class="restore-button">
        restore
      </v-btn>
    </div>
  </div>
</template>

<script>
export default {
  name: "CloseFriendsProfile",
  data: function () {
    return {
      clicked: false,
    }
  },
  props: {
    usernameProp: {
      type: String,
      required: true
    },
    pictureProp: {
      type: String,
      required: true,
    }
  },
  methods: {
    removeCloseFriend: function() {
      let dto = {
        profile: this.usernameProp,
      }
      this.axios.post('/users/remove/closefrined', dto,  {headers: this.getAHeader()})
      .then(r => {
        console.log(r);
        this.$emit('refresh');
        this.clicked = !this.clicked;
      })
      .cathc(r =>{
        console.log(r);
      })
    }
  },

}
</script>

<style scoped>

.profile {
  display: -webkit-inline-flex;
  flex-direction: row;
  height: auto;
  background-color: transparent;

  margin: 5px 5px;

}

.head {
  width: 60px;
  height: 60px;
  margin: 0;
  object-fit: cover;
  border-radius: 20%;
  cursor: pointer;

  border-style: solid;
  border-width: 2px;
  border-color: #323232;
  filter: brightness(1);

  transition: .3s;
  z-index: 0;
}

.head:hover {
  cursor: pointer;
  transition: .3s;
  filter: brightness(0.7);
}

.remove-button, .restore-button  {
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

.restore-button {
  color: #26a900;
  border-color: #26a900;
}

</style>