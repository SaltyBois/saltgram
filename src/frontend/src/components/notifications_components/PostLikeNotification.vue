<template>
  <div class="layout-div">
    <div class="post-header-left-side">
       <v-img  v-if="this.pictureProp"
                class="post-header-profile"
                :src="this.pictureProp"
                @click="$router.push('/user/' + usernameProp)"
                alt="Profile picture"/>
         <v-img v-else class="post-header-profile"
          @click="$router.push('/user/' + usernameProp)"
          :src="require('@/assets/profile_placeholder.png')"/>
      <b @click="$router.push('/user/' + usernameProp)" style="cursor: pointer">@{{usernameProp}}</b>
      <p style="margin-left: 5px; margin-right: 5px; margin-top: 15px">liked your post </p>      
    </div>
  </div>
</template>

<script>
export default {
  name: "PostLikeNotification",
  props: {
    usernameProp: {
      type: String,
      required: true
    },
    pictureProp: {
      type: String,
      required: true,
    }
  },
}
</script>

<style scoped>

.layout-div {
  height: 70px;
  width: 100%;

  display: flex;
  flex-direction: row;
  flex-flow: wrap;
  transition: 0.3s;
  cursor: pointer;
}

.layout-div:hover {
  background-color: #e0e0e0;
  transition: 0.3s;
}

.post-header-left-side {
  direction: ltr;
  flex-direction: row;
  text-align: -webkit-center;
  align-items: center;
  float: left;
  display: flex;
  justify-content: center
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

</style>