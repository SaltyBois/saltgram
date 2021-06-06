<template>
  <div class="inbox-message-div-left">
    <div style="float: left; text-align: -webkit-center; height: 60px; display: flex">
      <v-img  class="profile-chat-head-message"
              src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
              @click="$router.push('/user')"
              alt="Profile picture"/>
    </div>
    <div class="inbox-message-received">
      Received message template
    </div>
    <div class="message-time-received">
      15.4.2015. 17:45 time template
    </div>
  </div>
</template>

<script>
export default {
  name: "ReceivedMessage"
}
</script>

<style scoped>

  .inbox-message-div-left {
    display: inline-block;
    height: auto;
    width: 100%;
    align-content: center;
    margin: 3px 0;
  }

  .profile-chat-head-message {
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

  .message-time-received {
    text-align: -webkit-left;
    font: 12px bold;
    padding-top: 10px;
    padding-left: 5px;
  }

</style>