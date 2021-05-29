<template>
  <div>
    <transition name="fade" appear>
      <v-layout class="modal-overlay"
                align-center
                v-if="show"
                @click="toggleParrent"
                justify-center/>
    </transition>
    <transition name="fade" appear>
      <v-layout class="modal"
                justify-center
                align-center
                v-if="show"
                column>
        <!--    TODO(MILE): POST-->
<!--        <div id="posts-div">-->
          <div class="post-card">
            <div class="post-header">
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
            <div style="display: flex; justify-content: space-between">
              <v-layout class="post-content" align-center justify-center style="background-color: transparent">
                <v-img  class="post-content-media"
                        :src="mediaPath"
                        alt="Post content"/>
              </v-layout>
              <v-layout style="display: flex; flex-direction: column;" >
                <div class="post-description">
                  <div style=" padding: 5px;">
                    <p style="text-align: left; font-size: 10pt; margin-bottom: auto;">
                      <b>USERNAME</b>
                      Lorem Ipsum is simply dummy text of the printing and typesetting industry.
                    </p>

                  </div>
                </div>
                <div class="post-comment-section">
                  <p style="text-align: left; font-size: 10pt; margin-bottom: auto;">
                    FIRST COMMENT
                  </p>
                  <p style="text-align: left; font-size: 10pt; margin-bottom: auto;">
                    Second COMMENT
                  </p>

                </div>
                <!--  TODO(Mile): Emojis need to be included GENERICALLY  -->
                <div class="post-footer">
                  <div class="post-interactions"
                       style="background-color: transparent">
                    <div class="post-interactions-left-side">
                      <div style="width: 50px; height: 50px; text-align: -webkit-center">
                        <i class="fa fa-thumbs-o-up like" aria-hidden="true"/>
                      </div>
                      <div style="width: 50px; height: 50px; text-align: -webkit-center;">
                        <i class="fa fa-thumbs-o-up dislike" aria-hidden="true"/>
                      </div>
                      <div style="width: 50px; height: 50px; text-align: -webkit-center">
                        <i class="fa fa fa-comment-o like" aria-hidden="true"/>
                      </div>
                    </div>
                    <div class="post-interactions-right-side">
                      <div style="width: 50px; height: 50px; text-align: -webkit-center">
                        <i class="fa fa-folder-open-o like" aria-hidden="true"/>
                      </div>
                    </div>
                  </div>

                  <div style=" padding: 5px;">
                    <p style="text-align: left; font-size: 12pt; margin-bottom: auto;">
                      <b>1234</b> Likes  <b>532</b> Dislikes
                    </p>
                    <p style="text-align: left; font-size: 10pt; margin-bottom: auto; color: #858585">
                      Posted 1 hour ago
                    </p>
                  </div>

                  <div style="float: left; height: available; display: flex; flex-direction: row; width: 80%; margin-bottom: 10px">
                    <EmojiPicker @emoji="append" :search="search">
                      <div
                          class="emoji-invoker"
                          slot="emoji-invoker"
                          slot-scope="{ events: { click: clickEvent } }"
                          @click.stop="clickEvent">
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
                    <v-btn class="post-button" style="margin: 5px; width: 75px">
                      post
                    </v-btn>
                  </div>
                </div>
              </v-layout>
            </div>
          </div>
<!--        </div>-->
      </v-layout>
    </transition>
  </div>
</template>

<script>
export default {
  name: "PostView",
  props: {
      mediaPath: {
        required: true,
        type: String
      }
  },
  methods: {
    toggleParrent() {
      this.show = !this.show
    }
  },
  data: function () {
    return {
      show: false,
      search: '',
    }
  }
}
</script>

<style scoped>

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 98;
  background-color: rgba(0, 0, 0, 0.3);
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

#posts-div {
  display: flex;
  height: auto;
  min-width: 90%;
  max-width: 90%;
  max-height: 700px;
  flex-direction: column;
  text-align: -webkit-center;
}

.post-card {
  /*margin: 10px 10px;*/
  background-color: white;
  width: 98%;
  height: auto;

  border: #323232 solid 1px;

  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;
}

.post-header {
  height: 51px;

  border-bottom-style: solid;
  border-bottom-width: 1px;
  border-bottom-color: #373737;
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
  width: 30vw;
  height: 30vh;

  max-width:90%;
  max-height:90%;

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


.post-content {
  width: 30vw;
  height: 30vw;

  object-fit: cover;

  border-right: black 1px solid;
}

.post-interactions {
  height: 51px;

}

.post-description {
  height: auto;
}

.post-comment-section {
  height: 60%;
  padding: 5px;
  overflow-x: hidden;
  overflow-y: auto;
  scroll-behavior: smooth;
}

.post-footer {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  /*min-height: 60px;*/
  /*max-height: 60px;*/
  height: auto;

  border-top: black 1px solid;

  /*margin-bottom: 200px;*/

  /*display: flex;*/
  /*flex-direction: row;*/
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

.post-button {
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

.modal {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 99;

  width: 75%;
  height: 75%;
  background-color: #FFF;
  border-radius: 16px;

  padding: 5px;
}

</style>