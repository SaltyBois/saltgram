<template>
  <v-layout row
            style="justify-content: center; margin-top: 80px">
    <v-layout column
              align-center>
      <h4>Posts</h4>
      <h3><b>123</b></h3>
    </v-layout>
    <v-layout column
              align-center
              class="following-follower-div"
              @click="print; modalVisible=!modalVisible; title='Following';">
      <h4>Following</h4>
      <h3><b>1000</b></h3>
    </v-layout>
    <v-layout column
              align-center
              class="following-follower-div"
              @click="print; modalVisible=!modalVisible; title='Followers'; print">
      <h4>Followers</h4>
      <h3><b>10k</b></h3>
    </v-layout>
    <ModalListOfProfiles v-if="modalVisible" @toggle-visibility="toggleVisibility" :title="title"/>
  </v-layout>
</template>

<script>
import ModalListOfProfiles from "@/components/user_page_components/ModalListOfProfiles";

export default {
  name: "ProfileHeader",
  components: {ModalListOfProfiles},
  data: function () {
    return {
      modalVisible: false,

      title: ''
    }
  },
  methods: {
    toggleVisibility() {
      this.modalVisible = !this.modalVisible
    },
    print(){
      console.log('this.modalVisible:')
      console.log(this.modalVisible)
    }
  }
}
</script>

<style scoped>

.following-follower-div {
  background-color: transparent;
  transition: 0.3s;
  cursor: pointer;

  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;
}

.following-follower-div:hover {
  background-color: #858585;
  transition: 0.3s;
}

</style>