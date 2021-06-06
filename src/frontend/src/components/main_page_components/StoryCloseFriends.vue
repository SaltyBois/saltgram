<template>
  <div class="story-layout">
    <v-img  class="story-close-friends"
            src="https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg"
            alt="Profile picture"/>
    <b>Close Friend Story</b>
  </div>
</template>

<script>
export default {
name: "StoryCloseFriends"
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