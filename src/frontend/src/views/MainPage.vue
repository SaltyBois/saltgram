<template>
  <div style="background-color: #efeeee; height: 100%">
    <button class="go-up-button">
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
              <div class="story-layout">
                <div>
                  <v-img  class="no-story"
                          src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
                          alt="Profile picture">
                    <b class="story-plus-circle">
                      <h3 class="story-plus">+</h3>
                    </b>
                  </v-img>
                </div>
                <b>Add story</b>
              </div>
              <div class="story-layout">
                <v-img  class="story-close-friends"
                        src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
                        alt="Profile picture"/>
                <b>Your story</b>
              </div>
              <div class="story-layout">
                <v-img  class="story-seen"
                        src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
                        alt="Profile picture"/>
                <b>Your story</b>
              </div>
              <div class="story-layout">
                <v-img  class="story-close-friends"
                        src="https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg"
                        alt="Profile picture"/>
                <b>Username1</b>
              </div>
              <div class="story-layout">
                <v-img  class="story"
                        src="https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg"
                        alt="Profile picture"/>
                <b>Username1</b>
              </div>
            </v-layout>
          </v-layout >
        </div>
<!--        TODO: MILE-->
        <div id="posts-div"
             style="background-color: transparent;">
          <div style="margin-top: 5px">
            <v-btn class="follow-button" width="150px">+ new posts</v-btn>
          </div>
          <div class="post-card">
            <div class="post-header"
                 style="background-color: transparent">
              <div class="post-header-left-side">
                  <v-img  class="post-header-profile"
                          src="https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg"
                          @click="$router.push('/user')"
                          alt="Profile picture"/>
                <b @click="$router.push('/user')" style="cursor: pointer">Username1</b>
              </div>
              <div class="post-header-right-side">
                <b>...</b>
              </div>
            </div>
            <div class="post-content"
                 style="background-color: transparent">
              <v-img  class="post-content-media"
                      src="https://www.arabianbusiness.com/public/styles/square/public/images/2021/03/28/meme.jpg?itok=DeJVUtab"
                      alt="Profile picture"/>
            </div>
            <div class="post-interactions"
                 style="background-color: transparent">
              <div class="post-interactions-left-side">
                <div style="width: 50px; height: 50px; text-align: -webkit-center">
                  <i class="fa fa-thumbs-o-up like" aria-hidden="true"
                     />
                </div>
                <div style="width: 50px; height: 50px; text-align: -webkit-center;">
                  <i class="fa fa-thumbs-o-up dislike" aria-hidden="true"
                  />
                </div>

                <div style="width: 50px; height: 50px; text-align: -webkit-center">
                  <i class="fa fa fa-comment-o like" aria-hidden="true"
                  />
                </div>
              </div>
              <div class="post-interactions-right-side">
                <div style="width: 50px; height: 50px; text-align: -webkit-center">
                  <i class="fa fa-folder-open-o like" aria-hidden="true"
                  />
                </div>
              </div>
            </div>
            <div class="post-description"
                 style="background-color: transparent">
              <div style=" padding: 5px;">
                <p style="text-align: left; font-size: 12pt; margin-bottom: auto;">
                  <b>1234</b> Likes  <b>532</b> Dislikes
                <p style="text-align: left; font-size: 10pt; margin-bottom: auto;">
                  <b>USERNAME </b>
                  Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.
                </p>
              </div>
            </div>
            <div class="post-comment-section"
                 style="background-color: transparent">
              <div class="all-comments" >
                <p style="text-align: left; font-size: 10pt; margin-bottom: auto;">
                  View all <b>32</b> comments
                </p>
              </div>
              <p style="text-align: left; font-size: 10pt; margin-bottom: auto;">
                FIRST COMMENT
              </p>
              <p style="text-align: left; font-size: 10pt; margin-bottom: auto;">
                Second COMMENT
              </p>
              <p style="text-align: left; font-size: 10pt; margin-bottom: auto; color: #858585">
                Posted 1 hour ago
              </p>
            </div>
            <!--  TODO(Mile): Emojis need to be included GENERICALLY  -->
            <div class="post-footer"
                 style="background-color: transparent">
              <div style="float: left; height: available; display: flex; flex-direction: row; width: 80%">
                <EmojiPicker @emoji="append" :search="search">
                  <div
                      class="emoji-invoker"
                      slot="emoji-invoker"
                      slot-scope="{ events: { click: clickEvent } }"
                      @click.stop="clickEvent"
                  >
                    <svg height="24" viewBox="0 0 24 24" width="24" style="margin-top: 10px" xmlns="http://www.w3.org/2000/svg">
                      <path d="M0 0h24v24H0z" fill="none"/>
                      <path d="M11.99 2C6.47 2 2 6.48 2 12s4.47 10 9.99 10C17.52 22 22 17.52 22 12S17.52 2 11.99 2zM12 20c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8zm3.5-9c.83 0 1.5-.67 1.5-1.5S16.33 8 15.5 8 14 8.67 14 9.5s.67 1.5 1.5 1.5zm-7 0c.83 0 1.5-.67 1.5-1.5S9.33 8 8.5 8 7 8.67 7 9.5 7.67 11 8.5 11zm3.5 6.5c2.33 0 4.31-1.46 5.11-3.5H6.89c.8 2.04 2.78 3.5 5.11 3.5z"/>
                    </svg>
                  </div>

                  <div slot="emoji-picker" slot-scope="{ emojis, insert }" style="z-index: 10;">
                    <div class="emoji-picker" >
                      <div class="emoji-picker__search">
                        <input type="text" v-model="search" v-focus>
                      </div>
                      <div>
                        <div v-for="(emojiGroup, category) in emojis" :key="category">
                          <h5>{{ category }}</h5>
                          <div class="emojis">
                            <span
                                v-for="(emoji, emojiName) in emojiGroup"
                                :key="emojiName"
                                @click="insert(emoji)"
                                :title="emojiName"
                            >{{ emoji }}</span>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </EmojiPicker>
                <v-text-field label="Add a comment" style="width: available; padding: 5px" />
              </div>
              <div style="float: right; height: available; display: inline-block; ">
                <v-btn class="follow-button" style="margin: 5px; width: 75px">
                  post
                </v-btn>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div id="suggestions-div"
           style="background-color: transparent;">
        <div id="profile-info"
             style="background-color: transparent;">
          <v-img  class="profile-head"
                  src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
                  @click="$router.push('/user')"
                  alt="Profile picture"/>
          <b style="cursor: pointer;"
             @click="$router.push('/user')">USERNAME</b>
        </div>
        <div id="suggestions-info"
             style="background-color: transparent">
          <h4>Suggestions For You</h4>
          <div>
            <div class="suggested-profile">
              <v-img  class="suggestion-head"
                      src="https://i.pinimg.com/564x/4e/c4/f2/4ec4f2d69c9bc6b152abcb420252c3a8.jpg"
                      @click="$router.push('/user')"
                      alt="Profile picture"/>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <h5>@USERNAME</h5>
                <h5>Ime i prezime</h5>
                <small>45 mutual followers</small>
              </div>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <v-btn depressed
                       class="follow-button">
                  follow
                </v-btn>
              </div>
            </div>
            <div class="suggested-profile">
              <v-img  class="suggestion-head"
                      src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRkO7BR9U3GRatwXeEFkxxDRm7a5Kq7eGhJtg&usqp=CAU"
                      @click="$router.push('/user')"
                      alt="Profile picture"/>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <h5>@USERNAME</h5>
                <h5>Ime i prezime</h5>
                <small>45 mutual followers</small>
              </div>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <v-btn depressed
                       class="unfollow-button">
                  unfollow
                </v-btn>
              </div>
            </div>
            <div class="suggested-profile">
              <v-img  class="suggestion-head"
                      src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRYYR7aw1k6gkP_UGKGddjPee4ENRCEY8AO3A&usqp=CAU"
                      @click="$router.push('/user')"
                      alt="Profile picture"/>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <h5>@USERNAME</h5>
                <h5>Ime i prezime</h5>
                <small>45 mutual followers</small>
              </div>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <v-btn depressed
                       class="follow-button">
                  follow
                </v-btn>
              </div>
            </div>
            <div class="suggested-profile">
              <v-img  class="suggestion-head"
                      src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRYYR7aw1k6gkP_UGKGddjPee4ENRCEY8AO3A&usqp=CAU"
                      @click="$router.push('/user')"
                      alt="Profile picture"/>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <h5>@USERNAME</h5>
                <h5>Ime i prezime</h5>
                <small>45 mutual followers</small>
              </div>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <v-btn depressed
                       class="follow-button">
                  follow
                </v-btn>
              </div>
            </div>
            <div class="suggested-profile">
              <v-img  class="suggestion-head"
                      src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRYYR7aw1k6gkP_UGKGddjPee4ENRCEY8AO3A&usqp=CAU"
                      @click="$router.push('/user')"
                      alt="Profile picture"/>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <h5>@USERNAME</h5>
                <h5>Ime i prezime</h5>
                <small>45 mutual followers</small>
              </div>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <v-btn depressed
                       class="follow-button">
                  follow
                </v-btn>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import TopBar from "@/components/TopBar";
