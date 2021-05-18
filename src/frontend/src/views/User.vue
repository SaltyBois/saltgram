<template>
    <div id="user-main">
      <transition name="fade" appear style="z-index: 11">
        <portal-target name="drop-down-profile-menu" />
      </transition>
      <top-bar style="position: sticky; top: 0; z-index: 10"/>
        <div id="user-header">
            <div id="user-icon-logout">
              <v-layout align-center column>
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
                        <v-btn class="accent"
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
                        style="background: hotpink"
                        justify-space-between>
                <v-layout style="background: gold; height: 0.1%"
                          justify-center
                          column>
                  <v-layout id="logout-settings" class="mt-5 " column>

<!--                    <v-btn @click="showSettings=!showSettings; showProfileSettingsDialog=!showProfileSettingsDialog" class="accent"> TODO(Mile): UZETI OVO U OBZIR!!! -->
                    <v-btn @click=" showProfileSettingsDialog=!showProfileSettingsDialog" class="accent">
                      <i class="fa fa-cog" aria-hidden="true"></i>
                      settings
                    </v-btn>
                    <transition name="fade" appear>
                      <div class="modal-overlay" v-if="showProfileSettingsDialog" @click="showProfileSettingsDialog = false"></div>
                    </transition>
                    <transition name="slide" appear>
                      <v-layout class="modal"
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

                        <v-btn @click="logout; showProfileSettingsDialog = false" class="error mb-3">Logout</v-btn>
                        <v-btn @click="showProfileSettingsDialog = false" class="accent">
                          Cancel
                        </v-btn>
                      </v-layout>
                    </transition>
                  </v-layout>
                  <v-layout row>
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
                          style="background: chartreuse">
                  <h3><b>Ime i prezime</b></h3>
                  <h4 v-if="descriptionTextBoxReadOnly">{{this.descriptionTextBoxText}}</h4>
                  <v-text-field v-else
                                v-model="descriptionTextBoxText"
                                style="width: 300px; padding: 5px; margin: 5px; height: 50%"/>
                  <v-btn v-if="descriptionTextBoxReadOnly"
                         @click="changeDescription"
                         class="mt-3 mb-3 align-content-md-stretch accent">Change description</v-btn>
                  <v-btn v-else
                         @click="changeDescription"
                         class="mt-3 mb-3 align-content-md-stretch accent">Confirm description</v-btn>
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
        <div id="user-stories"
             style="background: darkolivegreen">
            Stories
        </div>
        <div id="user-media"
             style="background: lightskyblue">
            Media
        </div>
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
        height: 100vh;
    }

    #user-header {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
        background: cadetblue;
        margin-left: 10%;
        margin-right: 10%;
    }

    #user-icon-logout {
        display: inline-flex;
        flex-direction: row;
        justify-content: center;
        background-color: red;
        width: 100%;
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
    }

    #profile-image:hover {
      transition: .1s;
      border-width: 5px;
      border-color: cornflowerblue;
    }

    .modal-overlay {
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      z-index: 98;
      background-color: rgba(0, 0, 0, 0.3);
    }

    .modal {
      position: fixed;
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
        justify-content: center;
        text-align: center;
        height: 70px;
        margin-left: 10%;
        margin-right: 10%;
    }

    #user-media {
        justify-content: center;
        text-align: center;
        margin-left: 10%;
        margin-right: 10%;
    }
</style>