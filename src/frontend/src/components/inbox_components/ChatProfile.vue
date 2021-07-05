<template>
  <div style="float: left; text-align: -webkit-center; height: 60px; display: flex">
    <v-img  class="profile-chat-head"
            src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
            @click="$router.push('/user')"
            alt="Profile picture"/>
    <h3 style="margin-top: 15px; margin-right: 10px">Username template</h3>
  </div>
</template>

<script>
export default {
  name: "ChatProfile"
}
</script>

<style scoped>

  .profile-chat-head{
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

</style>