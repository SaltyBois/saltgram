<template>
  <div class="story-highlight-layout">
    <video  class="story-highlight"
            v-if="stories && stories[0].mimeType === 'video'"
            :controls="false"
            :playsinline="false"
            :preload="true"
            :autoplay="false"
            :src="stories[0].url"
            @click="$refs.storyView.$data.visible = true"/>
    <v-img v-else-if="stories && stories[0].mimeType === 'image'"
           class="story-highlight"
           @click="$refs.storyView.$data.visible = true"
           :src="stories[0].url"
           alt="Profile picture"/>
    <v-img v-else class="story-highlight" src="require('@/assets/profile_placeholder.png')" />
    <h5>{{name}}</h5>
    <StoryView :stories="stories" ref="storyView"/>
  </div>
</template>

<script>
import StoryView from "@/components/StoryView";
export default {
  name: "StoryHighlight",
  data: function () {
    return {
      
    }
  },
  props: {
    stories: {
      type: Array,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
  },
  components: {StoryView}
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