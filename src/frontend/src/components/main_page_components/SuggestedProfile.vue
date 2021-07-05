<template>
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
      <v-btn v-if="!clicked"
             @click="clicked=!clicked"
             depressed
             class="follow-button">
        follow
      </v-btn>
      <v-btn v-else
             @click="clicked=!clicked"
             depressed
             class="unfollow-button">
        unfollow
      </v-btn>
    </div>
  </div>
</template>

<script>
export default {
  name: "SuggestedProfile",
  data: function () {
    return {
      clicked: false
    }
  }
}
</script>

<style scoped>

  .suggested-profile {
    display: -webkit-inline-flex;
    flex-direction: row;
    height: auto;
    background-color: transparent;

  }

  .suggestion-head {
    width: 60px;
    height: 60px;
    margin: 0;
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

  .suggestion-head:hover {
    cursor: pointer;
    transition: .3s;
    filter: brightness(0.7);
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

</style>