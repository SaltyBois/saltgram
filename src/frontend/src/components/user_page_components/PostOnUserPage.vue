<template>
  <div>
    <PostView ref="postView" media-path="https://i.kym-cdn.com/entries/icons/original/000/032/100/cover4.jpg"/>
    <v-img  class="post"
            src="https://i.kym-cdn.com/entries/icons/original/000/032/100/cover4.jpg"
            @click="showPostFun"
            alt="Post"/>
  </div>
</template>

<script>
import PostView from "@/components/PostView";

export default {
  name: "PostOnUserPage",
  components: {PostView},
  methods: {
    showPostFun() {
      this.$refs.postView.$data.show = !this.$refs.postView.$data.show
    }
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