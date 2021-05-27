<template>
  <div style="background-color: #efeeee; min-height: fit-content">
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
                      style=";background-color: transparent;">

              <MyNoStory/>

              <MyCloseFriendStory/>

              <MyStory/>

              <MySeenStory/>

              <Story/>

              <StoryCloseFriends/>

              <StorySeen/>

              <StoryMuted/>

            </v-layout>
          </v-layout >
        </div>
<!--        TODO: MILE-->
        <div id="posts-div"
             style="background-color: transparent;">

          <div style="margin-top: 5px">
            <v-btn class="new-posts-button" width="150px">+ new posts</v-btn>
          </div>

          <PostOnMainPage v-for="index in 3" :key="index"/>

        </div>
      </div>
      <div id="suggestions-div"
           style="background-color: transparent;">

        <UserProfileHead/>

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
    }
  },
  methods: {
    insert(emoji) {
      this.input += emoji
    },
    append(emoji) {
      this.input += emoji
    },
    scrollUp() {
      window.scrollTo(0, 0)
    }
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
  }

</style>