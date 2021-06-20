<template>
  <div>
    <PostView ref="postView" :post="post"/>
    <v-img  class="post"
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
      //post: [],
    }

  },
  props: {
    post: { type: Object, required: true}
  },
  methods: {
    showPostFun() {
      this.$refs.postView.$data.show = !this.$refs.postView.$data.show
    }
  },
  mounted() {
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