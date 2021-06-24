<template>
  <div class="profile">
    <div style="width: 70px">
      <v-img  class="head"
              :src="pictureProp"
              @click="$router.push('/user/' + usernameProp)"
              alt="Profile picture"/>
    </div>
    <div style="margin: 0 3px; text-align: -webkit-left; width: 300px; padding-top: 5px">
      <h3>@{{usernameProp}}</h3>
    </div>
    <div style="margin: 0 3px; text-align: -webkit-center" v-if="usernameProp !== userProp">
      <v-btn v-if="!followingProp && pendingProp"
             depressed
             class="follow-button">
        pending
      </v-btn>
      <v-btn v-else-if="!followingProp"
            @click="follow()"
             depressed
             class="unfollow-button">
        follow
      </v-btn>
      <v-btn v-else
             @click="unfollow()"
             depressed
             class="unfollow-button">
        unfollow
      </v-btn>
    </div>
  </div>
</template>

<script>
export default {
  name: "ProfileInList",
  data: function () {
    return {
      clicked: false,
    }
  },
  props: {
    userProp: {
      type: String,
      required: true
    },
    usernameProp: {
      type: String,
      required: true
    },
    followingProp: {
      type: Boolean,
      required: true
    },
    pendingProp: {
      type: Boolean,
      required: true
    },
    pictureProp: {
      type: String,
      required: true
    }
  },
  methods: {
    follow: function() {
      let dto = {
        profile: this.usernameProp
      };
      this.axios.post("users/create/follow", dto, {headers: this.getAHeader()})
      .then(r => {
        if(r.data == "PENDING") {
          this.pendingProp = true;
        } else if(r.data == "Following") {
          this.followingProp = true;
        }
      })
    },
    unfollow: function() {
      this.axios.post("users/unfollow", {profile: this.usernameProp}, {headers: this.getAHeader()})
        .then(r => {
          console.log(r)
          this.followingProp = false;
          this.pendingProp = false;
        })
        .catch(r => {
          console.log(r)
        })
    }
  },
  mounted() {
    console.log(this.userProp)
    console.log(this.usernameProp)
  }
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

.follow-button, .unfollow-button  {
  margin: 10px 0;
  width: 100px;
  height: 50px;
  background-color: transparent;
  color: #016ddb;
  border-color: #016ddb;
  border-style: solid;
  border-width: 1px;
  text-align: -webkit-center;
}

.unfollow-button {
  color: #ff2626;
  border-color: #ff2626;
}

</style>