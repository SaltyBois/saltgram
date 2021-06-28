<template>
  <div class="main-panel">
    <div class="my-2">
      <v-btn small
             v-bind:class="category === 'profiles' ? 'primary' : ''"
             @click="category = 'profiles'"><i class="fa fa-address-book-o mr-1" />Profiles</v-btn>
      <v-divider vertical class="mx-1"/>
      <v-btn small
             v-bind:class="category === 'tags' ? 'primary' : ''"
             @click="category = 'tags'"><i class="fa fa-hashtag mr-1" />Tags</v-btn>
      <v-divider vertical class="mx-1"/>
      <v-btn small
             v-bind:class="category === 'locations' ? 'primary' : ''"
             @click="category = 'locations'"><i class="fa fa-map-marker mr-1" />Locations</v-btn>
    </div>
    <v-divider />
    <div v-if="processing">
      <v-progress-circular
          indeterminate
          class="mt-5"
          color="primary"/>
    </div>
    <div v-if="!processing" class="query-panel">
      <transition name="fade" appear>
        <div class="search-panel" v-if="category === 'profiles'" >
          <div v-for="(item, index) in searchedData.profiles" :key="index">
            <div>
              <ProfileSearch :username-prop="item.username" :profile-picture-address-prop="item.profilePictureURL" />
            </div>
            <v-divider/>
          </div>
        </div>
      </transition>
      <transition name="fade" appear>
        <div class="search-panel" v-if="category === 'tags'" >
          <div v-for="(item, index) in searchedData.tags" :key="index">
            <div>
              <TagSearch :name="item"/>
            </div>
            <v-divider/>
          </div>
        </div>
      </transition>
      <transition name="fade" appear>
        <div class="search-panel" v-if="category === 'locations'">
          <div v-for="index in 5" :key="index">
            <div>
              <LocationSearch />
            </div>
            <v-divider/>
          </div>
        </div>
      </transition>
    </div>
  </div>
</template>

<script>
import ProfileSearch from "@/components/topbar_components/ProfileSearch";
import TagSearch from "@/components/topbar_components/TagSearch";
import LocationSearch from "@/components/topbar_components/LocationSearch";

export default {
  name: "SearchPanel",
  components: { ProfileSearch, TagSearch, LocationSearch },
  data: function () {
    return {
      category: 'profiles',
      processing: false,
      searchedData: {
        profiles: [],
        tags: [],
        locations: []
      }
    }
  },
  methods: {

  },
  mounted() {
    console.log(this.searchedData.profiles)
  }

}
</script>

<style scoped>

.main-panel {
  background-color: #FFFFFF;
  text-align: -webkit-center;

  border: 1px solid black;
  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;
}

.query-panel {
  overflow-y: auto;
  /*height: 100%;*/
  /*overflow-clip: true;*/
}

.search-panel {
  transition: 0.3s;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform .5s;
}

.slide-enter,
.slide-leave-to {
  transform: translateY(-50%) translateX(100vw);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity .5s;
}

.fade-enter,
.fade-leave-to {
  opacity: 0;
}

</style>