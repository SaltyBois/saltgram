<template>
  <div class="post-card">
    <PostView ref="postView"
              :user-prop="user"
              :post="postElement"
              media-path="https://www.arabianbusiness.com/public/styles/square/public/images/2021/03/28/meme.jpg?itok=DeJVUtab" />
    <div class="post-header">
      <div class="post-header-left-side">
        <v-img  class="post-header-profile"
                v-if="user.profilePictureURL"
                :src="user.profilePictureURL"
                @click="$router.push('/user/' + user.username)"
                alt="Profile picture"/>
        <v-img  class="post-header-profile"
                v-else
                :src="require('@/assets/profile_placeholder.png')"
                alt="Profile picture"
                @click="$router.push('/user/' + user.username)"/>
        <b @click="$router.push('/user/' + user.username)" style="cursor: pointer">{{ user.username }}</b>
        <div v-if="postElement.post.sharedMedia.media[0].location.name == ''" @click="$router.push('/location/' + postElement.post.sharedMedia.media[0].location.name)" style="margin-left:3px; cursor:pointer">&#183;   {{postElement.post.sharedMedia.media[0].location.name}}</div>
      </div>
      <div class="post-header-right-side">
        <b style="font-size: 25px; padding-bottom: 5px; cursor: pointer" @click="$refs.postInfo.$data.showDialog = true">...</b>
        <PostInfo :username="user.username" ref="postInfo"/>
      </div>
    </div>
    <div class="post-content">
      <v-carousel class="post-content-media" :continuous="false" >
        <v-carousel-item v-for="(item, index) in postElement.post.sharedMedia.media" :key="index">
          <v-img contain
                 width="100%"
                 height="80%"
                 v-if="item.filename.endsWith('.jpg') || item.filename.endsWith('.png') || item.filename.endsWith('.jpeg')"
                 :src="item.url"/>
          <video controls
                 width="100%"
                 height="70%"
                 loop
                 v-else-if="item.filename.endsWith('.mp4')"
                 :src="item.url"/>
        </v-carousel-item>
      </v-carousel>
    </div>
    <div class="post-interactions">
      <div class="post-interactions-left-side">
        <div style="width: 50px; height: 50px; text-align: -webkit-center">
          <i class="fa fa-thumbs-o-up"
             @click="like()"
             v-bind:class="userReactionStatus === 'LIKE' ? 'liked' : 'like'"
             aria-hidden="true"/>
        </div>
        <div style="width: 50px; height: 50px; text-align: -webkit-center;">
          <i class="fa fa-thumbs-o-up"
             @click="dislike()"
             v-bind:class="userReactionStatus === 'DISLIKE' ? 'disliked' : 'dislike'"
             aria-hidden="true"/>
        </div>
        <div style="width: 50px; height: 50px; text-align: -webkit-center">
          <i class="fa fa fa-comment-o like" aria-hidden="true" @click="showPostFun"/>
        </div>
      </div>
      <div class="post-interactions-right-side">
        <div style="width: 50px; height: 50px; text-align: -webkit-center">
          <i class="fa fa-folder-open-o like" aria-hidden="true"/>
        </div>
      </div>
    </div>
    <div class="post-description">
      <div style=" padding: 5px;">
        <p style="text-align: left; font-size: 12pt; margin-bottom: auto;">
          <b>{{likes}}</b> Likes  <b>{{dislikes}}</b> Dislikes
        <p style="text-align: left; font-size: 12pt; margin-bottom: auto; margin-top: 10px">
          <b>{{ user.username }}</b>
        {{ postElement.post.sharedMedia.media[0].description}}
        </p>
      </div>
       <div style="display:flex; overflow:auto; white-space:nowrap" class="my-2">
          <p class="tag" v-for="tag in postElement.post.sharedMedia.media[0].tags" :key="tag.value" @click="$router.push('/tag/'+tag.value)" >
            #{{tag.value}}
          </p>
       </div>
        <div style="display:flex; overflow:auto; white-space:nowrap" class="my-2">
          <div v-for="tagged in postElement.taggedUsers" :key="tagged.username">
            <div class="tagged-user" @click="$router.push('/user/'+tagged.username)">
              <v-img  class="post-header-profile"
                      v-if="tagged.profilePictureURL"
                      :src="tagged.profilePictureURL"
                      alt="Profile picture"/>
              <v-img  class="post-header-profile"
                      v-else
                      :src="require('@/assets/profile_placeholder.png')"
                      alt="Profile picture"/>
              <b style="cursor: pointer" class="mt-3">{{ tagged.username }}</b>
          </div>
          </div>
    </div>
    </div>
    <div class="post-comment-section">
      <div class="all-comments" >
        <div v-for="(item, index) in comments" :key="index">
          <CommentOnPostView  v-if="index < 2" :comment="item" style="min-height: 10px; height: auto; max-height: 20px;"/>
        </div>

        <p v-if="comments.length > 2" style="text-align: left; font-size: 10pt; margin-bottom: auto; margin-top: 20px; cursor: pointer" @click="showPostFun">
          View all <b>{{ comments.length }}</b> comments
        </p>
        <p style="text-align: left; font-size: 10pt; margin-bottom: auto; margin-top: 20px; color: #858585">
        {{ new Date(postElement.post.sharedMedia.media[0].addedOn.substring(0, lastIndex).replace('CEST', '(CEST)')).toLocaleString('sr') }}
        </p>
      </div>

    </div>
    <!--  TODO(Mile): Emojis need to be included GENERICALLY  -->
    <div class="post-footer">
      <div style="float: left; height: available; display: flex; flex-direction: row; width: 80%">
        <v-text-field label="Add a comment" v-model="input" style="width: available; padding: 5px" />
      </div>
      <div style="float: right; height: available; display: inline-block; ">
        <v-btn class="follow-button"
               @click="sendComment()"
               style="margin: 5px; width: 100px">
          comment
        </v-btn>
      </div>
    </div>
  </div>
