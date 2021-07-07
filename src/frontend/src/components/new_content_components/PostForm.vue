<template>
  <div class="main-div">
    <div class="post-form-body">
      <div class="pt-3 pl-3 post-form-body-left-side">
        <h2>Add posts</h2>
        <div class="thumbnails">
          <div class="thumbnail"
          v-for="img in imageUrls" :key="img.url">

            <video  class="post"
                    v-if="img.filename.includes('.mp4') || img.filename.includes('.mkv')"
                    :controls="false"
                    :playsinline="false"
                    :muted="true"
                    :preload="true"
                    :autoplay="true"
                    width="128px"
                    height="128px"
                    :src="img.url"/>
            <v-img v-else
                   :src="img.url"
                   height="128px"
                   max-width="128px"/>
            <v-btn
            fab
            absolute
            depressed
            x-small
            @click="removeImg(img)"
            class="remove-btn">X</v-btn>
          </div>
          <v-file-input
          v-model="images"
          hide-input
          multiple
          prepend-icon="fa-plus-circle"
          @change="refreshPreview"
          style="display:none"
          ref="fileInput"
          />  
          <div class="add-thumbnail" @click="selectFiles">
            +
          </div>
        </div>
        <div
        v-if="campaign"
        id="campaign-info">
          <v-select
          label="Age group"
          :items="ageGroups"
          v-model="ageGroup"></v-select>
          <v-text-field v-model="website"
          label="Website"
          />
          <v-checkbox v-model="oneTime"
          label="One time"></v-checkbox>
          <v-menu
          v-model="dateMenu"
          :close-on-content-click="false"
          :nudge-right="40"
          transition="scale-transition"
          offset-y
          min-width="auto"
          >
              <template v-slot:activator="{ on, attrs }">
              <v-text-field
                  label="Campaign start"
                  v-model="campaignStart"
                  prepend-icon="fa-calendar"
                  readonly
                  v-bind="attrs"
                  v-on="on"
                  :max="maxDate"
              ></v-text-field>
              </template>
              <v-date-picker
              v-model="campaignStart"
              @input="dateMenu = false"
              ></v-date-picker>
          </v-menu>
          <v-menu
          v-if="!oneTime"
          v-model="dateMenu2"
          :close-on-content-click="false"
          :nudge-right="40"
          transition="scale-transition"
          offset-y
          min-width="auto"
          >
              <template v-slot:activator="{ on, attrs }">
              <v-text-field
                  label="Campaign end"
                  v-model="campaignEnd"
                  prepend-icon="fa-calendar"
                  readonly
                  v-bind="attrs"
                  v-on="on"
                  :max="maxDate"
              ></v-text-field>
              </template>
              <v-date-picker
              v-model="campaignEnd"
              @input="dateMenu = false"
              ></v-date-picker>
          </v-menu>
        </div>
      </div>
      <div class="post-form-body-right-side">
        <h2>Upload info</h2>
        <v-textarea
        v-model="description"
        no-resize
        outlined
        label="Add a description"
        dense
        hide-details
        style="flex: 0 1 auto"/>
        <!--<v-text-field
        
        outlined
        label="TODO Add location"
        dense
        hide-details
        style="flex: 0 1 auto"/>-->
        <geosearch @selected="selectLocation"></geosearch>
        <div v-if="selectedLocation">
          <v-text-field 
          v-model="location.name"
          label="Name"
          disabled/>
        </div>
        <v-combobox
        v-model="tags"
        chips
        clearable
        label="Tags"
        multiple
        append-icon=""
        solo>
          <template v-slot:selection="{ attrs, item }">
              <v-chip
              v-bind="attrs"
              close
              @click:close="removeTag(item)">
                {{item}}
              </v-chip>
          </template>
        </v-combobox>
        <v-spacer></v-spacer>
        <v-text-field v-model="query" label="Search profiles" @focus="filter()" @keyup="filter()"/>
        <v-list style="height:200px">
        <v-list-item v-for="p in filteredProfiles" :key="p.username">
          <template>
              <div class="post-header-left-side">
                <v-img  class="post-header-profile"
                        v-if="p.profilePictureURL"
                        :src="p.profilePictureURL"
                        alt="Profile picture"/>
                <v-img  class="post-header-profile"
                        v-else
                        :src="require('@/assets/profile_placeholder.png')"
                        alt="Profile picture"/>
                <b @click="$router.push('/user')" style="cursor: pointer">{{ p.username }}</b>
              </div>
              <div class="post-header-right-side">
                <v-checkbox v-model="p.checked" @change="updateFinalList(p.userId)">
                  
                </v-checkbox>
              </div>
          </template>
        </v-list-item>
        </v-list>
        <div class="d-flex flex-row">
          <v-checkbox
          v-model="closeFriends"
          label="Close friends"/>
          <v-checkbox
          v-if="role == 'agent'"
          v-model="campaign"
          label="Campaign"></v-checkbox>
        </div>
        <v-spacer></v-spacer>
        <v-btn color="accent" :disabled="!images.length" @click="uploadFiles" :loading="uploading">Upload</v-btn>
      </div>
    </div>
  </div>
