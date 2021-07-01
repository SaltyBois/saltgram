<template>
  <div class="profile">
    <div style="width: 70px">
      <v-img  class="head"
              src="https://i.pinimg.com/564x/4e/c4/f2/4ec4f2d69c9bc6b152abcb420252c3a8.jpg"
              @click="$router.push('/user/' + usernameProp)"
              alt="Profile picture"/>
    </div>
    <div style="margin: 0 3px; text-align: -webkit-left; width: auto; padding-top: 5px; overflow-x: hidden">
      <h3>@{{this.usernameProp}}</h3>
    </div>
    <div style="margin: 0 3px; text-align: -webkit-center">
      <v-btn v-if="!clicked"
             @click="addCloseFriend()"
             depressed
             class="add-button">
        add
      </v-btn>
      <v-btn v-else
             @click="clicked=!clicked"
             depressed
             class="remove-button">
        remove
      </v-btn>
    </div>
  </div>
</template>

<script>
export default {
  name: "FollowingProfile",
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
  },
  methods: {
    addCloseFriend: function() {
      let dto = {
        profile: this.usernameProp,
      }
      this.axios.post('/users/add/closefrined', dto,  {headers: this.getAHeader()})
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

.add-button, .remove-button  {
  margin: 10px 0;
  width: 100px;
  height: 50px;
  background-color: transparent;
  color: #26a900;
  border-color: #26a900;
  border-style: solid;
  border-width: 1px;
  text-align: -webkit-center;
}

.remove-button {
  color: #ff2626;
  border-color: #ff2626;
}

</style>