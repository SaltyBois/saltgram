<template>
  <div id="bar-main">
    <div id="margin-div">
      <div>
        <v-btn id="main-button"
               @click="$router.push('/main')"
               depressed
               height="100%">
          <v-img id="logo-image"
                 src="https://image.flaticon.com/icons/png/512/114/114928.png"/>
          <h3 id="home-title">Saltgram</h3>
        </v-btn>
      </div>

      <div id="space-div">
        <div class="ml-5 mr-5">
          <v-text-field type="text"
                        id="search-bar"
                        prepend-icon="fa fa-search"
                        v-model="searchQuery"
                        @click:prepend="$router.push('/user/' + searchQuery)"/>
        </div>
      </div>

      <div id="buttons-div">
        <v-btn  id="right-side-button0"
                @click="$router.push('/main')"
                depressed>
          <i class="fa fa-home icons" />
        </v-btn>
        <v-btn  id="right-side-button1"
                @click="$router.push('/newContent')"
                depressed>
          <i class="fa fa-plus-square icons" />
        </v-btn>
        <v-btn  id="right-side-button2"
                @click="$router.push('/notifications'); numberOfNewNotifications = 0"
                depressed>
          <i v-bind:class="numberOfNewNotifications !== 0 ? 'fa fa-heart icons heart' : 'fa fa-heart icons'"/>
          <div v-if="numberOfNewNotifications" class="number-of-notifications"><b>{{numberOfNewNotifications}}</b></div>
        </v-btn>
        <v-btn  id="right-side-button3"
                @click="$router.push('/inbox')"
                depressed>
          <i class="fa fa-commenting icons" />
          <div v-if="numberOfNewChats" class="number-of-chats"><b>{{numberOfNewChats}}</b></div>
        </v-btn>
        <v-btn  id="right-side-button4"
                depressed
                style="text-transform: none"
                @click="profileDropDownMenuActive=!profileDropDownMenuActive">
          <v-img  class="post-header-profile"
                  src="https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg"
                  alt="Profile picture"/>
          <b>@{{this.username}}</b>
        </v-btn>

        <portal to="drop-down-profile-menu">
          <transition name="fade" appear>
            <div class="modal-overlay-2"
                 v-if="profileDropDownMenuActive"
                 @click="profileDropDownMenuActive = false">
            </div>
          </transition>
        </portal>
          <transition name="fade" appear>
            <v-layout class="dropdown-menu"
                      v-if="profileDropDownMenuActive"
                      justify-center
                      align-content-center
                      wrap
                      column>
              <v-btn @click="profileDropDownMenuActive = false; $router.push('/user/' + username)" class="accent">
                <i class="fa fa-address-book mr-1"/>
                profile
              </v-btn>
              <v-btn @click="profileDropDownMenuActive = false;" class="accent mt-3" light>
                <i class="fa fa-folder-open mr-1"/>
                saved
              </v-btn>
              <v-btn @click="profileDropDownMenuActive = false; showProfileSettingsDialog = true; $router.push('/user/settings/' + username)" class="accent mt-3">
                <i class="fa fa-cog mr-1"/>
                settings
              </v-btn>

              <v-divider class="mt-3 mb-3"/>
              <v-btn @click="profileDropDownMenuActive = false; logout()" class="error">
                <i class="fa fa-lock mr-1"/>profile
                logout
              </v-btn>
            </v-layout>
          </transition>
          <transition name="fade" appear>
            <div class="arrow-up"
                 v-if="profileDropDownMenuActive"/>
          </transition>

      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "TopBar",
  data: function() {
    return {
      profileDropDownMenuActive: false,
      showProfileSettingsDialog: false,
      numberOfNewNotifications: 100,
      numberOfNewChats: 20,
      username: '',
      searchQuery: ''
    }
  },
  mounted() {
    if (this.$router.currentRoute.path.includes('/notifications')) this.numberOfNewNotifications = 0;
    if (this.$router.currentRoute.path.includes('/inbox')) this.numberOfNewChats = 0;
    this.loadingJWSOnMounted();
  },
  methods: {
    logout: function() {
      this.$store.state.jws = "";
      this.$router.push('/');
    },
    loadingJWSOnMounted() {
      this.refreshToken(this.getAHeader())
          .then(rr => {
            this.$store.state.jws = rr.data;
            this.axios.get("users", {headers: this.getAHeader()})
                .then(r =>{
                  this.username = r.data.username
                });

          }).catch(() => console.log('No User was founded !?!'));
    }
  }
}
</script>

