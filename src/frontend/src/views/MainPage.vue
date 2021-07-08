<template>
  <div class="background-div">
    <button class="go-up-button" @click="scrollUp">
      <i class="fa fa-level-up" style="transform: scale(2)"/>
    </button>
    <portal-target name="drop-down-profile-menu" />
    <portal-target name="settings-menu"/>
    <TopBar style="position: sticky; z-index: 2"/>
    <div id="main-div"
         style="background-color: transparent">

      <div id="media-div"
           style="background-color: transparent;">
        <div id="stories-div"
             style="background-color: transparent;">
          <v-layout id="user-stories"
                    style="background-color: transparent;"
                    column>
            <v-layout class="inner-story-layout"
                      style="background-color: transparent;">

              <MyNoStory v-if="!myStoriesExist" :user="loggedUser"/>

              <MyCloseFriendStory v-if="myStoriesExist && myStories.closeFriends" :user="loggedUser" :stories="myStories"/>

              <MyStory v-if="myStoriesExist && !myStories.closeFriends" :user="loggedUser" :stories="myStories" />

<!--              <MySeenStory/>-->
              <div v-for="(item, index) in pageStories" :key="index">
                <Story v-if="!item.closeFriends && item.storyElement.length !== 0" :user="item.user" :stories="item.storyElement" />

                <StoryCloseFriends v-else-if="item.closeFriends && item.storyElement.length !== 0" :user="item.user" :stories="item.storyElement"/>

<!--                <StorySeen/>-->

<!--                <StoryMuted />-->
              </div>
            </v-layout>
          </v-layout >
        </div>
<!--        TODO: MILE-->
        <div id="posts-div"
             style="background-color: transparent;">

<!--          <div style="margin-top: 5px">-->
<!--            <v-btn class="new-posts-button" width="150px">+ new posts</v-btn>-->
<!--          </div>-->

          <PostOnMainPage v-for="(item, index) in pagePostsSorted"
                          :user="item.user"
                          :post-element="item.postElement"
                          :key="index"/>

        </div>
      </div>
      <div id="suggestions-div"
           style="background-color: transparent;">

        <UserProfileHead :user="loggedUser"/>

        <div id="suggestions-info"
             style="background-color: transparent">
          <h4>Suggestions For You</h4>
          <div>

            <SuggestedProfile v-for="index in 5" :key="index"/>

          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import TopBar from "@/components/TopBar";
import EmojiPicker from 'vue-emoji-picker'
import PostOnMainPage from "@/components/main_page_components/PostOnMainPage";
import MyNoStory from "@/components/main_page_components/MyNoStory";
import MyCloseFriendStory from "@/components/main_page_components/MyCloseFriendStory";
import MySeenStory from "@/components/main_page_components/MySeenStory";
import MyStory from "@/components/main_page_components/MyStory";
import Story from "@/components/main_page_components/Story";
import StoryCloseFriends from "@/components/main_page_components/StoryCloseFriends";
import StoryMuted from "@/components/main_page_components/StoryMuted";
import UserProfileHead from "@/components/main_page_components/UserProfileHead";
import SuggestedProfile from "@/components/main_page_components/SuggestedProfile";
import StorySeen from "@/components/main_page_components/StorySeen";