import EmojiPicker from 'vue-emoji-picker'

export default {
  name: "MainPage",
  components: {
    TopBar,
    EmojiPicker,
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

  #profile-info {
    display: inline-block;
    text-align: -webkit-center;
    margin: 10px;
    height: 150px;

  }

  #suggestions-info {
    display: inline-block;
    text-align: -webkit-left;
    margin: 10px;
    width: available;
  }

  .profile-head, .suggestion-head {
    width: 100px;
    height: 100px;
    object-fit: cover;
    border-radius: 20%;
    margin: 10px;
    cursor: pointer;

    border-style: solid;
    border-width: 2px;
    border-color: #323232;
    filter: brightness(1);

    transition: .3s;
    z-index: 0;
  }

  .profile-head:hover, .suggestion-head:hover {
    cursor: pointer;
    transition: .3s;
    filter: brightness(0.7);
  }

  .suggestion-head {
    width: 60px;
    height: 60px;
    margin: 0;
  }

  .story-layout {
    padding: 5px 10px;
    width: 150px;
    text-align: -webkit-center;
  }

  .inner-story-layout {
    height: 150px;
    flex-direction: row;
    overflow-x: auto;
    overflow-y: hidden;
    white-space: nowrap;
  }

  .story, .story-close-friends, .no-story, .story-seen, .post-header-profile {
    width: 80px;
    height: 80px;
    object-fit: cover;
    border-radius: 20%;
    margin: 10px;
    cursor: pointer;


    filter: brightness(1);

    transition: .3s;
    z-index: 0;
  }

  .story:hover, .story-close-friends:hover, .no-story:hover {
    transition: .3s;
    filter: brightness(0.7);
  }

  .story {
    border-style: solid;
    border-width: 5px;
    border-color: red;
  }

  .story-close-friends {
    border-style: solid;
    border-width: 5px;
    border-color: #36c400;
  }

  .post-header-profile {
    width: 30px;
    height: 30px;
  }

  .no-story {
    border-style: solid;
    border-width: 2px;
    border-color: #000000;
  }

  .story-plus-circle {
    position: relative;
    display: inline-block;
    top: 45px;
    left: 22px;
    border-radius: 50%;
    border-style: solid;
    border-width: 2px;
    border-color: #000000;
    width: 30px;
    height: 30px;
    background-color: #017ba9;
    color: white;

    font-size: 25px;
    text-align: -webkit-center;
  }

  .story-plus {
    position: relative;
    bottom: 10px;
    left: 0;
  }

  .story-seen {
    border-style: solid;
    border-width: 3px;
    border-color: #535353;

    filter: brightness(1.5);
    opacity: 0.8;
  }

  .suggested-profile {
    display: -webkit-inline-flex;
    flex-direction: row;
    height: auto;
    background-color: transparent;

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

  .post-card {
    margin: 10px 50px;
    background-color: white;
    width: available;
    height: auto;

    border: #323232 solid 1px;
  }

  .post-header {
    height: 51px;

    border-bottom-style: solid;
    border-bottom-width: 1px;
    border-bottom-color: #373737;

  }

  .post-content {
    width: 30vw;
    height: 30vw;

    align-content: center;

    object-fit: cover;
  }

  .post-interactions, post-footer {
      height: 51px;

    border-top-style: solid;
    border-top-width: 1px;
    border-top-color: #373737;
  }

  .post-description {
    height: auto;
  }

  .post-comment-section {
    height: auto;
    padding: 5px;
  }

  .post-footer {
    height: 50px;
    /*display: flex;*/
    /*flex-direction: row;*/
  }

  .post-header-left-side, .post-header-right-side, .post-interactions-left-side, .post-interactions-right-side .post-footer-right-side {
    direction: ltr;
    flex-direction: row;
    text-align: -webkit-center;
    align-items: center;
    float: left;
    display: flex;
    justify-content: center
  }

  .post-header-right-side, .post-interactions-right-side, .post-footer-right-side {
    float: right;
    width: 50px;
    height: 50px;
  }

  .post-content-media {
    position: relative;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);

    max-width:100%;
    max-height:100%;

  }

  .like, .dislike {
    position: relative;
    top: 12px;
    left: 0;
    transform: scale(2);
    margin: 0 3px;

    transition: 0.2s;
  }

  .dislike {
    transform: scale(2) rotate(180deg);
  }

  .like:hover {
    transition: 0.2s;
    color: #016ddb;
    cursor: pointer;
  }

  .dislike:hover {
    transition: 0.2s;
    color: #ff0000;
    cursor: pointer;
  }

  .all-comments {

  }

  .smiley {
    position: relative;
    top: 10px;
    left: 0px;
  }

  .smiley-div {
    width: 30px;
    height: 50px;
    text-align: center;
    transform: rotate(90deg);
    cursor: pointer;
  }