</template>

<script>

export default {
  name: "PostForm",
  data: function () {
    return {
      website: '',
      campaignStart: '',
      campaignEnd: '',
      oneTime: false,
      dateMenu: false,
      dateMenu2: false,
      ageGroups: ['Pre 20s', '20s', '30s'],
      ageGroup: '',
      campaign: false,
      role: 'user',
      uploading: false,
      description: "",
      location: {
				country: "",
				state:   "",
				zipCode: "",
        city: "",
				street:  "",
        name: "",
      },
      tags: [],
      images: [],
      imageUrls: [],
      taggableProfiles: [],
      taggedProfiles: [],
      filteredProfiles: [],
      finalList: [],
      query: '',

      selectedLocation: false,
    }
  },
  mounted() {
          this.getRole();
          this.axios.get('users/taggableprofiles/get', {headers: this.getAHeader()})
            .then(r => {
                console.log(r.data);
                this.taggableProfiles = r.data;
            })
            .catch(r => console.log(r));
  },
  methods: {
    getRole: function() {
      this.refreshToken(this.getAHeader())
        .then(rr => {
          this.$store.state.jws = rr.data;
          this.axios.get('users/get/role', {headers: this.getAHeader()})
            .then(r => {
              this.role = r.data;
            }).catch(() => this.$router.push('/'));
        });
    },

    filter() {
      if(this.query === ''){
        this.filteredProfiles = this.taggableProfiles;
        return;
      }

      this.filteredProfiles = [];
      this.taggableProfiles.forEach(el => {
        if (el.username.includes(this.query)){
          this.filteredProfiles.push(el);
        } 
      });
    },

    updateFinalList(userId) {
      if(this.finalList.includes(userId)){
        const index = this.finalList.indexOf(userId);
        this.finalList.splice(index, 1);
      } else {
        this.finalList.push(userId);
      }
      console.log(this.finalList);
    },
    selectLocation: function(l) {
      this.selectedLocation = true;
      this.location.country = l.address.country_code;
      this.location.state = l.address.country;
      this.location.zipCode = l.address.postcode;
      this.location.city = l.address.city;
      this.location.street = l.address.road;
      this.location.name = l.display_name.split(',')[0];
    },

    uploadFiles: function() {
      this.uploading = true;
      let data = new FormData();
      this.images.forEach(img => {
        data.append('posts', img);
      });
      this.tags.forEach(tag => {
        data.append('tags', tag);
      });
      this.finalList.forEach(userTag => {
        data.append('userTags', userTag);
      });
      data.append('description', this.description)
      data.append('location', JSON.stringify(this.location))
      data.append('campaign', this.campaign)
      data.append('ageGroup', this.ageGroup)
      data.append('oneTime', this.oneTime)
      data.append('campaignStart', this.campaignStart)
      data.append('campaignEnd', this.campaignEnd)
      data.append('website', this.website)
      this.refreshToken(this.getAHeader())
        .then(rr => {
          this.$store.state.jws = rr.data;
          let config = {
            headers: {
              'Content-Type': 'multipart/form-data',
              'Authorization': 'Bearer ' + this.$store.state.jws,
            },
          };
          this.axios.post('content/post', data, config)
            .then(() => {
              this.uploading = false
              this.$router.push('/main')
            })
            .catch(r => console.log(r));
        }).catch(() => this.$router.push('/'));
    },

    selectFiles: function() {
      this.$refs.fileInput.$refs.input.click();
    },

    removeImg: function(img) {
      let index = this.imageUrls.indexOf(img);
      this.imageUrls.splice(index, 1);
      this.imageUrls = [...this.imageUrls];

      this.images.splice(index, 1);
      this.images = [...this.images];
    },

    refreshPreview: function(files) {
      // this.imageUrls = [];
      files = files.slice(0, 10);
      this.images.push(files)
      this.images = this.images.slice(0, 10);
      this.images.forEach(f => {
        let img = {
          url: URL.createObjectURL(f),
          filename: f.name
        }
        console.log(img)
        this.imageUrls.push(img);
      });
    },
    select: function(r) {
      this.selected = r;
      this.results = [];
      this.$emit('selected', r);
    },
  }
}
</script>

