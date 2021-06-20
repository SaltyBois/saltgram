<template>
  <div style="background-color: #efeeee; height: 100vh">
    <portal-target name="drop-down-profile-menu" />
    <portal-target name="settings-menu"/>
    <TopBar style="position: sticky; z-index: 2"/>
    <div id="main-div"
         style="background-color: transparent">

      <!--        TODO: MILE-->
      <div id="notifications-header-div"
           style="background-color: white;">
        <h2 style="letter-spacing: 1px">Notifications</h2>
      </div>

      <div class="top-notification-bar mt-3" v-if="!privateProfile">
        <v-btn v-bind:class="NotificationCategory === 0 ? 'primary' : 'accent'"
               @click="NotificationCategory = 0"
               class="mx-2 my-1"
               small>
          Regular notifications
        </v-btn>
        <v-btn v-bind:class="NotificationCategory === 1 ? 'primary' : 'accent'"
               @click="NotificationCategory = 1"
               class="mx-2 my-1"
               small>
          Follow request
        </v-btn>
      </div>

      <div class="notifications-body-div" v-if="NotificationCategory === 0">

        <CommentTagNotification/>

        <FollowNotification/>

        <PostCommentNotification/>

        <PostLikeNotification/>

        <PostTagNotification v-for="index in 5" :key="index"/>

      </div>

      <div class="notifications-body-div" v-else-if="NotificationCategory === 1">

        <RequestProfile v-for="item in this.followingRequests" :key="item" :username-prop="item.username" v-on:reload-requests="getFollowRequests()"/>

      </div>
    </div>
  </div>
</template>

<script>
import TopBar from "@/components/TopBar";
import CommentTagNotification from "@/components/notifications_components/CommentTagNotification";
import FollowNotification from "@/components/notifications_components/FollowNotification";
import PostCommentNotification from "@/components/notifications_components/PostCommentNotification";
import PostLikeNotification from "@/components/notifications_components/PostLikeNotification";
import PostTagNotification from "@/components/notifications_components/PostTagNotification";
import RequestProfile from "@/components/notifications_components/RequestProfile";

export default {
  name: "Notifications",
  components: {TopBar, CommentTagNotification, FollowNotification, PostCommentNotification, PostLikeNotification,
               PostTagNotification, RequestProfile},
  data: function () {
    return {
      privateProfile: false,
      NotificationCategory: 0,
      followingRequests: [],
    }
  },
  methods: {
    getFollowRequests: function() {
      this.axios.get("users/follow/requests/", {headers: this.getAHeader()})
      .then(r => {
        this.followingRequests = r.data;
        console.log(r.data);
      }).catch(err => {
        console.log(err);
      })
    },
  },
  mounted() {
    this.getFollowRequests();
  }
}
</script>

<style scoped>

#main-div {
  display: inline-block;
  margin: 15px 20% 0 20%;
  width: 60vw;
  height: 85vh;
  flex-direction: row;
  justify-content: center;
  align-content: center;


}

#notifications-header-div {
  display: flex;
  height: auto;
  flex-direction: column;
  text-align: -webkit-center;
  border: black 1px solid;

  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;

}

.notifications-body-div {
  background-color: #FFFFFF;
  height: 65vh;
  overflow-x: hidden;
  overflow-y: scroll;
  margin-top: 15px;
  /*display: flex;*/
  /*height: auto;*/
  /*flex-direction: column;*/
  /*text-align: -webkit-center;*/

  border: black 2px solid;

  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;

}

.top-notification-bar {
  position: static;
  display: inline-flex;
  padding: 0 25%;
  width: 100%;
  height: 40px;
  background-color: #FFFFFF;
  text-align: -webkit-center;

  border: 1px solid black;
  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;

}

</style>