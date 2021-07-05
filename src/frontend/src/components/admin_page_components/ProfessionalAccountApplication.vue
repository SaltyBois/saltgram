<template>
  <div class="profile-search-div">
    <div style="padding: 3%">
      <v-img  class="profile-search-head"
              v-if="profilePictureAddressProp"
              :src="profilePictureAddressProp"
              alt="Profile picture"/>
      <v-img  class="profile-search-head"
              v-else
              :src="require('@/assets/profile_placeholder.png')"
              alt="Profile picture"/>
    </div>
    <div style="padding: 3%;overflow: hidden">
      <b>
        {{usernameProp}}
      </b>
    </div>
  </div>
</template>

<script>
export default {
  name: "ProfessionalAccountApplication",
  props: {
    usernameProp: {
      type: String,
      required: true
    },
    profilePictureAddressProp: {
      type: String,
      required: true
    }
  }
}
</script>

<style scoped>

.profile-search-div {
  display: flex;
  flex-direction: row;
  height: 50px;
  width: 100%;
  background-color: transparent;
  cursor: pointer;
  transition: 0.3s;
}

.profile-search-div:hover {
  background-color: #949494;
  transition: 0.3s;
}

.profile-search-head {
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