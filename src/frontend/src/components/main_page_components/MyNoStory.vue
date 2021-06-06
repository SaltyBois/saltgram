<template>
  <div class="story-layout">
    <div @click="$router.push('/newContent'); ">
      <v-img  class="no-story"
              src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
              alt="Profile picture">
        <b class="story-plus-circle">
          <h3 class="story-plus">+</h3>
        </b>
      </v-img>
    </div>
    <b>Add story</b>
  </div>
</template>

<script>
export default {
  name: "MyNoStory",
  methods: {

  }
}
</script>

<style scoped>

  .story-layout {
    padding: 5px 10px;
    width: 150px;
    text-align: -webkit-center;
  }

  .no-story{
    width: 80px;
    height: 80px;
    object-fit: cover;
    border-radius: 20%;
    margin: 10px;
    cursor: pointer;
    border-style: solid;
    border-width: 2px;
    border-color: #000000;


    filter: brightness(1);

    transition: .3s;
    z-index: 0;
  }

  .no-story:hover {
    transition: .3s;
    filter: brightness(0.7);
  }

  .story-plus {
    position: relative;
    bottom: 10px;
    left: 0;
  }

  .story-plus-circle {
    position: relative;
    display: inline-block;
    top: 45px;
    left: 22px;
    border-radius: 50%;
    border-style: solid;
    border-width: 2px;
    border-color: #000000;
    width: 30px;
    height: 30px;
    background-color: #017ba9;
    color: white;

    font-size: 25px;
    text-align: -webkit-center;
  }
</style>