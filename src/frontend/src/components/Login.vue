<template>
    <div id="login-container">
      <div v-if="processing" id="login-and-logo">
        <v-progress-circular
        indeterminate
        color="primary"/>
      </div>
      <div v-else id="login-and-logo">
        <h1 id="home-title">Saltgram</h1>
        <v-form id="login" v-model="isFormValid">
          <span class="err">{{error}}</span>
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
          <p id="forgot-password"><router-link to="/forgotpassword">Forgot password?</router-link></p>
        </v-form>
      </div>
    </div>
</template>

<script>
export default {
  name: 'Home',
  data: function() {
    return {
      error: "",
      processing: false,
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
      this.error = "";
      this.processing = true;
      this.reCaptchaStatus = "submitting";
      this.$refs.recaptcha.reset();
      this.user.reCaptcha.token = token;
      this.user.reCaptcha.action = "login";

      this.axios.post("auth/login", this.user)
        .then(r => {
          console.log(r);
          this.axios.post("auth/jwt", r.data)
            .then(r => {
              this.$store.state.jws = r.data;
              this.$router.push("/user");
            })
            .catch(r => {
              console.log(r);
            });
        })
        .catch(r => {
          console.log(r);
          this.error = "Invalid username and/or password"
        })
        .finally(() => {
          this.processing = false;
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

  .err {
    border-color: #fff impor !important;
    color: #f00;
  }

  a {
    text-decoration: none;
  }
  #login-container {
      display: flex;
      flex-direction: column;
      justify-content: flex-start;
      text-align: center;
  }

  #login-and-logo {
      border: 1px solid #eee;
      padding: 1rem 2rem;
      background: #fff;
  }

  #login {
      display: flex;
      flex-direction: column;
      align-content: center;
  }

  #forgot-password {
    margin-top: 1rem;
    padding-top: 1rem;
    margin-bottom: 0px;
    border-top: 1px solid #eee;
  }

  #home-title {
    font-size: 2.5rem;
    font-family: "Lucida Handwriting", cursive;
  }

</style>