/*  TODO(NOTE): EMOJI TABLE CSS */
  .wrapper {
    position: relative;
    display: inline-block;
  }

  .regular-input {
    padding: 0.5rem 1rem;
    border-radius: 3px;
    border: 1px solid #ccc;
    width: 20rem;
    height: 12rem;
    outline: none;
  }

  .regular-input:focus {
    box-shadow: 0 0 0 3px rgba(66,153,225,.5);
  }

  .emoji-invoker {
    position: relative;
    top: 0;
    left: 0;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    cursor: pointer;
    transition: all 0.2s;
  }
  .emoji-invoker:hover {
    transform: scale(1.1);
  }
  .emoji-invoker > svg {
    fill: black;
  }

  .emoji-picker {
    position: relative;
    border: 1px solid #707070;
    width: 250px;
    height: 200px;
    overflow: scroll;
    padding: 0;
    box-sizing: border-box;
    border-radius: 5%;
    background: #fff;
    box-shadow: 1px 1px 8px #c7dbe6;
  }
  .emoji-picker__search {
    display: flex;
  }
  .emoji-picker__search > input {
    flex: 1;
    border-radius: 10rem;
    border: 1px solid #ccc;
    padding: 0.5rem 1rem;
    outline: none;
  }
  .emoji-picker h5 {
    margin-bottom: 0;
    color: #b1b1b1;
    text-transform: uppercase;
    font-size: 0.8rem;
    cursor: default;
  }
  .emoji-picker .emojis {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
  }
  .emoji-picker .emojis:after {
    content: "";
    flex: auto;
  }
  .emoji-picker .emojis span {
    padding: 0.2rem;
    cursor: pointer;
    border-radius: 5px;
  }
  .emoji-picker .emojis span:hover {
    background: #ececec;
    cursor: pointer;
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