<template>
  <div>
    <transition name="fade" appear>
      <v-layout class="modal-overlay"
                align-center
                v-if="show"
                @click="toggleParrent"
                justify-center/>
    </transition>
    <transition name="fade" appear>
      <v-layout class="modal"
                justify-center
                align-center
                v-if="show"
                column>
        <!--    TODO(MILE): POST-->
          <div class="post-card">
            <div class="post-header">
              <div class="post-header-left-side">
                <v-img  class="post-header-profile"
                        src="https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg"
                        @click="$router.push('/user')"
                        alt="Profile picture"/>
                <b @click="$router.push('/user')" style="cursor: pointer">Username1</b>
              </div>
              <div class="post-header-right-side">
                <b style="font-size: 25px; padding-bottom: 5px; cursor: pointer" @click="$refs.postInfo.$data.showDialog = true">...</b>
                <PostInfo username="Username1" ref="postInfo"/>
              </div>
            </div>
            <div style="display: flex; justify-content: space-between;">
              <v-layout class="post-content" align-center justify-center style="display: block; object-fit: contain">
                <v-carousel class="post-content-media" :continuous="false">
                  <v-carousel-item v-for="(item, index) in contentPlaceHolder.length" :key="index">
                    <v-img contain
                           :src="contentPlaceHolder[index]"/>
                   <!-- <video controls
                           loop
                           v-else-if="contentPlaceHolder[index].endsWith('.mp4')"
                           :src="contentPlaceHolder[index]"/>-->
                  </v-carousel-item>
                </v-carousel>
              </v-layout>
              <v-layout style="display: flex; flex-direction: column;" >
                <div class="post-description">
                  <div style=" padding: 5px;">
                    <p style="text-align: left; font-size: 10pt; margin-bottom: auto;">
                      <b>USERNAME</b>
                      {{this.description}}
                    </p>
                  </div>
                </div>
                <div class="post-comment-section">
                  <CommentOnPostView v-for="index in 10" :key="index"/>

                </div>
                <!--  TODO(Mile): Emojis need to be included GENERICALLY  -->
                <div class="post-footer">
                  <div class="post-interactions"
                       style="background-color: transparent">
                    <div class="post-interactions-left-side">
                      <div style="width: 50px; height: 50px; text-align: -webkit-center">
                        <i class="fa fa-thumbs-o-up like" aria-hidden="true"/>
                      </div>
                      <div style="width: 50px; height: 50px; text-align: -webkit-center;">
                        <i class="fa fa-thumbs-o-up dislike" aria-hidden="true"/>
                      </div>
                    </div>
                    <div class="post-interactions-right-side">
                      <div style="width: 50px; height: 50px; text-align: -webkit-center">
                        <i class="fa fa-folder-open-o like" aria-hidden="true"/>
                      </div>
                    </div>
                  </div>

                  <div style=" padding: 5px;">
                    <p style="text-align: left; font-size: 12pt; margin-bottom: auto;">
                      <b>1234</b> Likes  <b>532</b> Dislikes
                    </p>
                    <p style="text-align: left; font-size: 10pt; margin-bottom: auto; color: #858585">
                      Posted 1 hour ago
                    </p>
                  </div>

                  <div style="float: left; height: available; display: flex; flex-direction: row; width: 80%; margin-bottom: 10px">
                    <v-text-field label="Add a comment" style="width: available; padding: 5px" />
                    <v-btn class="post-button" style="margin: 5px; width: 75px">
                      post
                    </v-btn>
                  </div>
                </div>
              </v-layout>
            </div>
          </div>
      </v-layout>
    </transition>
  </div>
</template>

<script>
import CommentOnPostView from "@/components/CommentOnPostView";
import PostInfo from "@/components/PostInfo";

export default {
  name: "PostView",
  components: { CommentOnPostView, PostInfo },
  props: {
      mediaPath: {
        required: false,
        type: String
      },
      post: { type: Object, required: true}
  },
  methods: {
    toggleParrent() {
      this.show = !this.show
    },
    decrease() {
      if (this.iteratorContent > 0) this.iteratorContent -= 1;
      console.log(this.iteratorContent)
    },
    increase() {
      if (this.iteratorContent + 1 < this.contentPlaceHolder.length) this.iteratorContent += 1;
      console.log(this.iteratorContent)
    },
    loadingPost() {
      for(let i = 0; i < this.post.post.sharedMedia.media.length; i++){
        this.contentPlaceHolder.push(this.post.post.sharedMedia.media[i].url);
      }
      console.log(this.contentPlaceHolder);
      this.description = this.post.post.sharedMedia.media[0].description;
    },
  },
  mounted() {
    this.iteratorContent = 0
    this.loadingPost();
  },
  data: function () {
    return {
      show: false,
      search: '',
      iteratorContent: 0,
      contentPlaceHolder: [],
      description: '',
    }
  }
}
</script>

<style scoped>

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 98;
  background-color: rgba(0, 0, 0, 0.3);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity .5s;
}

.fade-enter,
.fade-leave-to {
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform .5s;
}

.slide-enter,
.slide-leave-to {
  transform: translateY(-50%) translateX(100vw);
}

.post-card {
  background-color: white;
  width: 98%;
  height: auto;
  max-height: 100%;

  border: #323232 solid 1px;

  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;
}

.post-header {
  height: 51px;

  border-bottom-style: solid;
  border-bottom-width: 1px;
  border-bottom-color: #373737;
}

.post-header-left-side, .post-header-right-side, .post-interactions-left-side, .post-interactions-right-side {
  direction: ltr;
  flex-direction: row;
  text-align: -webkit-center;
  align-items: center;
  float: left;
  display: flex;
  justify-content: center
}

.post-header-right-side, .post-interactions-right-side {
  float: right;
  width: 50px;
  height: 50px;
}

.post-content-media {
  justify-self: center;
  margin-top: 2.5%;
  max-height: 95%;
  max-width: 95%;
  min-height: available;
}

.like, .dislike {
  position: relative;
  top: 12px;
  left: 0;
  transform: scale(2);
  margin: 0 3px;

  transition: 0.2s;
}

.dislike {
  transform: scale(2) rotate(180deg);
}

.like:hover {
  transition: 0.2s;
  color: #016ddb;
  cursor: pointer;
}

.dislike:hover {
  transition: 0.2s;
  color: #ff0000;
  cursor: pointer;
}

.post-content {
  width: 30vw;
  height: 30vw;

  /*object-fit: cover;*/

  border-right: black 1px solid;
}

.post-interactions {
  height: 51px;

}

.post-description {
  height: auto;
}

.post-comment-section {
  max-height: 30vh;
  height: 60%;
  max-width: 60vw;
  padding: 5px;
  overflow-x: hidden;
  overflow-y: auto;
  scroll-behavior: smooth;
}

.post-footer {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: auto;
  min-height: 40%;
  border-top: black 1px solid;
}

.post-header-profile {
  width: 30px;
  height: 30px;
  object-fit: cover;
  border-radius: 20%;
  margin: 10px;
  cursor: pointer;


  filter: brightness(1);

  transition: .3s;
  z-index: 0;

}

.post-button {
  margin: 10px 0;
  width: 100px;
  height: 50px;
  background-color: transparent;
  color: #016ddb;
  border-color: #016ddb;
  border-style: solid;
  border-width: 1px;
  text-align: -webkit-center;
}

.modal {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 99;

  width: 75%;
  height: 75%;
  background-color: #FFF;
  border-radius: 16px;

  padding: 5px;
}

</style>