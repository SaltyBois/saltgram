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

          <ViewedChat/>

          <UnviewedChat/>

        </div>
      </div>
      <div id="inbox-content">
        <div class="inbox-content-header">

          <ChatProfile/>

          <div style="height: 60px; text-align: right">
            <i class="fa fa-info info-button" />
          </div>

        </div>
        <div class="inbox-content-body" ref="chatBox">

          <ReceivedMessage/>

          <SentMessage/>

          <ReceivedImageMessage/>

          <SentImageMessage/>

        </div>
        <div class="inbox-content-footer">
          <div class="inner-inbox-content-footer">

            <ChatInput/>

          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import TopBar from "@/components/TopBar";
import EmojiPicker from "vue-emoji-picker";
import ImageMessage from "@/components/inbox_components/ImageMessage";
import ViewedChat from "@/components/inbox_components/ViewedChat";
import UnviewedChat from "@/components/inbox_components/UnviewedChat";
import ChatProfile from "@/components/inbox_components/ChatProfile";
import ReceivedMessage from "@/components/inbox_components/ReceivedMessage";
import SentMessage from "@/components/inbox_components/SentMessage";
import ReceivedImageMessage from "@/components/inbox_components/ReceivedImageMessage";
import SentImageMessage from "@/components/inbox_components/SentImageMessage";
import ChatInput from "@/components/inbox_components/ChatInput";

export default {
name: "Inbox",
  components: {TopBar, EmojiPicker, ImageMessage, ViewedChat, UnviewedChat, ChatProfile, ReceivedMessage, SentMessage,
               ReceivedImageMessage, SentImageMessage, ChatInput},
  methods: {
  },
  mounted() {
    var content = this.$refs.chatBox;
    content.scrollTop = content.scrollHeight;
  },
  data: function () {
    return {

    }
  }
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

.inbox-content-header {
  height: 60px;
  width: 100%;

  float: top;

  display: block;
  vertical-align: center;
  border-bottom: black 1px solid ;
}

.inbox-content-body {
  overflow-y: auto;
  overflow-x: hidden;

  height: available;
  width: 100%;
}

.inbox-content-footer {
  height: 60px;
  width: 100%;

  border-top: black 1px solid;
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

</style>