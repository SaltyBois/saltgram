<template>
  <!--          TODO: JOVAN SETTINGS-->
  <div id="settings-container">
    <v-form id="settings-password" v-model="isFormValid">
      <h2>Change password:</h2>
      <b id="err">{{err}}</b>
      <v-text-field
          v-model="oldPassword"
          label="Old password"
          :rules="[rules.required, rules.min, different]"
          :append-icon="showPassword2 ? 'fa-eye' : 'fa-eye-slash'"
          :type="showPassword2 ? 'text' : 'password'"
          @click:append="showPassword1 = !showPassword1"
          required/>
      <v-text-field
          v-model="newPassword1"
          label="New password"
          hint="Min 8 characters, upper/lowercase, number and symbol"
          :rules="[rules.required, rules.min, rules.passMatch, passMatch, passStr]"
          :append-icon="showPassword2 ? 'fa-eye' : 'fa-eye-slash'"
          :type="showPassword2 ? 'text' : 'password'"
          @click:append="showPassword2 = !showPassword2"
          required/>

      <v-text-field
          v-model="newPassword2"
          label="Confirm new password"
          hint="Min 8 characters, upper/lowercase, number and symbol"
          :rules="[rules.required, rules.min, rules.passMatch, passMatch, passStr]"
          :append-icon="showPassword3 ? 'fa-eye' : 'fa-eye-slash'"
          :type="showPassword3 ? 'text' : 'password'"
          @click:append="showPassword3 = !showPassword3"
          required/>
      <b id="pass-str">Password strength: {{passScoreText}}</b>
      <password-meter :password="newPassword1" @score="onScore"/>
      <v-btn :disabled="!isFormValid" class="mt-5 primary" @click="changePassword">Change password</v-btn>
    </v-form>
  </div>
  <!--        TODO: END OF JOVAN SETTINGS-->
</template>

<script>
import passwordMeter from 'vue-simple-password-meter';

export default {
  name: "ChangePassword",
  components: {passwordMeter},
  data: function () {
    return {
      err: "",
      showPassword1: false,
      showPassword2: false,
      showPassword3: false,
      oldPassword: '',
      newPassword1: '',
      newPassword2: '',
      passScore: 0,
      passScoreText: '',
      isFormValid: false,
      showSettings: false,
      descriptionTextBoxReadOnly: true,
      radioButton: 'posts',
      rules: {
        required: v => !!v || "Required",
        min: v => v.length >= 8 || "Min 8 characters",
      },
    }
  },
  methods: {
    changePassword: function() {
      this.err = "";
      this.refreshToken(this.getAHeader())
          .then((r) => {
            this.$store.state.jws = r.data;
            let changeRequest = {
              oldPassword: this.oldPassword,
              newPassword: this.newPassword1,
            }
            this.axios.post("users/changepass", changeRequest, {headers: {"Authorization": "Bearer " + this.$store.state.jws}})
                .then(r => {
                  console.log(r);
                  this.showSettings = false;
                })
                .catch(r => {
                  console.log(r);
                  this.err = "Invalid old password!";
                });
          })
          .catch(r => {
            console.log(r);
            this.$router.push("/")
          });
    },
    onScore: function({score, strength}) {
      // console.log("Password score: " + strength);
      this.passScore = score;
      this.passScoreText = strength;
    },
  },
  computed: {
    different: function() {
      return (this.newPassword1 !== this.oldPassword && this.newPassword2 !== this.oldPassword) || "Cannot use old password"
    },

    passStr: function() {
      console.log("pass score: " + this.passScore);
      return this.passScore > 3 || "Use a stronger password!";
    },

    passMatch: function() {
      return this.newPassword1 === this.newPassword2 || "Passwords must match"
    },
  },
}
</script>

<style scoped>

#settings-container {
  padding: 20px;
  text-align: center;
}

#settings-password {
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-width: 25rem;
}

#err {
  color: #f00;
}

#pass-str {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
}

</style>