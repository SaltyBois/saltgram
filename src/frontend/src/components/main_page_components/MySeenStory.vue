<template>
  <div class="story-layout">
    <v-img  class="story-seen"
            src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
            alt="Profile picture"/>
    <b>Your story</b>
  </div>
</template>

<script>
export default {
  name: "MySeenStory"
}
</script>

<style scoped>

  .story-layout {
    padding: 5px 10px;
    width: 150px;
    text-align: -webkit-center;
  }

  .story-seen{
    width: 80px;
    height: 80px;
    object-fit: cover;
    border-radius: 20%;
    margin: 10px;
    cursor: pointer;

    border-style: solid;
    border-width: 3px;
    border-color: #535353;

    /*filter: brightness(1.5);*/
    opacity: 0.8;

    transition: .3s;
    z-index: 0;
  }

</style>