export default {
  name: "MainPage",
  components: {
    TopBar,
    EmojiPicker,
    PostOnMainPage,
    MyNoStory,
    MyCloseFriendStory,
    MySeenStory,
    MyStory,
    Story,
    StoryCloseFriends,
    StoryMuted,
    UserProfileHead,
    SuggestedProfile,
    StorySeen,
  },
  data() {
    return {
      input: '',
      search: '',
      followingUsers: [],
      pagePosts: [],
      pagePostsSorted: [],
      pageStories: [],
      myStories: {},
      myStoriesExist: null,
      postsInfo: {},
      loggedUser: {},
    }
  },
  methods: {
    scrollUp() {
      window.scrollTo(0, 0)
    },
    getLoggedUserInfo() {
      this.axios.get('users', {headers: this.getAHeader()})
          .then(r => {
            this.loggedUser = r.data;
            // console.log(this.loggedUser)
            this.getFollowingUsers();
          })
          .catch(err => {
            console.log(err)
          })
    },
    getFollowingUsers: function() {
      this.axios.get("users/following/main/", {headers: this.getAHeader()})
          .then(r => {
            // console.log(r.data);
            this.followingUsers = r.data;
            // console.log(this.followingUsers)
            this.myStories = {};
            this.getUserPosts(this.loggedUser, 0)
            this.getUserStories(this.loggedUser, -1)
            for(let i = 0; i < this.followingUsers.length; ++i) {
              let currentUser = this.followingUsers[i];
              this.getUserPosts(currentUser, i);
              this.getUserStories(currentUser, i);
            }
          })
    },
    getUserPosts: function (currentUser, i) {
      // console.log('currentUser Username: ' + currentUser.username)
      // console.log('currentUser ID: ' + currentUser.id)
      this.axios.get("content/post/" + currentUser.id, {headers: this.getAHeader()})
          .then(r => {
            // console.log(r.data);
            r.data.forEach(el => {
              let postElement = {
                user: currentUser,
                postElement: el
              }
              this.pagePosts.push(postElement)
            })
            // console.log('i: ', i)
            // console.log('this.followingUsers.length: ', this.followingUsers.length)
            if (this.followingUsers !== null) {
              if ( i + 1 === this.followingUsers.length) {
                this.sortPosts()
              }
            }
          }).catch(err => {
            console.log(err)
        // this.$router.push('/');
      })
    },
    sortPosts() {
      this.pagePosts.sort(function (a,b) {
        let index1 = a.postElement.post.sharedMedia.media[0].addedOn.indexOf('CEST') + 4
        let index2 = b.postElement.post.sharedMedia.media[0].addedOn.indexOf('CEST') + 4
        let d1 = new Date(a.postElement.post.sharedMedia.media[0].addedOn.substring(0, index1).replace('CEST', '(CEST)'))
        let d2 = new Date(b.postElement.post.sharedMedia.media[0].addedOn.substring(0, index2).replace('CEST', '(CEST)'))
        if (d1 < d2) {
          return 1;
        }
        if (d1 > d2) {
          return -1;
        }
        // dates must be equal
        return 0;
      });
      this.pagePostsSorted = this.pagePosts;
      // console.log(this.pagePosts)
    },
    getUserStories(currentUser, i) {
        // console.log(currentUser, i)
      this.axios.get("content/story/" + currentUser.id, {headers: this.getAHeader()})
          .then(r => {
            console.log(r.data)
            let validStories = []
            const oneDay = 60 * 60 * 24 * 1000;
            r.data.forEach(el1 => {
              el1.stories.forEach(el2 => {
                let index = el2.addedOn.indexOf('CEST') + 4
                let storyDate = new Date(el2.addedOn.substring(0, index).replace('CEST', '(CEST)'))
                // console.log(storyDate)
                if ((Date.now() - storyDate) < oneDay) validStories.push(el1)
              })
            })

            // console.log(validStories)

            let storyElement = {
              user: currentUser,
              storyElement: validStories,
              closeFriends: false,
            }

            storyElement.storyElement.forEach(el => {
              if (el.closeFriends) {
                storyElement.closeFriends = true;
              }
            })

            if (i === -1) {
              this.myStories = storyElement;
              if (this.myStories.storyElement.length !== 0) this.myStoriesExist = true;
              else this.myStoriesExist = false;
              console.log('IDE GAS: ', this.myStories)
              // console.log(this.myStoriesExist)
            }
            else this.pageStories.push(storyElement);

          }).catch(err => {
            console.log(err)
            this.$router.push('/');
      })
    }
  },
  mounted() {
    this.getLoggedUserInfo()
    // window.setTimeout(function () {this.sortPosts()}, 100)
  },
  directives: {
    focus: {
      inserted(el) {
        el.focus()
      },
    },
  },
}
</script>

<style scoped>

  .background-div {
    background-color: #e9e9e9;
    min-height: auto;
    height: 500vh;

  }

  #main-div {
    display: flex;
    margin: 0 10%;
    height: 1000px;
    flex-direction: row;
  }

  #media-div {
    display: flex;
    flex-direction: column;

    width: 70%;
  }

  #stories-div {
    display: flex;
    height: 150px;
    flex-direction: row;

  }

  #posts-div {
    display: flex;
    height: auto;
    flex-direction: column;
    text-align: -webkit-center;
  }

  #suggestions-div {
    /*position: flex;*/
    left: 70%;
    display: flex;
    flex-direction: column;
    width: 30%;
    height: border-box;
  }

  #suggestions-info {
    display: inline-block;
    text-align: -webkit-left;
    margin: 10px;
    width: available;
  }

  .inner-story-layout {
    height: 150px;
    flex-direction: row;
    overflow-x: auto;
    overflow-y: hidden;
    white-space: nowrap;
  }

  .new-posts-button  {
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

  .go-up-button {
    position: fixed;
    color: #016ddb;
    background-color: white;
    border: #016ddb solid 1px ;
    width: 50px;
    height: 50px;
    top: 90%;
    left: 95%;
    border-radius: 50%;
    transition: 0.3s;
  }

  .go-up-button:hover {
    transform: scale(1.2);
    transition: 0.3s;
  }

</style>