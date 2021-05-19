<template>
  <div id="user-main">
    <div id="user-header">
      <div id="user-icon-logout">
        <v-layout align-center column style="width: 40%;">
          <h2 style="text-align: center; margin-top: 10px">USERNAME</h2>
            <v-img  id="profile-image"
                    src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
                    alt="Profile picture"
                    @click="showProfileImageDialog = true"/>
            <transition name="fade" appear>
              <div class="modal-overlay" v-if="showProfileImageDialog" @click="showProfileImageDialog = false"></div>
            </transition>
            <transition name="slide" appear>

              <v-layout class="modal"
                        v-if="showProfileImageDialog"
                        justify-center
                        column>
                  <v-btn class="primary"
                         @click="$refs.file.click(); showProfileImageDialog = false">Upload New Profile Photo</v-btn>

                <v-divider class="mt-5 mb-5"/>
                <v-btn @click="showProfileImageDialog = false" class="accent">
                  Cancel
                </v-btn>
              </v-layout>
            </transition>

          <input type="file"
                 ref="file"
                 style="display: none"
                 @change="onSelectedFile"
                 accept="image/*">
        </v-layout>

        <v-layout column
                  style="width: 70%"
                  justify-space-between>
          <v-layout style="height: 40%;"
                    justify-center
                    column>

            <v-layout row
                      style="justify-content: center; margin-top: 80px">
              <v-layout column
                        align-center>
                <h4>Posts</h4>
                <h3><b>123</b></h3>
              </v-layout>
              <v-layout column
                        align-center>
                <h4>Following</h4>
                <h3><b>1000</b></h3>
              </v-layout>
              <v-layout column
                        align-center>
                <h4>Followers</h4>
                <h3><b>10k</b></h3>
              </v-layout>
            </v-layout>

          </v-layout>

          <v-layout column
                    style="height: 60%; margin: 20px">
            <h3><b>Ime i prezime</b></h3>
            <h4 id="description">
              Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et
              dolore magna aliqua. Aliquam ut porttitor leo a diam. Porttitor eget dolor morbi non arcu risus quis.
              Gravida cum sociis natoque penatibus et. At in tellus integer feugiat scelerisque.
              Tellus orci ac auctor augue mauris. Mi bibendum neque egestas congue quisque egestas.
              Scelerisque eleifend donec pretium vulputate sapien nec sagittis. At varius vel pharetra vel turpis nunc.
            </h4>
          </v-layout>
        </v-layout>


          <div v-if="showSettings" id="settings-container">
              <v-form id="settings-password" v-model="isFormValid">
                  <h2>Change password:</h2>
                  <b id="err">{{err}}</b>
                  <v-text-field
                  v-model="oldPassword"
                  label="Old password"
                  :rules="[rules.required, rules.min, different]"
                  :append-icon="showPassword2 ? 'fa-eye' : 'fa-eye-slash'"
                  :type="showPassword2 ? 'text' : 'password'"
                  @click:append="showPassword1 = !showPassword1"
                  required/>
                  <v-text-field
                  v-model="newPassword1"
                  label="New password"
                  hint="Min 8 characters, upper/lowercase, number and symbol"
                  :rules="[rules.required, rules.min, rules.passMatch, passMatch, passStr]"
                  :append-icon="showPassword2 ? 'fa-eye' : 'fa-eye-slash'"
                  :type="showPassword2 ? 'text' : 'password'"
                  @click:append="showPassword2 = !showPassword2"
                  required/>

                  <v-text-field
                  v-model="newPassword2"
                  label="Confirm new password"
                  hint="Min 8 characters, upper/lowercase, number and symbol"
                  :rules="[rules.required, rules.min, rules.passMatch, passMatch, passStr]"
                  :append-icon="showPassword3 ? 'fa-eye' : 'fa-eye-slash'"
                  :type="showPassword3 ? 'text' : 'password'"
                  @click:append="showPassword3 = !showPassword3"
                  required/>
                  <b id="pass-str"><div>Password strength: </div><div>{{passScoreText}}</div></b>
                  <password-meter :password="newPassword1" @score="onScore"/>
                  <v-btn :disabled="!isFormValid" @click="changePassword">Change password</v-btn>
              </v-form>
          </div>
      </div>
      <div id="user-info">
          <h2 id="username">{{user.username}}</h2>
          <p>{{user.fullName}}</p>
          <p>{{user.email}}</p>
      </div>
    </div>

