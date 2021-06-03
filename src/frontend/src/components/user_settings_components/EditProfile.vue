<template>
  <div class="main-div">
    <div class="item-container">
      <div class="profile-head-layout">
        <ImageMessage v-if="showContent" :image-src="this.profileImage" @toggle-image-message="showContent = false"/>
        <v-img  class="head"
                @click="showContent = true"
                :src="this.profileImage"
                alt="Profile picture"/>
        <b style="color: #2b80e0; margin-top: 5px; cursor:pointer;" @click="$refs.file.click()">Change profile photo</b>
      </div>
      <div style="padding-top: 15px; margin-left: 5px">
        <h1 style="text-align: left; justify-content: center" >Username</h1>
      </div>

      <input type="file"
             ref="file"
             style="display: none"
             @change="onSelectedFile($event)"
             accept="image/*,video/*">

    </div>
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px;">Name</h3>
        </div>
        <div style="width: 50%;">
          <v-text-field outlined value="Name and Lastname" style="width: 400px"/>
        </div>
      </div>
    </div>
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px; ">Username</h3>
        </div>
        <div style="width: 50%;">
          <v-text-field outlined value="Username" style="width: 400px"/>
        </div>
      </div>
    </div>
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px; ">Description</h3>
        </div>
        <div style="width: 50%;">
          <v-textarea outlined height="80px" no-resize value="Lorem ipsum" style="width: 400px;"/>
        </div>
      </div>
    </div>
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px;">E-mail</h3>
        </div>
        <div style="width: 50%;">
          <v-text-field outlined no-resize value="bezbednovic@gmail.com" style="width: 400px;"/>
        </div>
      </div>
    </div>
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px;">Phone number</h3>
        </div>
        <div style="width: 50%;">
          <v-text-field outlined no-resize value="+381 123 4567" style="width: 400px;"/>
        </div>
      </div>
    </div>
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px;">Gender</h3>
        </div>
        <div style="width: 50%;">
          <v-select outlined :items="genderRoles" v-model="genderRoles[0]" style="width: 400px"/>
        </div>
      </div>
    </div>
    <div class="item-container " style="height: auto;">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 40%;">
          <h3 style="margin-top: 14px;">Date of Birth</h3>
        </div>
        <div style="width: 60%;">
          <v-date-picker show-current />
        </div>
      </div>
    </div>
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px;">Web Site</h3>
        </div>
        <div style="width: 50%;">
          <v-text-field outlined no-resize value="Saltgram.rs" style="width: 400px;"/>
        </div>
      </div>
    </div>
    <div class="item-container ">
      <v-btn class="primary mb-5">Confirm changes</v-btn>
    </div>
  </div>
</template>

<script>
import ImageMessage from "@/components/inbox_components/ImageMessage";

export default {
  name: "EditProfile",
  components: {ImageMessage},
  data: function () {
    return {
      profileImage: 'https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg',
      showContent: false,
      genderRoles: [ 'Male', 'Female' ]
    }
  },
  methods: {
    onSelectedFile(event) {
      var files = event.target.files || event.dataTransfer.files;
      if (!files.length)
        return;
      console.log(files.length)
      console.log(files[0])
      this.item.image = URL.createObjectURL(files[0])
      console.log(this.item.image)
      if (files[0]['type'].includes('image')) this.typeContent = 'image';
      else this.typeContent = 'video';
      console.log(this.typeContent)
      this.isUploadedContent = true;
    },
  }
}
</script>

<style scoped>

.main-div {
  display: inline-flex;
  flex-direction: column;
  height: 100%;
  overflow-y: auto;
}

.item-container {
  height: 100px;
  display: inline-flex;
  flex-direction: row;
  justify-content: center;
  align-content: center;
  text-align: -webkit-center;
}

.profile-head-layout {
  margin-top: 5px;
  width: auto;
  align-content: center;
}

.head {
  width: 70px;
  height: 70px;
  object-fit: cover;
  border-radius: 20%;
  cursor: pointer;

  border-style: solid;
  border-width: 1px;
  border-color: black;

  filter: brightness(1);

  transition: .3s;
  z-index: 0;
}

</style>