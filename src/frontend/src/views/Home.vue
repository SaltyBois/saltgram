<template>
  <div id="home-main">
    <div id="main-container">
      <div>
        <v-img id="logo" src="https://image.flaticon.com/icons/png/512/114/114928.png" contain/>
      </div>
      <div id="login-container">
        <h1 id="home-title">Saltgram</h1>
        <v-form id="login" v-model="isFormValid">
          <v-text-field 
          v-model="user.username"
          label="Username"
          required/>
          <v-text-field 
          v-model="user.password"
          label="Password"
          hint="At least 8 characters"
          :rules="[rules.required, rules.min]"
          :append-icon="showPassword ? 'fa-eye' : 'fa-eye-slash'"
          :type="showPassword ? 'text' : 'password'"
          @click:append="showPassword = !showPassword"
          required/>
           <vue-recaptcha
              ref="recaptcha"
              @verify="onCaptchaVerified"
              @expired="onCaptchaExpired"
              size="invisible"
              :sitekey="sitekey">
            </vue-recaptcha>
          <v-btn :disabled="!isFormValid" class="accent" @click="login">Log in</v-btn>
        </v-form>
        <div id="sign-up">
          <p>Don't have an account? <router-link to="/register">Sign up</router-link></p>
          <p><router-link to="/forgotpassword">Forgot password</router-link></p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>

export default {
  name: 'Home',
  data: function() {
    return {
      isFormValid: false,
      captchaResponse: "",
      reCaptchaStatus: "submitting",
      user: {
        username: "",
        password: "",
        reCaptcha: {
          token: "",
          action: "",
        },
      },
      rules: {
        required: value => !!value || "Required",
        min: value => value.length >= 8 || "Min 8 characters",
      },
      showPassword: false,
    }
  },
  methods: {
    login: function() {
      this.$refs.recaptcha.execute();
    },

    onCaptchaVerified: function(token) {
      this.reCaptchaStatus = "submitting";
      this.$refs.recaptcha.reset();
      this.user.reCaptcha.token = token;
      this.user.reCaptcha.action = "login";

      this.axios.post("http://localhost:8081/login", this.user)
        .then(r => {
          console.log(r);
          this.axios.post("http://localhost:8081/auth/jwt", r.data)
            .then(r => {
              localStorage.setItem("jws", r.data);
              this.$router.push("/user");
            })
            .catch(r => {
              console.log(r);
            });
        })
        .catch(r => {
          console.log(r);
        });
    },

    onCaptchaExpired: function() {
      this.$refs.recaptcha.reset();
    },
  },
  computed: {
    sitekey: function() {
      return process.env.VUE_APP_RECAPTCHA_SITE_KEY;
    }
  },
}
</script>

<style scoped>
  #home-title {
    font-family: "Lucida Handwriting", cursive;
  }

  #home-main {
    display: grid;
    place-items: center;
    height: 100vh;
  }

  #main-container {
    display: flex;
    flex-direction: row;
  }

  #login-container {
    display: flex;
    flex-direction: column;
    align-content: center;
    text-align: center;
  }

  #login {
    display: flex;
    flex-direction: column;
    align-content: center;
  }

  #logo {
    height: 40%;
  }
</style>