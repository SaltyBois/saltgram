<template>
  <div id="home-main">
    <div id="main-container">
      <div>
        <v-img id="logo" src="https://image.flaticon.com/icons/png/512/114/114928.png" contain/>
      </div>
      <div id="login-register">
        <sg-login v-if="login" @change="switchView" />
        <sg-register v-else  @change="switchView" />
        <div id="switch-login">
          <p id="switch-up">
              <span v-if="login">Don't have an account?</span>
              <span v-else>Have an account?</span>
              <a href="#" @click="login = !login">
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
    changeView: function(value) {
      console.log(value);
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
    flex-direction: row;
  }

  #logo {
    height: 40%;
  }

  #login-register {
    display: flex;
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
    font-weight: 500;
  }
</style>