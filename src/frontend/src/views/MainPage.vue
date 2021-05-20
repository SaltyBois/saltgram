<template>
  <div style="background-color: #efeeee;">
    <portal-target name="drop-down-profile-menu" />
    <portal-target name="settings-menu"/>
    <TopBar style="position: sticky; z-index: 2"/>
    <div id="main-div"
         style="background-color: transparent">

      <div id="media-div"
           style="background-color: transparent;">
        <div id="stories-div"
             style="background-color: transparent;">
          <v-layout id="user-stories"
                    style="background-color: transparent;"
                    column>
            <v-layout class="inner-story-layout"
                      style=";background-color: transparent;">
              <div class="story-layout">
                <div>
                  <v-img  class="no-story"
                          src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
                          alt="Profile picture">
                    <b class="story-plus-circle">
                      <h3 class="story-plus">+</h3>
                    </b>
                  </v-img>
                </div>
                <b>Your story</b>
              </div>
              <div class="story-layout">
                <v-img  class="story-close-friends"
                        src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
                        alt="Profile picture"/>
                <b>Your story</b>
              </div>
              <div class="story-layout">
                <v-img  class="story-seen"
                        src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
                        alt="Profile picture"/>
                <b>Your story</b>
              </div>
              <div class="story-layout">
                <v-img  class="story-close-friends"
                        src="https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg"
                        alt="Profile picture"/>
                <b>Username1</b>
              </div>
              <div class="story-layout">
                <v-img  class="story"
                        src="https://i.pinimg.com/736x/4d/8e/cc/4d8ecc6967b4a3d475be5c4d881c4d9c.jpg"
                        alt="Profile picture"/>
                <b>Username1</b>
              </div>
            </v-layout>
          </v-layout >
        </div>
        <div id="posts-div"
             style="background-color: mediumslateblue;">

        </div>
      </div>
      <div id="suggestions-div"
           style="background-color: transparent;">
        <div id="profile-info"
             style="background-color: transparent;">
          <v-img  class="profile-head"
                  src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
                  @click="$router.push('/user')"
                  alt="Profile picture"/>
          <b style="cursor: pointer;"
             @click="$router.push('/user')">USERNAME</b>
        </div>
        <div id="suggestions-info"
             style="background-color: transparent">
          <h4>Suggestions For You</h4>
          <div>
            <div class="suggested-profile">
              <v-img  class="suggestion-head"
                      src="https://i.pinimg.com/564x/4e/c4/f2/4ec4f2d69c9bc6b152abcb420252c3a8.jpg"
                      @click="$router.push('/user')"
                      alt="Profile picture"/>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <h5>@USERNAME</h5>
                <h5>Ime i prezime</h5>
                <small>45 mutual followers</small>
              </div>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <v-btn depressed
                       class="follow-button">
                  follow
                </v-btn>
              </div>
            </div>
            <div class="suggested-profile">
              <v-img  class="suggestion-head"
                      src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRkO7BR9U3GRatwXeEFkxxDRm7a5Kq7eGhJtg&usqp=CAU"
                      @click="$router.push('/user')"
                      alt="Profile picture"/>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <h5>@USERNAME</h5>
                <h5>Ime i prezime</h5>
                <small>45 mutual followers</small>
              </div>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <v-btn depressed
                       class="unfollow-button">
                  unfollow
                </v-btn>
              </div>
            </div>
            <div class="suggested-profile">
              <v-img  class="suggestion-head"
                      src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRYYR7aw1k6gkP_UGKGddjPee4ENRCEY8AO3A&usqp=CAU"
                      @click="$router.push('/user')"
                      alt="Profile picture"/>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <h5>@USERNAME</h5>
                <h5>Ime i prezime</h5>
                <small>45 mutual followers</small>
              </div>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <v-btn depressed
                       class="follow-button">
                  follow
                </v-btn>
              </div>
            </div>
            <div class="suggested-profile">
              <v-img  class="suggestion-head"
                      src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRYYR7aw1k6gkP_UGKGddjPee4ENRCEY8AO3A&usqp=CAU"
                      @click="$router.push('/user')"
                      alt="Profile picture"/>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <h5>@USERNAME</h5>
                <h5>Ime i prezime</h5>
                <small>45 mutual followers</small>
              </div>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <v-btn depressed
                       class="follow-button">
                  follow
                </v-btn>
              </div>
            </div>
            <div class="suggested-profile">
              <v-img  class="suggestion-head"
                      src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRYYR7aw1k6gkP_UGKGddjPee4ENRCEY8AO3A&usqp=CAU"
                      @click="$router.push('/user')"
                      alt="Profile picture"/>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <h5>@USERNAME</h5>
                <h5>Ime i prezime</h5>
                <small>45 mutual followers</small>
              </div>
              <div style="margin: 0 3px; text-align: -webkit-center">
                <v-btn depressed
                       class="follow-button">
                  follow
                </v-btn>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import TopBar from "@/components/TopBar";

