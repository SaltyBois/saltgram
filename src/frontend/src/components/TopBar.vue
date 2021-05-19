<template>
  <div id="bar-main">
    <div id="margin-div">
      <div>
        <v-btn id="main-button"
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
                        prepend-icon="fa fa-search"/>
        </div>
      </div>

      <div id="buttons-div">
        <v-btn  id="right-side-button0"
                depressed>
          <i class="fa fa-home" />
        </v-btn>
        <v-btn  id="right-side-button1"
                depressed>
          <i class="fa fa-plus-square" />
        </v-btn>
        <v-btn  id="right-side-button2"
                depressed>
          <i class="fa fa-heart" />
        </v-btn>
        <v-btn  id="right-side-button3"
                depressed>
          <i class="fa fa-commenting" />
        </v-btn>
        <v-btn  id="right-side-button4"
                depressed
                @click="profileDropDownMenuActive=!profileDropDownMenuActive">
          <i class="fa fa-address-card mr-3 ml-3" />
          <b>@USERNAME</b>
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
              <v-btn @click="profileDropDownMenuActive = false" class="accent">
                <i class="fa fa-address-book mr-1"/>
                profile
              </v-btn>
              <v-btn @click="profileDropDownMenuActive = false" class="accent mt-3" light>
                <i class="fa fa-folder-open mr-1"/>
                saved
              </v-btn>
              <v-btn @click="profileDropDownMenuActive = false; showProfileSettingsDialog = true; console.log(showProfileSettingsDialog)" class="accent mt-3">
                <i class="fa fa-cog mr-1"/>
                settings
              </v-btn>

              <v-divider class="mt-3 mb-3"/>
              <v-btn @click="profileDropDownMenuActive = false" class="error">
                <i class="fa fa-lock mr-1"/>profile
                logout
              </v-btn>
            </v-layout>
          </transition>
          <transition name="fade" appear>
            <div class="arrow-up"
                 v-if="profileDropDownMenuActive"/>
          </transition>

        <portal to="settings-menu">
          <transition name="fade" appear>
            <div class="modal-overlay-2" v-if="showProfileSettingsDialog" @click="showProfileSettingsDialog = false"></div>
          </transition>
        </portal>
          <transition name="slide" appear>
            <v-layout class="modal2"
                      v-if="showProfileSettingsDialog"
                      justify-center
                      column>
              <v-btn @click="showProfileSettingsDialog = false" class="accent mb-3">Edit profile</v-btn>
              <v-btn @click="showProfileSettingsDialog = false" class="accent mb-3">Change Password</v-btn>
              <v-btn @click="showProfileSettingsDialog = false" class="accent mb-3">Apps and websites</v-btn>
              <v-btn @click="showProfileSettingsDialog = false" class="accent mb-3">Email and sms</v-btn>
              <v-btn @click="showProfileSettingsDialog = false" class="accent mb-3">Push notifications</v-btn>
              <v-btn @click="showProfileSettingsDialog = false" class="accent mb-3">Manage contacts</v-btn>
              <v-btn @click="showProfileSettingsDialog = false" class="accent mb-3">Privacy and security</v-btn>
              <v-btn @click="showProfileSettingsDialog = false" class="accent mb-3">Emails from Instagram</v-btn>
              <v-btn @click="showProfileSettingsDialog = false" class="primary mb-3">Switch to Professional Account</v-btn>

              <v-divider class="mb-3"/>

              <v-btn @click="showProfileSettingsDialog = false" class="error mb-3">Logout</v-btn>
              <v-btn @click="showProfileSettingsDialog = false" class="accent">
                Cancel
              </v-btn>
            </v-layout>
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
  outline-offset: 0px;
  background: transparent;
  align-content: center;
  outline-color: black;
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

.modal2 {
  position: absolute;
  top: 400%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 100;
  border-style: solid;
  border-color: black;
  border-width: 2px;

  width: 100%;
  max-width: 400px;
  background-color: #FFF;
  border-radius: 16px;

  padding: 25px;
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

</style>