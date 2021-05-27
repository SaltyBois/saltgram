<template>
  <div id="user-main">
    <portal-target name="drop-down-profile-menu" />
    <portal-target name="settings-menu"/>
    <TopBar style="position: sticky; z-index: 2"> </TopBar>
    <div id="user-header">
      <div id="user-icon-logout">

        <ProfileImage/>

        <v-layout column
                  style="width: 70%"
                  justify-space-between>
          <v-layout style="height: 40%;"
                    justify-center
                    column>

            <ProfileHeader/>

          </v-layout>

          <NameAndDescription/>

        </v-layout>

<!--          TODO: JOVAN SETTINGS-->
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
<!--        TODO: END OF JOVAN SETTINGS-->
      </div>
    </div>

<!--        TODO: STORY HIGHLIGHTS-->
    <v-layout id="user-stories"
              column>
      <v-layout class="inner-story-layout"
                style="margin: 10px">
        <StoryHighlight v-for="index in 10" :key="index"/>
      </v-layout>
    </v-layout >

    <!--  TODO: LAYOUT FOR TOGGLING: POSTS, SAVED, TAGGED  -->
    <v-layout id="radio-button-layout">
      <v-radio-group row  v-model="radioButton">
        <v-radio label="Posts"  value="posts"/>
        <v-radio label="Saved"  value="saved"/>
        <v-radio label="Tagged" value="tagged"/>
      </v-radio-group>
    </v-layout>

<!--        TODO: POSTS -->
    <transition name="fade">
      <v-layout class="user-media"
                v-if="radioButton === 'posts'"
                column>
        <PostOnUserPage/>

      </v-layout>
    </transition>


    <!--        TODO: SAVED -->
    <transition name="fade">
      <v-layout class="user-media"
                v-if="radioButton === 'saved'"
                column>
        <PostOnUserPage/>
        <PostOnUserPage/>
      </v-layout>
    </transition>

    <!--        TODO: TAGGED -->
    <transition name="fade">
      <v-layout class="user-media"
                v-if="radioButton === 'tagged'"
                column>
        <PostOnUserPage/>
        <PostOnUserPage/>
        <PostOnUserPage/>
      </v-layout>
    </transition>
  </div>
</template>

<script>
import passwordMeter from 'vue-simple-password-meter';
import TopBar from "@/components/TopBar";
import ProfileImage from "@/components/user_page_components/ProfileImage";
import ProfileHeader from "@/components/user_page_components/ProfileHeader";
import NameAndDescription from "@/components/user_page_components/NameAndDescription";
import StoryHighlight from "@/components/user_page_components/StoryHighlight";
import PostOnUserPage from "@/components/user_page_components/PostOnUserPage";

export default {
    components: {
      TopBar, passwordMeter, ProfileImage, ProfileHeader, NameAndDescription, StoryHighlight, PostOnUserPage
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
            radioButton: 'posts',
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
            return (this.newPassword1 !== this.oldPassword && this.newPassword2 !== this.oldPassword) || "Cannot use old password"
        },

        passStr: function() {
            console.log("pass score: " + this.passScore);
            return this.passScore > 3 || "Use a stronger password!";
        },

        passMatch: function() {
            return this.newPassword1 === this.newPassword2 || "Passwords must match"
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

    .user-media {
      --w:400px;
      --n:3;
      --m:2;

      margin: 0 10%;
      display:grid;
      grid-template-columns:repeat(auto-fit,minmax(clamp(100%/(var(--n) + 1) + 0.1%,(var(--w) - 100vw)*1000,100%/(var(--m) + 1) + 0.1%),1fr)); /*this */
      gap:10px;

      transition: 0.3s;
    }

    #radio-button-layout {
      height: 70px;
      text-align: -webkit-center;
      justify-content: center;
      margin: 0 10%;
      float: left;
     }

    .inner-story-layout {
      height: 150px;
      flex-direction: row;
      overflow-x: auto;
      overflow-y: hidden;
      white-space: nowrap;

    }

    .inner-post-layout > div {
      display: inline-block;
      color: white;
      text-align: center;
      padding: 14px;
      text-decoration: none;
    }

</style>