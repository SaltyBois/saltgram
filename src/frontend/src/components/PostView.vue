<template>
  <div :key="renderKey">
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
              <div v-if="isUserLoggedIn" class="post-header-right-side">
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
                  <CommentOnPostView v-for="(item, index) in comments" :key="index" :comment="item"/>
                </div>
                <!--  TODO(Mile): Emojis need to be included GENERICALLY  -->
                <div class="post-footer">
                  <div v-if="isUserLoggedIn" class="post-interactions"
                       style="background-color: transparent">
                    <div class="post-interactions-left-side">
                      <div style="width: 50px; height: 50px; text-align: -webkit-center">
                        <i class="fa fa-thumbs-o-up " 
                          @click="like()" 
                          aria-hidden="true" 
                          v-bind:class="userReactionStatus === 'LIKE' ? 'liked' : 'like'"/>
                      </div>
                      <div style="width: 50px; height: 50px; text-align: -webkit-center;">
                        <i class="fa fa-thumbs-o-up " 
                          @click="dislike()" 
                          aria-hidden="true"
                          v-bind:class="userReactionStatus === 'DISLIKE' ? 'disliked' : 'dislike'"/>
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
                      <b>{{likes}}</b> Likes  <b>{{dislikes}}</b> Dislikes
                    </p>
                    <p style="text-align: left; font-size: 10pt; margin-bottom: auto; color: #858585">
                      Posted 1 hour ago
                    </p>
                  </div>                                    
                    <div v-if="isUserLoggedIn" style="float: left; height: available; display: flex; flex-direction: row; width: 80%; margin-bottom: 10px">
                    <v-text-field v-model="commentContent" label="Add a comment" style="width: available; padding: 5px" />
                    <v-btn class="post-button" style="margin: 5px; width: 75px" @click="sendComment()">
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
  computed: {
    isUserLoggedIn: function() {
      return this.$store.state.jws;
    }
  },
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
    loadingComments() {
           this.axios.get("content/comment/" + this.post.post.id)
           .then(r => {
              console.log(r);
              this.comments = r.data;
              console.log(this.comments);
            }).catch(err => {
              console.log(err)
              this.$router.push('/');
            })
    },
    sendComment() {
      let com = {content: this.commentContent, postId: this.post.post.id};
       this.axios.post("content/comment", com, {headers: this.getAHeader()})
           .then(r => {
              console.log(r);
              this.commentContent = '';
              this.renderKey++;
            }).catch(err => {
              console.log(err)
              this.$router.push('/');
            })
    },
    like() {
      let reaction = {reactionType: 'LIKE', postId: this.post.post.id};
      if(this.userReactionStatus === ''){
        this.axios.post("content/reaction", reaction, {headers: this.getAHeader()})
            .then(r => {
                console.log(r);
                this.userReactionStatus = 'LIKE';
              this.renderKey++;
              }).catch(err => {
                console.log(err)
                this.$router.push('/');
              })
      } else if (this.userReactionStatus === 'DISLIKE') {
        let putReaction = {reactionType: 'LIKE', id: this.reactionId}
        this.axios.put("content/reaction", putReaction, {headers: this.getAHeader()})
           .then(r => {
              console.log(r);
              this.userReactionStatus = 'LIKE';
              this.renderKey++;
            }).catch(err => {
              console.log(err)
              this.$router.push('/');
            })
      }
    },
    dislike() {
      let reaction = {reactionType: 'DISLIKE', postId: this.post.post.id};
      if(this.userReactionStatus === ''){
       this.axios.post("content/reaction", reaction, {headers: this.getAHeader()})
           .then(r => {
              console.log(r);
              this.userReactionStatus = 'DISLIKE';
              this.renderKey++;
            }).catch(err => {
              console.log(err)
              this.$router.push('/');
            })
      } else if (this.userReactionStatus === 'LIKE'){
        let putReaction = {reactionType: 'DISLIKE', id: this.reactionId}
        this.axios.put("content/reaction", putReaction, {headers: this.getAHeader()})
           .then(r => {
              console.log(r);
              this.userReactionStatus = 'DISLIKE';
              this.renderKey++;
            }).catch(err => {
              console.log(err)
              this.$router.push('/');
            })
      }
    },
    loadingReactions() {
           this.axios.get("content/reaction/" + this.post.post.id)
           .then(r => {
              console.log(r);
              this.reactions = r.data;
              if(this.reactions === null){
                this.reactions = [];
              }
              console.log(this.reactions);
              for(let i = 0; i < this.reactions.length; i++){
                if(this.reactions[i].reactionType === 'LIKE'){
                  this.likes += 1;
                } else {
                  this.dislikes += 1;
                }
              }
              if(this.$store.state.jws){
                this.checkIfReacted();
              }
            }).catch(err => {
              console.log(err)
              this.$router.push('/');
            })
    },
    checkIfReacted() {
      this.axios.get("users", {headers: this.getAHeader()})
           .then(r => {
              console.log(r);
              this.user = r.data;
              console.log(this.user.id);
              for(let i = 0; i < this.reactions.length; i++){
                console.log(this.reactions[i].userId)
                if(this.reactions[i].userId == this.user.id){
                  this.userReactionStatus = this.reactions[i].reactionType;
                  this.reactionId = this.reactions[i].id; 
                  break;
                }
              }
            }).catch(err => {
              console.log(err)
              this.$router.push('/');
            })

    },
  },
  mounted() {
    this.iteratorContent = 0
    this.loadingPost();
    this.loadingComments();
    this.loadingReactions();
  },
  data: function () {
    return {
      show: false,
      search: '',
      iteratorContent: 0,
      contentPlaceHolder: [],
      description: '',
      comments: [],
      commentContent: '',
      reactions: [],
      user: '',
      userReactionStatus: '',
      likes: 0,
      dislikes: 0,
      reactionId: '',
      renderKey: 0,
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

.liked, .disliked {
  position: relative;
  top: 12px;
  left: 0;
  transform: scale(2);
  margin: 0 3px;
  transition: 0.2s;
  color: #016ddb;
  cursor: pointer;
}

.disliked {
    color: #ff0000;
    transform: scale(2) rotate(180deg);
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