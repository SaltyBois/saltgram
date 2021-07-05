<template>
  <div style="background-color: #efeeee; height: 100vh">
    <TopBar style="position: sticky; z-index: 2"/>
    <div id="main-div">
      <div class="post-or-story-div">
        <v-btn class="content-buttons"
               v-bind:class="{'accent': !this.btnPostClicked, 'primary': this.btnPostClicked}"
               @click="showForm = 'post'; btnPostClicked = true">
          New Post
        </v-btn>
        <v-btn class="content-buttons"
               v-bind:class="{'accent': btnPostClicked, 'primary': !btnPostClicked}"
               @click="showForm = 'story'; btnPostClicked = false">
          New Story
        </v-btn>
      </div>
      <div class="rest-area">

<!--        <transition name="fade">-->
          <PostForm v-if="showForm === 'post'" />
<!--        </transition>-->

<!--        <transition name="fade">-->
        <StoryForm v-if="showForm === 'story'" />
<!--        </transition>-->

      </div>

    </div>
  </div>
</template>

<script>
import TopBar from "@/components/TopBar";
import PostForm from "@/components/new_content_components/PostForm";
import StoryForm from "@/components/new_content_components/StoryForm";

export default {
  name: "NewContent",
  components: {TopBar, PostForm, StoryForm},
  data: function () {
    return {
      showForm: 'post',
      btnPostClicked: true
    }
  }
}
</script>

<style scoped>

#main-div {
  background-color: #FFFFFF;
  display: flex;
  flex-direction: row;
  margin: 10px 10%;
  height: 80vh;

  border: black 1px solid;

  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;
}

.post-or-story-div {
  display: flex;
  flex-direction: column;
  height: 20%;
  width: 200px;

  justify-content: center;
  align-content: center;
}

.content-buttons {
  margin: 5px;
}

.rest-area {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;

  justify-content: center;
  align-content: center;

  border-left: black 1px solid;
}


.fade-enter-active,
.fade-leave-active {
  transition: opacity 0s;
}

.fade-enter,
.fade-leave-to {
  opacity: 0;
}

</style>