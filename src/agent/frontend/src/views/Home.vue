<template>
  <div>
    <main-navigation>
      <router-link to="/signin">Sign in</router-link>
    </main-navigation>
    <div id="home-main">
      <div id="home-container">
        <h1>Welcome</h1>
        <p>Join SaltAgent</p>
        <v-btn to="/signup" id="signup" depressed color="accent">Sign up</v-btn>
      </div>
    </div>
  </div>
</template>

<script>

  export default {
    name: 'Home',
  }
</script>


<style scoped>
  #home-main {
    display: grid;
    place-items: center;
    background: var(--v-primary-base);
    height: 100vh;
  }

  #home-container {
    display: flex;
    flex-direction: column;
    text-align: center;
  }

  #home-container h1 {
    font-size: 4rem;
    color: var(--v-secondary-lighten2);
    animation: 1.5s ease-in-out 0s 1 bottomReveal;
  }  

  #home-container p {
    font-family: monospace;
    font-size: 2rem;
    color: var(--v-secondary-lighten2);
    animation: 1.5s ease-in-out 0.5s 1 bottomReveal;
  }

  #signup {
    animation: 1.5s ease-in-out 1s 1 bottomReveal;
    animation-fill-mode: both;
  }

  @keyframes bottomReveal {
    0% {
      transform: translate(0, 100px);
      opacity: 0;
    }

    100% {
      transform: translate(0, 0);
      opacity: 1;
    }
  }
</style>