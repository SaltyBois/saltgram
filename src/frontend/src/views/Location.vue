<template>
  <div>
    <div class="reactions-main">
      <TopBar style="position: sticky; z-index: 2"/>
      <div id="reactions-header-div">
        <h1 style="letter-spacing: 1px">Location {{$route.params.name}}</h1>
      </div>
      <v-layout class="user-media"
                column>
        <PostOnUserPage v-for="(item, index) in content" :key="index" :post="item" :user="item.user"/>
      </v-layout>
    </div>
  </div>
</template>

<script>
import TopBar from "@/components/TopBar";
import PostOnUserPage from "@/components/user_page_components/PostOnUserPage";

export default {
  name: "Location",
  components: { TopBar, PostOnUserPage },
  data() {
    return {
      content: [],
      user: {},
    }
  },
  mounted() {
    this.loadLocationContent();
  },
  methods: {
    loadLocationContent() {
      this.refreshToken(this.getAHeader())
            .then(rr => {
              this.$store.state.jws = rr.data;
                this.axios.get('content/location/' + this.$route.params.name, {headers: this.getAHeader()})
                .then(r => {
                  this.content = r.data
                  console.log(this.content)
                })
                .catch(r => console.log(r));
        })
    }
  }
}
</script>

<style scoped>

.reactions-main {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-content: center;
  /* text-align: center; */
  background: #efeeee;
  min-height: 100vh;
  height: auto;
  /*margin-left: 10%;*/
  /*margin-right: 10%;*/
}

.user-media {
  --w:400px;
  --n:3;
  --m:2;

  margin: 5px 10%;
  display:grid;
  grid-template-columns:repeat(auto-fit,minmax(clamp(100%/(var(--n) + 1) + 0.1%,(var(--w) - 100vw)*1000,100%/(var(--m) + 1) + 0.1%),1fr)); /*this */
  gap:10px;
}

#reactions-header-div {
  align-self: center;
  display: flex;
  height: auto;
  flex-direction: column;
  text-align: -webkit-center;
  /*border: black 1px solid;*/

  width: 400px;
  /*border-radius: 10px;*/
  -top: 10px;
}

</style>