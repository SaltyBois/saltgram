<template>
  <div class="tag-search-div" @click="goToTag()">
    <div style="padding: 3%">
      <i class="fa fa-hashtag tag-search-head pt-1"/>
    </div>
    <div style="padding: 3%;overflow: hidden">
      <b>
        {{name}}
      </b>
    </div>
  </div>
</template>

<script>
export default {
name: "TagSearch",
props: {
  name: {type: String, required: true},
},
methods: {
  goToTag() {
    this.$router.push("/tag/" + this.name);
  }
},
}

</script>

<style scoped>

.tag-search-div {
  display: flex;
  flex-direction: row;
  height: 45px;
  width: 100%;
  background-color: transparent;
  cursor: pointer;
  transition: 0.3s;
}

.tag-search-div:hover {
  background-color: #949494;
  transition: 0.3s;
}

.tag-search-head {
  width: 30px;
  height: 30px;
  object-fit: cover;
  border-radius: 20%;
  cursor: pointer;

  border-style: solid;
  border-width: 2px;
  border-color: #323232;
  filter: brightness(1);

  transition: .3s;
  z-index: 0;
}

</style>