<!--        TODO: STORY HIGHLIGHTS-->
    <v-layout id="user-stories"
              style="background: darkolivegreen"
              row>
      <v-layout class="inner-story-layout"
                style="background-color: darkorchid; margin: 10px">
        <v-container class="story-highlight-layout"
                     style="background-color: red"
                     column
                     fluid>
          <v-img  class="story-highlight"
                  src="https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg"
                  alt="Profile picture"/>
          <h5>Highlights 1</h5>
        </v-container>
        <v-container class="story-highlight-layout"
                     style="background-color: blue;"
                     column>
          <v-img  class="story-highlight"
                  src="https://filmdaily.co/wp-content/uploads/2020/05/coughing-cat-meme-lede.jpg"
                  alt="Profile picture"/>
          <h5>Highlights 2</h5>
        </v-container>
        <v-container class="story-highlight-layout"
                     style="background-color: coral;"
                     column>
          <v-img  class="story-highlight"
                  src="https://www.arabianbusiness.com/public/styles/square/public/images/2021/03/28/meme.jpg?itok=DeJVUtab"
                  alt="Profile picture"/>
          <h5>Highlights 3</h5>
        </v-container>
      </v-layout>






    </v-layout >
<!--        TODO: POSTS -->
    <v-layout id="user-media"
              style="background: lightskyblue"
              column>
        <v-layout class="inner-post-layout" row style="background-color: red">

          <v-layout class="post-layout" style="background-color: blue">
            <v-img  class="post"
                    src="https://i.kym-cdn.com/entries/icons/original/000/032/100/cover4.jpg"
                    alt="Profile picture"/>
          </v-layout>

          <v-layout class="post-layout" style="background-color: gold">

          </v-layout>

          <v-layout class="post-layout" style="background-color: forestgreen">

          </v-layout>

          <v-layout class="post-layout" style="background-color: deeppink">

          </v-layout>
        </v-layout>
    </v-layout>
<!--        <div style="background: red; height: 300px">-->

<!--        </div>-->
<!--        <div style="background: blue; height: 300px">-->

<!--        </div>-->
  </div>
</template>

<script>
import passwordMeter from 'vue-simple-password-meter';
import TopBar from "@/components/TopBar";
export default {
    components: {
      TopBar,
        passwordMeter,
    },
    data: function() {
        return {
            err: "",
            showPassword1: false, 
            showPassword2: false, 
            showPassword3: false, 
            oldPassword: '',
            newPassword1: '',
            newPassword2: '',
            passScore: 0,
            passScoreText: '',
            isFormValid: false,
            showSettings: false,
            profilePicture: '',
            descriptionTextBoxReadOnly: true,
            descriptionTextBoxText: '[[Description]]',
            showProfileImageDialog: false,
            showProfileSettingsDialog: false,
            user: {},
            rules: {
                required: v => !!v || "Required",
                min: v => v.length >= 8 || "Min 8 characters",
            },
        }
    },

    methods: {
        changePassword: function() {
            this.err = "";
            this.refreshToken()
                .then((r) => {
                    this.$store.state.jws = r.data;
                    let changeRequest = {
                        oldPassword: this.oldPassword,
                        newPassword: this.newPassword1,
                    }
                    this.axios.post("users/changepass", changeRequest, {headers: {"Authorization": "Bearer " + this.$store.state.jws}})
                        .then(r => {
                            console.log(r);
                            this.showSettings = false;
                        })
                        .catch(r => {
                            console.log(r);
                            this.err = "Invalid old password!";
                        });
                })
                .catch(r => {
                    console.log(r);
                    this.$router.push("/")
                });
        },

        onSelectedFile(event) {
          console.log(event)
            this.profilePicture = event.target.files[0]
            console.log(this.profilePicture)
        },

        changeDescription() {
          this.descriptionTextBoxReadOnly = !this.descriptionTextBoxReadOnly

        },

        onScore: function({score, strength}) {
            // console.log("Password score: " + strength);
            this.passScore = score;
            this.passScoreText = strength;
        },

        logout: function() {
            this.$store.state.jws = "";
            this.$router.go();
        },

        refreshToken: async function() {
            let jws = this.$store.state.jws
            if (!jws) {
                this.$router.push("/")
            }

            return this.axios.get("auth/refresh", {headers: {"Authorization": "Bearer " + jws}})
        },

        getUserInfo: function() {
            let jws = this.$store.state.jws
            if (!jws) {
                this.$router.push("/");
            }

            this.axios.get("users", {headers:{"Authorization": "Bearer " + jws}})
                .then(r => {
                    console.log(r);
                    this.user = r.data;
                })
                .catch(r => {
                    console.log(r);
                    // NOTE(Jovan): Try to refresh
                    this.this.refreshToken()
                        .then(r => {
                            console.log(r);
                            this.$store.state.jws = r.data;
                            this.$router.go()
                        })
                        .catch(r => {
                            console.log(r);
                            this.$router.push("/");
                        });
                });
        },
    },
    mounted() {
        // this.getUserInfo(); // TODO UNCOMMENT THIS
    },
    computed: {
        different: function() {
            return (this.newPassword1 != this.oldPassword && this.newPassword2 != this.oldPassword) || "Cannot use old password"
        },

        passStr: function() {
            console.log("pass score: " + this.passScore);
            return this.passScore > 3 || "Use a stronger password!";
        },

        passMatch: function() {
            return this.newPassword1 == this.newPassword2 || "Passwords must match"
        },
    },
}
</script>

