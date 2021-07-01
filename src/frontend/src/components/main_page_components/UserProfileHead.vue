<template>
  <div id="profile-info"
       style="background-color: transparent;">
    <v-img  class="profile-head"
            :src="user.profilePictureURL"
            @click="$router.push('/user/' + user.username)"
            alt="Profile picture"/>
    <b style="cursor: pointer;"
       @click="$router.push('/user/' + user.username)">{{ user.username }}</b>
  </div>
</template>

<script>
export default {
  name: "UserProfileHead",
  props: {
    user: { type: Object, required: true}
  }
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