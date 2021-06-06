<template>
  <div id="profile-info"
       style="background-color: transparent;">
    <v-img  class="profile-head"
            src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
            @click="$router.push('/user')"
            alt="Profile picture"/>
    <b style="cursor: pointer;"
       @click="$router.push('/user')">USERNAME</b>
  </div>
</template>

<script>
export default {
name: "UserProfileHead"
}
</script>

<style scoped>

  #profile-info {
    display: inline-block;
    text-align: -webkit-center;
    margin: 10px;
    height: 150px;

  }

  .profile-head {
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

  .profile-head:hover {
    cursor: pointer;
    transition: .3s;
    filter: brightness(0.7);
  }

</style>