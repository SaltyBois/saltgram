<template>
  <div id="home-main">
    <div id="main-container">
      <div>
        <img id="logo" src="https://image.flaticon.com/icons/png/512/114/114928.png"/>
      </div>
      <div id="login-register">
          <sg-login v-if="login" />
          <sg-register v-else />
          <div id="switch-login">
            <p id="switch-up">
                <span v-if="login">Don't have an account?</span>
                <span v-else>Have an account?</span>
                <a href="#" @click="changeView($event)">
                  <span v-if="login">Sign up</span>
                  <span v-else>Log in</span>
                </a>
            </p>
          </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      login: false,
    }
  },
  methods: {
    changeView: function(e) {
      e.preventDefault();
      this.login = !this.login;
    },
  },
}
</script>

<style scoped>
  #home-main {
    display: grid;
    place-items: center;
    height: 100vh;
    background: #fafafa;
  }

  #main-container {
    display: flex;
    flex-shrink: 0;
    flex-grow: 0;
    flex-direction: row;
    justify-content: center;
    align-content: center;
  }

  #logo {
    display: grid;
    place-items: center;
    flex-basis: 30%;
    width: 30rem;
    height: 30rem;
  }

  #login-register {
    display: flex;
    flex-basis: 70%;
    flex-direction: column;
  }

  #switch-up {
    display: flex;
    flex-direction: row;
    justify-content: space-evenly;
    margin-top: 1rem;
    padding: 1rem 0rem;
    border: 1px solid #eee;
    background: #fff;
  }

  #switch-up a {
    text-decoration: none;
    font-weight: 500;
  }
</style>