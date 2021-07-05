<template>
  <div class="story-layout">
    <StoryView ref="storyView" :stories="realStories"/>
    <v-img  class="story"
            @click="toggle()"
            v-if="user.profilePictureURL"
            :src="user.profilePictureURL"
            alt="Profile picture"/>
    <v-img  class="story"
            @click="toggle()"
            v-else
            :src="require('@/assets/profile_placeholder.png')"
            alt="Profile picture"/>
    <b>Your story</b>
  </div>
</template>

<script>
import StoryView from "@/components/StoryView";

export default {
  name: "MyStory",
  components: {StoryView},
  props: {
    user: { type: Object, required: true},
    stories: { type: Object, required: true},
  },
  data: function () {
    return {
      realStories: []
    }
  },
  methods: {
    toggle() {
      this.$refs.storyView.toggleView();
    }
  },
  mounted() {
    this.stories.storyElement.forEach(el => {
      el.closeFriends = this.stories.closeFriends
      this.realStories.push(el.stories[0])
    })
    // console.log(this.stories)
  }
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