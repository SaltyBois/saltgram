<template>
  <div class="story-layout">
    <v-img  class="story"
            src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
            alt="Profile picture"/>
    <b>Your story</b>
  </div>
</template>

<script>
export default {
name: "MyStory"
}
</script>

<style scoped>

  .story-layout {
    padding: 5px 10px;
    width: 150px;
    text-align: -webkit-center;
  }

  .story {
    width: 80px;
    height: 80px;
    object-fit: cover;
    border-radius: 20%;
    margin: 10px;
    cursor: pointer;

    border-style: solid;
    border-width: 5px;
    border-color: red;

    filter: brightness(1);

    transition: .3s;
    z-index: 0;
  }

  .story:hover {
    transition: .3s;
    filter: brightness(0.7);
  }

</style>