<style scoped>
    #user-main {
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        align-content: center;
        /* text-align: center; */
        background: #efeeee;
        height: auto;
    }

    #user-header {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        /*background: cadetblue;*/
        margin-left: 10%;
        margin-right: 10%;
        height: auto;
    }

    #user-icon-logout {
        display: inline-flex;
        flex-direction: row;
        justify-content: center;
        /*background-color: red;*/
        width: 100%;
      height: auto;
    }


    /*
    */
    #profile-image {
      width: 300px;
      height: 300px;
      object-fit: cover;
      border-radius: 20%;
      margin: 10px;
      cursor: pointer;

      border-style: solid;
      border-width: 10px;
      border-color: cornflowerblue;
      transition: .1s;
      z-index: 0;
    }

    #profile-image:hover {
      transition: .1s;
      border-width: 5px;
      border-color: cornflowerblue;
    }

    .modal-overlay {
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      z-index: 98;
      background-color: rgba(0, 0, 0, 0.3);
    }

    .modal {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      z-index: 99;

      width: 100%;
      max-width: 400px;
      background-color: #FFF;
      border-radius: 16px;

      padding: 25px;
    }

    /*.modal2, .modal {*/
    /*  position: absolute;*/
    /*  top: 0;*/
    /*  right: 0;*/
    /*}*/

    .fade-enter-active,
    .fade-leave-active {
      transition: opacity .5s;
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

    /*
     */

    #username {
        font-weight: 400;
        font-size: 2rem;
        font-family: sans-serif;
    }

    #user-info {
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        /* text-align:left; */
        /* padding: 1rem 2rem;
        background: #fff;
        border: 1px solid #eee; */
    }

    #logout-settings {
        display: flex;
        flex-direction: row;
        justify-content: center;
    }

    #settings-password {
        display: flex;
        flex-direction: column;
        justify-content: center;
        min-width: 25rem;
    }

    #err {
        color: #f00;
    }

    #pass-str {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
    }

    #user-stories {
        flex-direction: column;
        height: auto;
        margin: 0 10%;
    }

    #user-media {
        justify-content: center;
        text-align: center;
        height: auto;
        margin-left: 10%;
        margin-right: 10%;
    }

    #description {

    }

    .story-highlight {
      width: 80px;
      height: 80px;
      object-fit: cover;
      border-radius: 20%;
      margin: 10px;
      cursor: pointer;

      border-style: solid;
      border-width: 2px;
      border-color: #323232;

      transition: .3s;
      z-index: 0;
    }

    /*.story-highlight:hover {*/
    /*  border-style: solid;*/
    /*  border-width: 5px;*/
    /*  border-color: #323232;*/

    /*  transition: .3s;*/
    /*}*/

    .story-highlight-layout {
      padding: 5px 10px;
      width: 150px;
      text-align: -webkit-center;
    }

    .inner-story-layout {
      display: flex;
      flex-wrap: nowrap;
      flex-direction: row;
      align-content: flex-start;
      width: fit-content;
    }

    .inner-post-layout {
      margin: 10px;
      display: flex;
      scroll-behavior: smooth;
    }

    .post-layout {
      width: 200px;
      height: 200px;
      margin: 5px;
      display: block;
      text-align: -webkit-center;
      justify-content: center;
    }

    .post {
      width: 180px;
      height: 180px;
      object-fit: cover;
      border-radius: 20%;
      margin: 10px;
      cursor: pointer;

      border-style: solid;
      border-width: 2px;
      border-color: #323232;

      transition: .3s;
      z-index: 0;
    }
</style>