<template>
  <div class="chat-div">
    <div style="padding: 3%">
      <v-img  class="chat-head"
              src="https://i.pinimg.com/474x/ab/62/39/ab6239024f15022185527618f541f429.jpg"
              @click="$router.push('/user')"
              alt="Profile picture"/>
    </div>
    <div style="padding: 3%;overflow: hidden">
      <b>
        <h3>Username1</h3>
        <p>Message content 	&#183; now</p>
      </b>
    </div>
    <div class="blue-dot">
      <p disabled="true" style="transform: scale(5); color: #016ddb; margin-top: 15px">&#183;</p>
    </div>
  </div>
</template>

<script>
export default {
  name: "UnviewedChat"
}
</script>

<style scoped>

  .chat-div {
    display: flex;
    flex-direction: row;

    height: 80px;
    width: 100%;
    background-color: transparent;
    cursor: pointer;
    transition: 0.3s;
  }

  .chat-div:hover {
    background-color: #949494;
    transition: 0.3s;
  }

  .chat-div:checked {
    background-color: #4b90b5;
    transition: 0.3s;
  }

  .chat-head {
    width: 60px;
    height: 60px;
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

  .blue-dot {
    position: relative;
    top: 0;
    left: 25%;

    width: auto;
    height: auto;

    float: right;
    padding: 3%;

  }

</style>