<style scoped>

.main-div {
  width: 100%;
  height: 100%;
}

.post-form-header {
  height: 10%;
  width: 100%;
  text-align: center;
}

.post-form-body {
  display: flex;
  height: 100%;
  width: 100%;
}

.post-form-body-left-side {
  text-align: -webkit-center;
  width: 60%;
  height: 100%;
}

.post-form-body-right-side {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  width: 40%;
  margin: 10px;
}
.post-form-body-right-side > * {
  padding-top: 10px;
}

.post-header-left-side, .post-header-right-side, .post-interactions-left-side, .post-interactions-right-side {
  direction: ltr;
  flex-direction: row;
  text-align: -webkit-center;
  align-items: center;
  float: left;
  display: flex;
  justify-content: center
}

.post-header-right-side, .post-interactions-right-side {
  float: right;
  width: 50px;
  height: 50px;
}

.content-shape {
  /*display: flex;*/
  width: min-content;
  height: min-content;
  min-width: 100px;
  min-height: 100px;
  max-width: 60vh;
  max-height: 60vh;

  object-fit: contain;

  text-align: -webkit-center;
  justify-content: center;

  margin-top: 10px;
  border: 1px black solid;

  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;
}

.post-header-profile {
  width: 30px;
  height: 30px;
  object-fit: cover;
  border-radius: 10px;
  border: black solid 1px;
  margin: 10px;
  cursor: pointer;


  filter: brightness(1);

  transition: .3s;
  z-index: 0;

}

.content-item {
  width: 100%;
  height: 100%;

  max-width: 60vh;
  max-height: 60vh;
  display: block;
  text-align: -webkit-center;
  justify-content: center;

  object-fit: scale-down;


  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;

  cursor: pointer;
}

.no-content {
  position: relative;
  top: 10%;
  left: 0;
  transform: scale(2.5);
}

.thumbnail {
  position: relative;
  display: inline-block;
  margin: 10px;
}

.add-thumbnail {
  cursor: pointer;
  display: grid;
  place-items: center;
  font-weight: 500;
  font-size: 3rem;
  width: 128px;
  height: 128px;
  background: #eee;
  margin: 10px;
}

.remove-btn {
  top: 0;
  right: 0;
}

.thumbnails {
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  align-items: center;
  flex-wrap: wrap;
}
#results {
        position: absolute;
        top: 100%;
        left: 0;
        right: 0;
        z-index: 99;
        border: solid 1px #eee;
    }



    .result {
        cursor: pointer;
        border: solid 1px #eee;
    }

    .result:hover {
        background: #eee;
    }

</style>