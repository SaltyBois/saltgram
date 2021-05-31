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

      <div class="notifications-body-div">
        <CommentTagNotification/>

        <FollowNotification/>

        <PostCommentNotification/>

        <PostLikeNotification/>

        <PostTagNotification v-for="index in 5" :key="index"/>

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

export default {
  name: "Notifications",
  components: {TopBar, CommentTagNotification, FollowNotification, PostCommentNotification, PostLikeNotification,
               PostTagNotification},
  methods: {

  },
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
  height: 70vh;
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

</style>