<style scoped>

#bar-main {
  position: absolute;
  top: 0;
  left: 0;
  display: inline-block;
  flex-direction: row;
  align-content: center;
  /* text-align: center; */
  background: #FFFFFF;
  height: 90px;
  width: 100%;
  outline-width: 3px;
  padding-left: 10%;
  padding-right: 10%;
  border-style: solid;
  border-width: 0;
  border-bottom-color: #484848;
  border-bottom-width: 2px;
}

#home-title {
  font-size: 30px;
  font-family: "Lucida Handwriting", cursive;
  text-transform: capitalize;
}

#logo-image {
  width: 50px;
  height: 50px;
}

#main-button {
  outline-offset: 0px;
  background: transparent;
  align-content: center;
  height: 60px;
}

#margin-div {
  margin: 10px 15px 10px 15px;
  align-content: center;
  flex-direction: row;
  display: flex;
  justify-content: space-between;
  /*background-color: cornflowerblue;*/
}

#space-div {
  width: 100%;
  /*background-color: coral;*/
  justify-content: center;
  align-content: center;
}

#buttons-div {
  right: 0px;
  /*background-color: brown;*/
  display: flex;
  alignment: center;
}

#right-side-button0 {
  background: transparent;
  align-content: center;
  width: 50px;
  height: 100%;
  padding: 5px;
}

#right-side-button1 {
  outline-offset: 0px;
  background: transparent;
  align-content: center;
  outline-color: black;
  width: 50px;
  height: 100%;
  padding: 5px;
}

#right-side-button2 {
  outline-offset: 0px;
  background: transparent;
  align-content: center;
  outline-color: black;
  width: 50px;
  height: 100%;
  padding: 5px;
}

#right-side-button3 {
  outline-offset: 0px;
  background: transparent;
  align-content: center;
  outline-color: black;
  width: 50px;
  height: 100%;
  padding: 5px;
}

#right-side-button4 {
  background: transparent;
  align-content: center;
  outline-color: black;
  width: auto;
  height: 100%;
  padding: 5px;
  overflow: auto;
}

#search-bar {
  width: 200px;
  margin: 10px;
  background: #fafafa;
}

.modal-overlay-2 {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 2;
  background-color: transparent;
  opacity: 0.3;
  height: 100%;
}

.dropdown-menu {
  position: absolute;
  top: 105%;
  right: 10%;
  bottom: 0;
  z-index: 99;
  background-color: white;
  border-style: solid;
  border-width: 3px;
  border-color: #858585;
  width: 15%;
  height: 250px;
  border-radius: 5%;
  align-content: center;
  justify-content: center;
}


.arrow-up {
  position: absolute;
  top: 75%;
  right: 12.5%;
  bottom: 0;
  z-index: 98;
  width: 0;
  height: 0;
  border-left: 40px solid transparent;
  border-right: 40px solid transparent;
  border-bottom: 40px solid #858585;
  border-radius: 10%;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity .25s;
}

.fade-enter,
.fade-leave-to {
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform .5s;
}

.slide-enter,
.slide-leave-to {
  transform: translateY(-50%) translateX(100vw);
}

.icons {
  transform: scale(2.5);
  transition: 0.2s;
  text-align: -webkit-center;
}

/*.icons:hover {*/
/*  transform: scale(2.7);*/
/*  transition: 0.2s;*/
/*}*/

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

.heart {
  color: #ff0051;
  transition: 0.2s;
}

.number-of-notifications {
  position: absolute;
  bottom: 0;
  width: 100%;
  height: 20px;
  font-size: 18px;
  color: black;
  text-align: -webkit-center;
  letter-spacing: 0px;
  transition: 0.2s;
}

.number-of-chats {
  position: absolute;

  bottom: 10px;
  left: 30px;
  width: 29px;
  height: 29px;
  background-color: #14b1ff;
  border-radius: 50%;
  border: solid 2px white ;
  font-size: 18px;
  color: black;
  text-align: -webkit-center;
  letter-spacing: 0px;
  transition: 0.2s;
  padding-right: 1px;
  padding-top: 2px;

}

</style>