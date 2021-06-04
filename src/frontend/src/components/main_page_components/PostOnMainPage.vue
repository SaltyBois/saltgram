<template>
  <div class="post-card">
    <PostView ref="postView" media-path="https://www.arabianbusiness.com/public/styles/square/public/images/2021/03/28/meme.jpg?itok=DeJVUtab" />
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

    <transition name="fade" appear>
      <div class="left-btn" v-if="iteratorContent > 0" @click="decrease()">
        <i class="fa fa-sign-out ml-2 mt-3" style="transform: scale(1.4) rotate(180deg);"/>
      </div>
    </transition>
    <div class="top-right-album" v-if="contentPlaceHolder.length !== 1" >
      <b class="top-right-album-letters">{{iteratorContent + 1}}/{{contentPlaceHolder.length}}</b>
    </div>
    <transition name="fade" appear>
      <div class="right-btn" v-if="iteratorContent < contentPlaceHolder.length - 1" @click="increase()">
        <i class="fa fa-sign-out mt-2 ml-1" style="transform: scale(1.4)" />
      </div>
    </transition>


    <div class="post-content">

      <v-img  v-for="(item, index) in contentPlaceHolder.length"
              :key="index"
              class="post-content-media"
              v-bind:style="index === iteratorContent ? '' : 'display: none'"
              :src="contentPlaceHolder[index]"
              alt="Post content"/>
<!--      <div class="post-content">-->
<!--        <v-img  class="post-content-media"-->
<!--                src="https://skinnyms.com/wp-content/uploads/2015/04/9-Best-Grumpy-Cat-Memes-750x500.jpg"-->
<!--                alt="Post content"/>-->
<!--      </div>-->

    </div>
    <div class="post-interactions">
      <div class="post-interactions-left-side">
        <div style="width: 50px; height: 50px; text-align: -webkit-center">
          <i class="fa fa-thumbs-o-up like" aria-hidden="true"/>
        </div>
        <div style="width: 50px; height: 50px; text-align: -webkit-center;">
          <i class="fa fa-thumbs-o-up dislike" aria-hidden="true"/>
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
          <b>1234</b> Likes  <b>532</b> Dislikes
        <p style="text-align: left; font-size: 10pt; margin-bottom: auto;">
          <b>USERNAME </b>
          Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.
        </p>
      </div>
    </div>
    <div class="post-comment-section">
      <div class="all-comments" >

        <CommentOnPostView v-for="index in 2" :key="index" style="min-height: 10px; height: auto; max-height: 20px;"/>
        <p style="text-align: left; font-size: 10pt; margin-bottom: auto; cursor: pointer" @click="showPostFun">
          View all <b>32</b> comments
        </p>
        <p style="text-align: left; font-size: 10pt; margin-bottom: auto; color: #858585">
          Posted 1 hour ago
        </p>
      </div>

    </div>
    <!--  TODO(Mile): Emojis need to be included GENERICALLY  -->
    <div class="post-footer">
      <div style="float: left; height: available; display: flex; flex-direction: row; width: 80%">
        <v-text-field label="Add a comment" style="width: available; padding: 5px" />
      </div>
      <div style="float: right; height: available; display: inline-block; ">
        <v-btn class="follow-button" style="margin: 5px; width: 75px">
          post
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

  },
  methods: {
    insert(emoji) {
      this.input += emoji
    },
    append(emoji) {
      this.input += emoji
    },
    showPostFun() {
      this.$refs.postView.$data.show = !this.$refs.postView.$data.show
    },
    decrease() {
      if (this.iteratorContent > 0) this.iteratorContent -= 1;
    },
    increase() {
      if (this.iteratorContent + 1 < this.contentPlaceHolder.length) this.iteratorContent += 1;
    }
  },
  mounted() {
    this.iteratorContent = 0
  },
  data: function () {
    return {
      input: '',
      search: '',
      showPost: false,
      iteratorContent: 0,
      contentPlaceHolder: [
        'https://skinnyms.com/wp-content/uploads/2015/04/9-Best-Grumpy-Cat-Memes-750x500.jpg',
        'https://i.kym-cdn.com/entries/icons/original/000/035/692/cover1.jpg',
        'https://www.thehonestkitchen.com/blog/wp-content/uploads/2019/07/CatMemes-copy-10.jpg',
        'https://i.ytimg.com/vi/KHa4OOvYLx0/maxresdefault.jpg'
      ]
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
  border-radius: 20%;
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

</style>