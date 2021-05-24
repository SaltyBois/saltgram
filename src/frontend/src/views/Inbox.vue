<template>
  <div style="background-color: #efeeee; height: 100%">
    <portal-target name="drop-down-profile-menu" />
    <portal-target name="settings-menu"/>
    <top-bar style="position: sticky; z-index: 2"/>
    <div id="main-inbox-div">
      <div id="inbox-list">
        <div class="inbox-list-header" >
          <div class="inbox-list-header-username">
            <h2 style="float:left;">Username</h2>
          </div>
          <div class="new-message">
            <i class="fa fa-envelope-o new-message-button"/>
          </div>
        </div>
        <div class="inbox-list-conversations">
          <div class="chat-div">
            <div style="padding: 3%;">
              <v-img  class="chat-head"
                      src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRYYR7aw1k6gkP_UGKGddjPee4ENRCEY8AO3A&usqp=CAU"
                      @click="$router.push('/user')"
                      alt="Profile picture"/>
            </div>
            <div style="padding: 3%;overflow: hidden">
              <h4>Username1</h4>
              <p>Message content &#183; 10m ago</p>
            </div>
          </div>
          <div class="chat-div">
            <div style="padding: 3%">
              <v-img  class="chat-head"
                      src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
                      @click="$router.push('/user')"
                      alt="Profile picture"/>
            </div>
            <div style="padding: 3%;overflow: hidden">
              <b>
                <h3>Username1</h3>
                <p>Message content 	&#183; now</p>
              </b>
            </div>
            <div class="blue-dot">
              <p disabled="true" style="transform: scale(5); color: #016ddb; margin-top: 15px">&#183;</p>
            </div>
          </div>
        </div>
      </div>
      <div id="inbox-content">
        <div class="inbox-content-header">
          <div style="float: left; text-align: -webkit-center; height: 60px; display: flex">
            <v-img  class="profile-chat-head"
                    src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
                    @click="$router.push('/user')"
                    alt="Profile picture"/>
            <h3 style="margin-top: 15px; margin-right: 10px">Username1</h3>
          </div>
          <div style="height: 60px; text-align: right">
            <i class="fa fa-info info-button" />
          </div>

        </div>
        <div class="inbox-content-body" ref="chatBox">
          <div class="inbox-message-div-left">
            <div style="float: left; text-align: -webkit-center; height: 60px; display: flex">
              <v-img  class="profile-chat-head-message"
                      src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
                      @click="$router.push('/user')"
                      alt="Profile picture"/>
            </div>
            <div class="inbox-message-received">
              AAAAAajsdjhkasdhjklasd;kjhask;hdjas;jkhdas
            </div>
            <div class="message-time-received">
              15.4.2015. 17:45
            </div>
          </div>
          <div class="inbox-message-div-right">
            <div class="inbox-message-sent">
              BBBBklhdsakjadjhasdsa;jkdskadhkasd;jksad;jhksa;jkh
            </div>
            <div class="message-time-sent">
              15.4.2015. 17:46
            </div>
          </div>
          <div class="inbox-message-div-right">
            <div class="inbox-message-sent">
              Hey! How are you Garry?
            </div>
            <div class="message-time-sent">
              15.4.2015. 17:46
            </div>
          </div>
          <div class="inbox-message-div-left">
            <div style="float: left; text-align: -webkit-center; height: 60px; display: flex">
              <v-img  class="profile-chat-head-message"
                      src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
                      @click="$router.push('/user')"
                      alt="Profile picture"/>
            </div>
            <div class="inbox-message-received">
              Hey I am at home. Where are you? How are you doing?
            </div>
            <div class="message-time-received">
              15.4.2015. 17:45
            </div>
          </div>
          <div class="inbox-message-div-left">
            <div style="float: left; text-align: -webkit-center; height: 60px; display: flex">
              <v-img  class="profile-chat-head-message"
                      src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
                      @click="$router.push('/user')"
                      alt="Profile picture"/>
            </div>
            <div class="inbox-message-received">
              Hey I am at home. Where are you? How are you doing?
            </div>
            <div class="message-time-received">
              15.4.2015. 17:45
            </div>
          </div>
          <div class="inbox-message-div-right">
            <div class="inbox-message-sent">
              Fine thanks. I am at gym right now.
            </div>
            <div class="message-time-sent">
              15.4.2015. 17:46
            </div>
          </div>
          <div class="inbox-message-div-right">
            <div class="inbox-message-sent">
              Fine thanks. I am at gym right now.
            </div>
            <div class="message-time-sent">
              15.4.2015. 17:46
            </div>
          </div>
          <div class="inbox-message-div-right">
            <div class="inbox-message-sent">
              Fine thanks. I am at gym right now.
            </div>
            <div class="message-time-sent">
              15.4.2015. 17:46
            </div>
          </div>
          <div class="inbox-message-div-right">
            <div class="inbox-message-sent">
              Fine thanks. I am at gym right now.
            </div>
            <div class="message-time-sent">
              15.4.2015. 17:46
            </div>
          </div>
          <div class="inbox-message-div-right">
            <div class="inbox-message-sent">
              Fine thanks. I am at gym right now.
            </div>
            <div class="message-time-sent">
              15.4.2015. 17:46
            </div>
          </div>
          <div class="inbox-message-div-right">
            <div class="inbox-message-sent">
              Fine thanks. I am at gym right now.
            </div>
            <div class="message-time-sent">
              15.4.2015. 17:46
            </div>
          </div>
          <div class="inbox-message-div-right">
            <div class="inbox-message-sent">
              <v-img src="https://metro.co.uk/wp-content/uploads/2020/08/IMG_8718-ec27.jpg?quality=90&strip=all"/>
            </div>
            <div class="message-time-sent">
              15.4.2015. 17:46
            </div>
          </div>
        </div>
        <div class="inbox-content-footer">
          <div class="inner-inbox-content-footer">
            <EmojiPicker @emoji="append" :search="search">
              <div
                  class="emoji-invoker"
                  slot="emoji-invoker"
                  slot-scope="{ events: { click: clickEvent } }"
                  @click.stop="clickEvent">
                <svg height="24" viewBox="0 0 24 24" width="24" style="margin-top: 15px" xmlns="http://www.w3.org/2000/svg">
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
            <i class="fa fa-photo picture-button" @click="$refs.file.click()"/>
            <input type="file"
                   ref="file"
                   style="display: none"
                   @change="onSelectedFile"
                   accept="image/*">
            <v-text-field label="Enter a message" style="padding-left: 5px; margin-bottom: 5px;" />

            <div style="float: right; height: available; display: inline-block;">
              <v-btn class="follow-button" style="margin: 10px 5px 5px;width: 75px">
                Send
              </v-btn>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import TopBar from "@/components/TopBar";