export default {
  name: "MainPage",
  components: {
    TopBar,
  },
}
</script>

<style scoped>
  #main-div {
    display: flex;
    margin: 0 10%;
    height: 1000px;
    flex-direction: row;
  }

  #media-div {
    display: flex;
    flex-direction: column;

    width: 70%;
  }

  #stories-div {
    display: flex;
    height: 150px;
    flex-direction: row;
  }

  #posts-div {
    display: flex;
    height: 500px;
    flex-direction: column;
  }

  #suggestions-div {
    display: flex;
    flex-direction: column;
    width: 30%;
    height: border-box;
  }

  #profile-info {
    display: inline-block;
    text-align: -webkit-center;
    margin: 10px;
    height: 150px;

  }

  #suggestions-info {
    display: inline-block;
    text-align: -webkit-left;
    margin: 10px;
    width: available;
  }

  .profile-head, .suggestion-head {
    width: 100px;
    height: 100px;
    object-fit: cover;
    border-radius: 20%;
    margin: 10px;
    cursor: pointer;

    border-style: solid;
    border-width: 2px;
    border-color: #323232;
    filter: brightness(1);

    transition: .3s;
    z-index: 0;
  }

  .profile-head:hover, .suggestion-head:hover {
    cursor: pointer;
    transition: .3s;
    filter: brightness(0.7);
  }

  .suggestion-head {
    width: 60px;
    height: 60px;
    margin: 0;
  }

  .story-layout {
    padding: 5px 10px;
    width: 150px;
    text-align: -webkit-center;
  }

  .inner-story-layout {
    height: 150px;
    flex-direction: row;
    overflow-x: auto;
    overflow-y: hidden;
    white-space: nowrap;

  }

  .story, .story-close-friends, .no-story, .story-seen {
    width: 80px;
    height: 80px;
    object-fit: cover;
    border-radius: 20%;
    margin: 10px;
    cursor: pointer;


    filter: brightness(1);

    transition: .3s;
    z-index: 0;
  }

  .story:hover, .story-close-friends:hover, .no-story:hover {
    transition: .3s;
    filter: brightness(0.7);
  }

  .story {
    border-style: solid;
    border-width: 5px;
    border-color: red;
  }

  .story-close-friends {
    border-style: solid;
    border-width: 5px;
    border-color: #36c400;
  }

  .no-story {
    border-style: solid;
    border-width: 2px;
    border-color: #000000;
  }

  .story-plus-circle {
    position: relative;
    display: inline-block;
    top: 45px;
    left: 22px;
    border-radius: 50%;
    border-style: solid;
    border-width: 2px;
    border-color: #000000;
    width: 30px;
    height: 30px;
    background-color: #017ba9;
    color: white;

    font-size: 25px;
    text-align: -webkit-center;
  }

  .story-plus {
    position: relative;
    bottom: 10px;
    left: 0px;
  }

  .story-seen {
    border-style: solid;
    border-width: 3px;
    border-color: #535353;

    filter: brightness(1.5);
    opacity: 0.8;
  }

  .suggested-profile {
    display: -webkit-inline-flex;
    flex-direction: row;
    height: 60px;
    background-color: transparent;

  }

  .follow-button, .unfollow-button  {
    margin: 10px 0;
    width: 100px;
    height: 50px;
    background-color: transparent;
    color: #016ddb;
    border-color: #016ddb;
    border-style: solid;
    border-width: 1px;
    text-align: -webkit-center;
  }

  .unfollow-button {
    color: #ff2626;
    border-color: #ff2626;
  }

  .suggestion-info {
    display: flex;
    flex-direction: row;
  }

  .post-div {

  }


</style>