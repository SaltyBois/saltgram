<template>
  <div class="story-layout">
    <v-img  class="story-close-friends"
            v-if="user.profilePictureURL"
            :src="user.profilePictureURL"
            alt="Profile picture"/>
    <v-img  class="story-close-friends"
            v-else
            :src="require('@/assets/profile_placeholder.png')"
            alt="Profile picture"/>
    <b>{{ user.username }}</b>
  </div>
</template>

<script>
export default {
  name: "StoryCloseFriends",
  props: {
    user: { type: Object, required: true},
    stories: { type: Array, required: true}
  }
}
</script>

<style scoped>

  .story-layout {
    padding: 5px 10px;
    width: 150px;
    text-align: -webkit-center;
  }

  .story-close-friends{
    width: 80px;
    height: 80px;
    object-fit: cover;
    border-radius: 20%;
    margin: 10px;
    cursor: pointer;
    border-style: solid;
    border-width: 5px;
    border-color: #36c400;

    filter: brightness(1);

    transition: .3s;
    z-index: 0;
  }

  .story-close-friends:hover{
    transition: .3s;
    filter: brightness(0.7);
  }

</style>