import EmojiPicker from "vue-emoji-picker";
export default {
name: "Inbox",
  components: {TopBar, EmojiPicker},
  methods: {
    onSelectedFile(event) {
      console.log(event)
      this.profilePicture = event.target.files[0]
      console.log(this.profilePicture)
    },
  },
  mounted() {
    var content = this.$refs.chatBox;
    content.scrollTop = content.scrollHeight;
  },

}
</script>

<style scoped>

#main-inbox-div {
  background-color: white;
  /*border-radius: 0.5%;*/
  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;
  display: flex;
  margin: 20px 10%;
  height: 80vh;
  flex-direction: row;
  border: #323232 solid 1px;
}

#inbox-list {
  display: flex;
  flex-direction: column;
  width: 40%;

}

#inbox-content {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
  width: 60%;
  /*display: grid;*/
  /*grid-template-rows: 60px 500px 60px;*/
  /*flex-flow: column;*/
  /*float: left;*/
}

.inbox-list-header {
  height: 60px;
  display: block;

  border-bottom-style: solid;
  border-bottom-color: black;
  border-bottom-width: 1px;

  border-right-style: solid;
  border-right-color: black;
  border-right-width: 1px;
}

.inbox-list-header-username {
  width: 75%;
  float: left;
  text-align: -webkit-center;
  margin-top: 10px;
  margin-left: 5px;
  overflow-x: -moz-hidden-unscrollable;
  scroll-behavior: smooth;
  white-space: nowrap;
  overflow: hidden;
}

