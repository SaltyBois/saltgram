<template>
  <div>
    <PostView ref="postView" :post="post" :key="reloadKey" @reload="reloadPostView"/>
    <video  class="post"
            v-if="post.post.sharedMedia.media[0].mimeType === 1"
            :controls="false"
            :playsinline="false"
            :preload="true"
            :autoplay="false"
            :src="post.post.sharedMedia.media[0].url"
            @click="showPostFun"/>
    <v-img  class="post"
            v-else
            :src="post.post.sharedMedia.media[0].url"
            @click="showPostFun"
            alt="Post"/>
  </div>
</template>

<script>
import PostView from "@/components/PostView";

export default {
  name: "PostOnUserPage",
  components: {PostView},
  data() {
    return {
      reloadKey: 0,
      //post: [],

    }

  },
  props: {
    post: { type: Object, required: true}
  },
  methods: {
    showPostFun() {
      this.$refs.postView.$data.show = !this.$refs.postView.$data.show
    },

    reloadPostView: function() {
      ++this.reloadKey;
      this.showPostFun();
    },
  },
  mounted() {
    // console.log(this.post.post.sharedMedia.media[0])
    //console.log(this.post);
  }
}
</script>

<style scoped>

  .post {
    width: 300px;
    height: 300px;
    object-fit: cover;
    border-radius: 20%;
    margin: 10px;
    cursor: pointer;

    border-style: solid;
    border-width: 2px;
    border-color: #323232;
    background-color: transparent;
    filter: brightness(1);

    transition: .3s;
    z-index: 0;
  }

  .post:hover {
    transition: .3s;
    filter: brightness(0.7);
  }

</style>