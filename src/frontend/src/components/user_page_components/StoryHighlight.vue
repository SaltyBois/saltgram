<template>
  <div class="story-highlight-layout">
    <v-img  class="story-highlight"
            src="https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg"
            alt="Profile picture"/>
    <h5>Highlights 1</h5>
  </div>
</template>

<script>
export default {
  name: "StoryHighlight"
}
</script>

<style scoped>

.story-highlight {
  width: 80px;
  height: 80px;
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

.story-highlight:hover {
  transition: .3s;
  filter: brightness(0.7);
}

.story-highlight-layout {
  padding: 5px 10px;
  width: 150px;
  text-align: -webkit-center;
}

</style>