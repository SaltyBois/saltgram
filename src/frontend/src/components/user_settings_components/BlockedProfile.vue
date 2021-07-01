<template>
  <div style="width: 100%">
    <div class="profile">
      <div style="width: 70px">
        <v-img  class="head"
                src="https://i.pinimg.com/564x/4e/c4/f2/4ec4f2d69c9bc6b152abcb420252c3a8.jpg"
                @click="$router.push('/user')"
                alt="Profile picture"/>
      </div>
      <div style="margin: 0 3px; text-align: -webkit-left; width: 50%; padding-top: 5px">
        <h3>@{{this.usernameProp}}</h3>
      </div>
      <div style="margin: 0 3px; text-align: -webkit-center">
        <v-btn v-if="!blocked"
               @click="blockProfile()"
               depressed
               class="remove-button">
          Block
        </v-btn>
        <v-btn v-else
               @click="unblockProfile()"
               depressed
               class="restore-button">
          Unblock
        </v-btn>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "BlockedProfile",
  data: function () {
    return {
      blocked: true,
    }
  },
  props: {
    usernameProp: {
      type: String,
      required: true
    },
  },
  methods: {
    unblockProfile: function() {
      let dto = {
        profile: this.usernameProp,
      }
      this.axios.post('/users/unblock/profile', dto,  {headers: this.getAHeader()})
      .then(r => {
        console.log(r);
        this.blocked = false;
        //this.$emit('get-blocked');
      })
      .cathc(r =>{
        console.log(r);
      })
    },
    blockProfile: function() {
      let dto = {
        profile: this.usernameProp,
      }
      this.axios.post('/users/block/profile', dto,  {headers: this.getAHeader()})
      .then(r => {
        console.log(r);
        this.blocked = true;
        //this.$emit('get-blocked');
      })
      .cathc(r =>{
        console.log(r);
      })
    },
  }
}
</script>

<style scoped>

.profile {
  display: flex;
  flex-direction: row;
  width: 100%;
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
  color: #0eb4ff;
  border-color: #0eb4ff;
}

</style>