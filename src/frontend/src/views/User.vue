<template>
  <div id="user-main">
    <portal-target name="drop-down-profile-menu" />
    <portal-target name="settings-menu"/>
    <TopBar style="position: sticky; z-index: 2"/>
    <div id="user-header">
      <div id="user-icon-logout">
        <v-layout align-center
                  justify-center>
          <ProfileImage ref="profileImage" 
          :is-my-profile-prop="this.isMyProfile" 
          :username="this.profile.username" 
          image-src="Insert image source" 
          @toggle-following="toggleFollow" 
          />
        </v-layout>
        <v-layout column
                  style="width: 70%"
                  justify-space-between>
          <v-layout style="height: 40%;"
                    justify-center
                    column>

            <ProfileHeader :following-prop="this.profile.following" :followers-prop="this.profile.followers" :user-prop="this.user"/>

          </v-layout>

          <NameAndDescription :name="this.profile.fullName" :description="this.profile.description" :web-site="this.profile.webSite"/>

        </v-layout>

      </div>
    </div>

    <div v-if="!isContentVisible" class="private-account">
      <i class="fa fa-lock" style="transform: scale(2.5)"/>
      <h3>This user is private</h3>

    </div>

<!--        TODO: STORY HIGHLIGHTS-->
    <v-layout id="user-stories"
              v-if="isContentVisible"
              column>
      <v-layout class="inner-story-layout"
                style="margin: 10px">
        <StoryHighlight v-for="index in 10" :key="index"/>
      </v-layout>
    </v-layout >

    <!--  TODO: LAYOUT FOR TOGGLING: POSTS, SAVED, TAGGED  -->
    <v-layout id="radio-button-layout"
              v-if="isContentVisible">
      <v-radio-group row  v-model="radioButton">
        <v-radio label="Posts"  value="posts"/>
        <v-radio label="Saved"  value="saved"/>
        <v-radio label="Tagged" value="tagged"/>
      </v-radio-group>
    </v-layout>

<!--        TODO: POSTS -->
    <transition name="fade">
      <v-layout class="user-media"
                v-if="radioButton === 'posts' && isContentVisible"
                column>
        <PostOnUserPage/>
      </v-layout>
    </transition>


    <!--        TODO: SAVED -->
    <transition name="fade">
      <v-layout class="user-media"
                v-if="radioButton === 'saved' && isContentVisible"
                column>
        <PostOnUserPage/>
        <PostOnUserPage/>
      </v-layout>
    </transition>

    <!--        TODO: TAGGED -->
    <transition name="fade">
      <v-layout class="user-media"
                v-if="radioButton === 'tagged' && isContentVisible"
                column>
        <PostOnUserPage/>
        <PostOnUserPage/>
        <PostOnUserPage/>
      </v-layout>
    </transition>
  </div>
</template>

<script>
import TopBar from "@/components/TopBar";
import ProfileImage from "@/components/user_page_components/ProfileImage";
import ProfileHeader from "@/components/user_page_components/ProfileHeader";
import NameAndDescription from "@/components/user_page_components/NameAndDescription";
import StoryHighlight from "@/components/user_page_components/StoryHighlight";
import PostOnUserPage from "@/components/user_page_components/PostOnUserPage";

export default {
    components: {
      TopBar, ProfileImage, ProfileHeader, NameAndDescription, StoryHighlight, PostOnUserPage
    },
    data: function() {
      return {
        profile : {
          privateUser: true,
          description: '',
          fullName: '',
          followers: '',
          following: '',
          followersList:[],
          followingList: [],
          username: '',
          webSite: ''
        },
        isMyProfile: false,
        radioButton: 'posts',
        followingUser: false,
        pendingRequest: false,
      }
    },
    computed: {
      isContentVisible() {
          return !(!this.isMyProfile && this.profile.privateUser && !this.followingUser);

      }
    },
    methods: {
        getUserInfo: function() {
            this.refreshToken(this.getAHeader())
                .then(rr => {
                    this.$store.state.jws = rr.data;
                    this.axios.get("users", {headers: this.getAHeader()})
                        .then(r =>{ 
                          this.user = r.data
                          this.getUser();
                          });
                      
                }).catch(() => this.$router.push('/'));
        },

        getUser: function() {
            this.isMyProfile = this.user.username === this.$route.params.username;
            this.$refs.profileImage.$data.isMyProfile = this.isMyProfile
            this.followingUser = false;
            this.privateUser = true;


            this.axios.get("users/profile/" + this.$route.params.username, {headers: this.getAHeader()})
            .then(r => {
              // console.log(r.data)
              this.profile.privateUser = !r.data.isPublic;
              this.profile.followingUser = r.data.isFollowing;
              this.profile.username = r.data.username;
              this.profile.followers = r.data.followers;
              this.profile.following = r.data.following;
              this.profile.fullName = r.data.fullName;
              this.profile.description = r.data.description;
              this.profile.webSite = r.data.webSite;
            }).catch(err => {
              console.log(err)
              console.log('Pushing Back to Login Page after fetching profile')
              this.$router.push('/');
            })

            if(!this.isMyProfile) {
              this.axios.get("users/check/follow/" + this.$route.params.username, {headers: this.getAHeader()})
              .then(r => {
                this.followingUser = r.data
                this.$refs.profileImage.$data.following = r.data
              })
        
              if(!this.followingUser) {
                this.axios.get("users/check/followrequest/" + this.$route.params.username, {headers: this.getAHeader()})
                .then(r => {
                  //this.pendingRequest = r.data
                  this.$refs.profileImage.$data.waitingForResponse = r.data
                })
              }
            }

          // this.axios.get("users/get/followers/" + this.$route.params.username, {headers: this.getAHeader()})
          //     .then(r => {
          //       this.profile.followersList = r.data
          //     }).catch(() => {
          //   console.log('Pushing Back to Login Page after failed fetching followers')
          //   this.$router.push('/');
          // })
          //
          // this.axios.get("users/get/following/" + this.$route.params.username, {headers: this.getAHeader()})
          //     .then(r => {
          //       this.profile.followingList = r.data
          //       console.log('this.profile.followingList')
          //       console.log(this.profile.followingList)
          //       this.$refs.profileImage.$data.following = this.profile.followingList[this.$route.params.username] !== 'true';
          //       console.log('this.$refs.profileImage.$data.following')
          //       console.log(this.$refs.profileImage.$data.following)
          //       // for (let i = 0; i < this.profile.followingList.length; ++i) {
          //       //   if (this.$route.params.username === this.profile.followingList[i]) {
          //       //     this.$refs.profileImage.$data.following = true
          //       //     break;
          //       //   }
          //       // }
          //     }).catch(r => {
          //   console.log('Pushing Back to Login Page after failed fetching following')
          //   this.$router.push('/');
          // })


        },
        toggleFollow(follow) {
          this.followingUser = follow
        },
    },
    mounted() {
         this.getUserInfo(); // TODO UNCOMMENT THIS
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
        min-height: 100vh;
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

    .private-account {
      text-align: -webkit-center;
      padding-top: 50px;
      height: 100px;
    }

</style>