.new-message {
  width: 60px;
  height: 60px;
  float: right;
  text-align: -webkit-center;
}

.new-message-button {
  transform: scale(2);
  margin-top: 20px;
  cursor: pointer;

  transition: 0.3s;
}

.new-message-button:hover {
  color: #016ddb;
  transition: 0.3s;
}

.inbox-list-conversations {
  overflow-y: auto;
  overflow-x: hidden;

  height: 100%;
  width: 100%;

  border-right: black 1px solid;
}

.chat-head {
  width: 60px;
  height: 60px;
  object-fit: cover;
  border-radius: 20%;
  cursor: pointer;

  border-style: solid;
  border-width: 2px;
  border-color: #323232;
  filter: brightness(1);

  transition: .3s;
  z-index: 0;
}

.chat-head:hover {
  cursor: pointer;
  transition: .3s;
  filter: brightness(0.7);
}

.chat-div {
  display: flex;
  flex-direction: row;

  height: 80px;
  width: 100%;
  background-color: transparent;
  cursor: pointer;
  transition: 0.3s;
}

.chat-div:hover {
  background-color: #949494;
  transition: 0.3s;
}

.chat-div:checked {
  background-color: #4b90b5;
  transition: 0.3s;
}

.blue-dot {
  position: relative;
  top: 0%;
  left: 25%;

  width: auto;
  height: auto;

  float: right;
  padding: 3%;

}

.inbox-content-header {
  height: 60px;
  width: 100%;

  float: top;

  display: block;
  vertical-align: center;

  /*background-color: red;*/


  border-bottom: black 1px solid ;
}

.inbox-content-body {
  overflow-y: auto;
  overflow-x: hidden;

  height: available;
  width: 100%;


  /*background-color: #002aff;*/
}

.inbox-content-footer {
  height: 60px;
  width: 100%;
  /*float: bottom;*/
  /*background-color: gold;*/

  border-top: black 1px solid;
}

.profile-chat-head, .profile-chat-head-message {
  width: 30px;
  height: 30px;
  object-fit: cover;
  border-radius: 20%;
  cursor: pointer;
  margin-top: 15px;
  margin-left: 10px;
  margin-right: 10px;

  border-style: solid;
  border-width: 2px;
  border-color: #323232;
  filter: brightness(1);

  transition: .3s;
  z-index: 0;
}

.profile-chat-head-message {
  margin: 5px 5px 0 5px;
}

.profile-chat-head:hover {
  cursor: pointer;
  transition: .3s;
  filter: brightness(0.7);
}

.info-button {
  transform: scale(2);
  margin-right: 20px;
  margin-top: 20px;
  cursor: pointer;

  transition: 0.3s;
}

.info-button:hover {
  color: #016ddb;
  transition: 0.3s;
}

.inner-inbox-content-footer {
  height: 60px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

.follow-button {
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

.inbox-message-div-left, .inbox-message-div-right {
  display: inline-block;
  height: auto;
  width: 100%;
  align-content: center;
  margin: 3px 0;
  /*padding: 20px;*/
}

.inbox-message-div-right {
  margin: 3px 0;
}

.inbox-message-received {
  background-color: #dbdbdb;
  color: black;
  border: #323232 2px solid;
  float: left;
  max-width: 60%;
  height: auto;
  overflow-wrap: anywhere;
  padding: 5px;


  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;
}

.inbox-message-sent {
  background-color: #016ddb;
  color: white;
  border: #323232 2px solid;
  float: right;
  max-width: 60%;
  overflow-wrap: anywhere;
  padding: 5px;
  margin-right: 5px;

  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;
}

.message-time-received {
  text-align: -webkit-left;
  font: 12px bold;
  padding-top: 10px;
  padding-left: 5px;
}

.message-time-sent {
  float: right;
  text-align: right;
  font: 12px bold;
  padding-top: 10px;
  padding-right: 5px;
}

.picture-button {
  transform: scale(1.5);
  margin-top: 24px;
  margin-right: 10px;
  cursor: pointer;
  transition: 0.3s;
}

.picture-button:hover {
  transform: scale(1.6);
  transition: 0.3s;
  color: #016ddb;
}

</style>