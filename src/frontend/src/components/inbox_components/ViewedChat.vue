<template>
  <div class="chat-div">
    <div style="padding: 3%;">
      <v-img  class="chat-head"
              src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRYYR7aw1k6gkP_UGKGddjPee4ENRCEY8AO3A&usqp=CAU"
              @click="$router.push('/user')"
              alt="Profile picture"/>
    </div>
    <div style="padding: 3%;overflow: hidden">
      <h4>Username1</h4>
      <p>Message content &#183; 10m ago</p>
    </div>
  </div>
</template>

<script>
export default {
name: "ViewedChat"
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

.chat-head:hover {
  cursor: pointer;
  transition: .3s;
  filter: brightness(0.7);
}

</style>