<template>
  <div class="layout-div">
    <div class="post-header-left-side">
      <v-img  class="post-header-profile"
              src="https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg"
              @click="$router.push('/user')"
              alt="Profile picture"/>
      <b @click="$router.push('/user')" style="cursor: pointer">Username1 </b>
      <p2 style="margin-left: 5px; margin-right: 5px"> has tagged you in comment in </p2>
      <b @click="$router.push('/user')" style="cursor: pointer">Username2 </b>
      <p2 style="margin-left: 5px"> post</p2>
      <v-img  class="post-header-profile"
              src="https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg"
              @click="$router.push('/user')"
              alt="Profile picture"/>
      <p4 style="color: #858585">1h ago</p4>
    </div>
  </div>
</template>

<script>
export default {
  name: "CommentTagNotification"
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

.post-header-right-side {
  direction: ltr;
  flex-direction: row;
  text-align: -webkit-center;
  align-items: center;
  float: right;
  display: flex;
  justify-content: center
}

</style>