</template>

<script>
import PostView from "@/components/PostView";
import CommentOnPostView from "@/components/CommentOnPostView";
import PostInfo from "@/components/PostInfo";


export default {
  name: "PostOnMainPage",
  components: {
    PostView, CommentOnPostView, PostInfo
  },
  props: {
    postElement: { type: Object, required: true},
    user: { type: Object, required: true}
  },
  methods: {
    showPostFun() {
      this.$refs.postView.$data.show = !this.$refs.postView.$data.show
    },
    getComments() {
      this.comments = [];
      this.axios.get("content/comment/" + this.postElement.post.id)
          .then(r => {
            // console.log(r);
            this.comments = r.data;
            // console.log(this.comments);
          }).catch(err => {
        console.log(err)
        this.$router.push('/');
      })
    },
    getReactions() {
      this.reactions = [];
      this.likes = 0;
      this.dislikes = 0;
      this.axios.get("content/reaction/" + this.postElement.post.id)
          .then(r => {
            // console.log(r);
            this.reactions = r.data;
            if(this.reactions === null){
              this.reactions = [];
            }
            // console.log(this.reactions);
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
            // console.log(r);
            this.loggedUser = r.data;
            // console.log(this.user.id);
            for(let i = 0; i < this.reactions.length; i++){
              // console.log(this.reactions[i].userId)
              if(this.reactions[i].userId == this.loggedUser.id){
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
    like() {
      let reaction = {reactionType: 'LIKE', postId: this.postElement.post.id};
      if(this.userReactionStatus === ''){
        this.axios.post("content/reaction", reaction, {headers: this.getAHeader()})
            .then(() => {
              // console.log(r);
              this.userReactionStatus = 'LIKE';
              this.getReactions();
            }).catch(err => {
          console.log(err)
          this.$router.push('/');
        })
      } else if (this.userReactionStatus === 'DISLIKE') {
        let putReaction = {reactionType: 'LIKE', id: this.reactionId}
        this.axios.put("content/reaction", putReaction, {headers: this.getAHeader()})
            .then(() => {
              // console.log(r);
              this.userReactionStatus = 'LIKE';
              this.getReactions();
            }).catch(err => {
          console.log(err)
          this.$router.push('/');
        })
      }
    },
    dislike() {
      let reaction = {reactionType: 'DISLIKE', postId: this.postElement.post.id};
      if(this.userReactionStatus === ''){
        this.axios.post("content/reaction", reaction, {headers: this.getAHeader()})
            .then(() => {
              // console.log(r);
              this.userReactionStatus = 'DISLIKE';
              this.getReactions();
            }).catch(err => {
          console.log(err)
          this.$router.push('/');
        })
      } else if (this.userReactionStatus === 'LIKE'){
        let putReaction = {reactionType: 'DISLIKE', id: this.reactionId}
        this.axios.put("content/reaction", putReaction, {headers: this.getAHeader()})
            .then(() => {
              // console.log(r);
              this.userReactionStatus = 'DISLIKE';
              this.getReactions();
            }).catch(err => {
          console.log(err)
          this.$router.push('/');
        })
      }
    },
    sendComment() {
      let com = {content: this.input, postId: this.postElement.post.id};
      this.axios.post("content/comment", com, {headers: this.getAHeader()})
          .then(() => {
            // console.log(r);
            this.input = '';
            this.getComments();
          }).catch(err => {
        console.log(err)
        this.$router.push('/');
      })
    },
  },
  mounted() {
    this.lastIndex = this.postElement.post.sharedMedia.media[0].addedOn.indexOf('CEST') + 4
    console.log(this.lastIndex)
    console.log(this.postElement.post.sharedMedia.media[0].addedOn.substring(0, this.lastIndex).replace('CEST', '(CEST)'))
    
    let datetime = this.postElement.post.sharedMedia.media[0].addedOn.split('.')[0];
    let date = datetime.split(' ')[0];
    let time = datetime.split(' ')[1];
    let dateParts = date.split('-');
    let timeParts = time.split(':');

    console.log(new Date(dateParts[0], dateParts[1] - 1, dateParts[2], timeParts[0], timeParts[1]));
    // console.log(this.postElement.post.sharedMedia.media[0].url)
    this.getComments();
    this.getReactions();
  },
  data: function () {
    return {
      input: '',
      showPost: false,
      comments: [],
      reactions: {},
      likes: 0,
      dislikes: 0,
      loggedUser: {},
      userReactionStatus: false,
      reactionId: false,
      lastIndex: 0,
    }
  }
}
</script>

<style scoped>

.post-card {
  margin: 10px 50px;
  background-color: white;
  width: available;
  height: auto;

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

.post-header-left-side, .post-header-right-side, .post-interactions-left-side, .post-interactions-right-side .post-footer-right-side {
  direction: ltr;
  flex-direction: row;
  text-align: -webkit-center;
  align-items: center;
  float: left;
  display: flex;
  justify-content: center
}

.post-header-right-side, .post-interactions-right-side, .post-footer-right-side {
  float: right;
  width: 50px;
  height: 50px;
}

.post-header-profile {
  width: 30px;
  height: 30px;
  object-fit: cover;
  border-radius: 10px;
  border: black solid 1px;
  margin: 10px;
  cursor: pointer;


  filter: brightness(1);

  transition: .3s;
  z-index: 0;
}

.post-content {
  width: 30vw;
  height: 30vw;

  align-content: center;

  object-fit: cover;
}

.post-interactions {
  height: 51px;

  border-top-style: solid;
  border-top-width: 1px;
  border-top-color: #373737;
}

.post-description {
  height: auto;
}

.post-comment-section {
  height: auto;
  padding: 5px;
}

.post-footer {
  height: 50px;
  /*display: flex;*/
  /*flex-direction: row;*/
}

.post-header-left-side, .post-header-right-side, .post-interactions-left-side, .post-interactions-right-side .post-footer-right-side {
  direction: ltr;
  flex-direction: row;
  text-align: -webkit-center;
  align-items: center;
  float: left;
  display: flex;
  justify-content: center
}

.post-header-right-side, .post-interactions-right-side, .post-footer-right-side {
  float: right;
  width: 50px;
  height: 50px;
}

.post-content-media {
  position: relative;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);

  max-width:100%;
  max-height:100%;

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

.all-comments {
  display: grid;
}

.follow-button, .unfollow-button  {
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

.left-btn, .right-btn {
  float: left;
  z-index: 3;
  margin-top: 25%;
  margin-left: 30px;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 1px black solid;
  background-color: black;
  color: #FFFFFF;
  text-align: center;
  justify-content: center;
  cursor: pointer;
  opacity: 0.3;
}

.left-btn:hover, .right-btn:hover {
  opacity: 0.9;
  transform: scale(1.3);
  transition: 0.3s;
}


.right-btn {
  float: right;
  margin-left: 0;
  margin-right: 30px;
  transform: rotate(0deg);
}

.top-right-album {
  float: right;
  margin-right: 3px;
  margin-top: 3px;
  background-color: black;
  color: white;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  padding-top: 8px;
  padding-right: 1px;
  opacity: 0.9;
  transform: rotate(0deg);
}

.top-right-album-letters {
  transition: slide 0.3s;

}

.fade-enter-active,
.fade-leave-active {
  transition: opacity .3s;
}

.fade-enter,
.fade-leave-to {
  opacity: 0;
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

.tag {
  border: rgb(15, 15, 202) solid 1px;
  border-radius: 5px;
  background: rgb(156, 240, 255);
  width: auto;
  transition: 0.3s;
  cursor: pointer;
  margin: 0 3px;
  padding: 0 3px;
}

.tag:hover{
  background: rgb(85, 228, 253);
    transition: 0.3s;
}

.tagged-user{
  border: rgb(0, 0, 0) solid 1px;
  border-radius: 5px;
  background: rgb(156, 240, 255);
  width: auto;
  transition: 0.3s;
  cursor: pointer;
  margin: 0 3px;
  padding: 0 3px;
  display: flex;
}

.tagged-user:hover{
  background: rgb(85, 228, 253);
    transition: 0.3s;
}

</style>