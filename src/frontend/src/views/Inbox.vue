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
            <div style="padding: 3%">
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
        <div class="inbox-content-body">

        </div>
        <div class="inbox-content-footer">

        </div>
      </div>
    </div>
  </div>
</template>

<script>
import TopBar from "@/components/TopBar";
export default {
name: "Inbox",
  components: {TopBar}
}
</script>

<style scoped>

#main-inbox-div {
  background-color: white;
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

.profile-chat-head {
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

</style>