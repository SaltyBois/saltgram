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
             accept="image/*">

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
    <div class="item-container mb-10">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <v-btn class="mx-2" v-bind:class="!privateProfile ? 'primary' : 'accent'" @click="privateProfile = false"><i class="fa fa-unlock mr-1"/>Public profile</v-btn>
        </div>
        <div style="width: 50%;">
          <v-btn class="mx-2" v-bind:class="privateProfile ? 'primary' : 'accent'" @click="privateProfile = true"><i class="fa fa-lock mr-1"/>Private profile</v-btn>
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
      profilePicture: 'https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg',
      showContent: false,
      genderRoles: [ 'Male', 'Female' ],
      privateProfile: false,
      isUploadedContent: false,
    }
  },
  methods: {
    onSelectedFile(event) {

      this.profilePicture = event.target.files[0]
      console.log(this.profilePicture)

      this.refreshToken(this.getAHeader())
        .then(rr => {
          this.$store.state.jws = rr.data

          let data = new FormData();
          data.append('profileImg', this.profilePicture);
          let config = {
            headers: {
              'Content-Type': 'multipart/form-data',
              'Authorization': 'Bearer ' + this.$store.state.jws,
            },
          };
          // Vratiti nazad
          this.axios.post("content/profilepicture", data, config)
            .then(() => this.isUploadedContent = true)
            .catch(r => console.log(r));
        }).catch(() => this.$router.push('/'));


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