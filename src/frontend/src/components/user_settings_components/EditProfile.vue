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
        <h1 style="text-align: left; justify-content: center" >{{ user.username }}</h1>
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
          <v-text-field outlined v-model="profile.fullName" style="width: 400px"/>
        </div>
      </div>
    </div>
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px; ">Username</h3>
        </div>
        <div style="width: 50%;">
          <v-text-field outlined v-model="profile.username" style="width: 400px"/>
        </div>
      </div>
    </div>
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px; ">Description</h3>
        </div>
        <div style="width: 50%;">
          <v-textarea outlined height="80px" no-resize v-model="profile.description" style="width: 400px;"/>
        </div>
      </div>
    </div>
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px;">E-mail</h3>
        </div>
        <div style="width: 50%;">
          <v-text-field outlined no-resize v-model="profile.email" style="width: 400px;"/>
        </div>
      </div>
    </div>
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px;">Phone number</h3>
        </div>
        <div style="width: 50%;">
          <v-text-field outlined no-resize v-model="profile.phoneNumber" style="width: 400px;"/>
        </div>
      </div>
    </div>
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px;">Gender</h3>
        </div>
        <div style="width: 50%;">
          <v-select outlined :items="genderRoles" v-model="profile.gender" style="width: 400px"/>
        </div>
      </div>
    </div>
    <div class="item-container " style="height: auto;">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 40%;">
          <h3 style="margin-top: 14px;">Date of Birth</h3>
        </div>
        <div style="width: 60%;">
          <v-date-picker show-current v-model="profile.dateOfBirth" :max="maxDate" />
        </div>
      </div>
    </div>
    <div class="item-container ">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <h3 style="margin-top: 14px;">Web Site</h3>
        </div>
        <div style="width: 50%;">
          <v-text-field outlined no-resize v-model="profile.webSite" style="width: 400px;"/>
        </div>
      </div>
    </div>
    <div class="item-container mb-10">
      <div style="display: inline-flex; flex-direction: row; margin-top: 20px; width: 70%">
        <div style="width: 50%;">
          <v-btn class="mx-2" v-bind:class="!profile.privateUser ? 'primary' : 'accent'" @click="profile.privateUser  = false"><i class="fa fa-unlock mr-1"/>Public profile</v-btn>
        </div>
        <div style="width: 50%;">
          <v-btn class="mx-2" v-bind:class="profile.privateUser  ? 'primary' : 'accent'" @click="profile.privateUser  = true"><i class="fa fa-lock mr-1"/>Private profile</v-btn>
        </div>
      </div>
    </div>
    <div class="item-container ">
      <v-btn class="primary mb-5" @click="changeUserInfo">Confirm changes</v-btn>
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
      genderRoles: [ 'Male', 'Female' ],
      profile : {
        privateUser: true,
        description: '',
        fullName: '',
        username: '',
        webSite: '',
        phoneNumber: '',
        gender: '',
        dateOfBirth: '',
        email: ''
      },
      user: '',
      maxDate: ''
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
    getUserInfo: function() {
      this.refreshToken(this.getAHeader())
          .then(rr => {
            this.$store.state.jws = rr.data;
            this.axios.get("users", {headers: this.getAHeader()})
                .then(r =>{
                  this.user = r.data
                  this.profile.email = this.user.email
                  this.getProfileInfo();
                });

          }).catch(() => this.$router.push('/'));
    },
    getProfileInfo: function() {
      this.axios.get("users/profile/" + this.$route.params.username, {headers: this.getAHeader()})
          .then(r => {
            console.log(r.data)
            this.profile.privateUser = !r.data.isPublic;
            this.profile.username = r.data.username;
            this.profile.fullName = r.data.fullName;
            this.profile.description = r.data.description;
            this.profile.webSite = r.data.webSite;
            this.profile.phoneNumber = r.data.phoneNumber;
            this.profile.gender = r.data.gender;
            this.profile.dateOfBirth = r.data.dateOfBirth * 1000;
            // console.log(this.profile.dateOfBirth);
            this.dateFunc();
          }).catch(err => {
        console.log(err)
        console.log('Pushing Back to Login Page after fetching profile')
        this.$router.push('/');
      })
    },
    dateFunc() {
      let dateFull = new Date(this.profile.dateOfBirth);
      var formatedDate = dateFull.getFullYear() + '-';

      let month = dateFull.getMonth() + 1
      if (month < 10) {
        formatedDate += '0' + month + '-'
      } else {
        formatedDate += month + '-'
      }

      let date = dateFull.getDate()

      if (date < 10) {
        formatedDate += '0' + date
      } else {
        formatedDate += date
      }
      // console.log(formatedDate)
      this.profile.dateOfBirth = formatedDate;
    },
    maxDateFunc() {
      let now = new Date();
      this.maxDate = now.getFullYear() + '-';

      let month = now.getMonth() + 1
      if (month < 10) {
        this.maxDate += '0' + month + '-'
      }
      else {
        this.maxDate += month + '-'
      }

      let date = now.getDate()

      if (date < 10) {
        this.maxDate += '0' + date
      }
      else {
        this.maxDate += date
      }
      // console.log(this.maxDate)
    },
    changeUserInfo() {
      console.log(this.profile.dateOfBirth)
      let parts = this.profile.dateOfBirth.split('-')
      this.profile.dateOfBirth = new Date(parts[0], parts[1] - 1, parts[2])
      console.log(this.profile.dateOfBirth)
      let profileData = {
            privateProfile: this.profile.privateUser,
            description: this.profile.description,
            fullName: this.profile.fullName,
            username: this.profile.username,
            webSite: this.profile.webSite,
            phoneNumber: this.profile.phoneNumber,
            gender: this.profile.gender,
            dateOfBirth: this.profile.dateOfBirth,
            email: this.profile.email
        }
      this.axios.put("users/profile/" + this.$route.params.username, profileData, {headers: this.getAHeader()} )
      .then(r => {
        console.log(r)
      })
      .catch(() => this.$router.push('/'))
    }
  },
  mounted() {
    this.maxDateFunc();
    this.getUserInfo();
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