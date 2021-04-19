<template>
  <div id="home-main">
    <div id="main-container">
      <div>
        <v-img id="logo" src="https://image.flaticon.com/icons/png/512/114/114928.png" contain/>
      </div>
      <sg-login />
    </div>
  </div>
</template>

<style scoped>
  #home-main {
    display: grid;
    place-items: center;
    height: 100vh;
    background: #fafafa;
  }

  #main-container {
    display: flex;
    flex-direction: row;
  }

  #logo {
    height: 40%;
  }
</style>