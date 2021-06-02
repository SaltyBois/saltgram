<template>
  <div class="story-layout">
    <v-img  class="story-seen"
            src="https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg"
            alt="Profile picture"/>
    <b>Muted Story</b>
  </div>
</template>

<script>
export default {
name: "StoryMuted"
}
</script>

<style scoped>

.story-layout {
  padding: 5px 10px;
  width: 150px;
  text-align: -webkit-center;
  filter: brightness(1.5);
  opacity: 0.5;

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

  transition: .3s;
  z-index: 0;
}


</style>