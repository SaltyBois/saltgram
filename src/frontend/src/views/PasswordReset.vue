<template>
    <div id="reset-main">
        <v-form id="reset-container" v-model="isFormValid">
            <v-text-field
            v-model="oldPassword"
            label="Old password"
            :rules="[rules.required, rules.min, different]"
            :append-icon="showPassword1 ? 'fa-eye' : 'fa-eye-slash'"
            :type="showPassword1 ? 'text' : 'password'"
            @click:append="showPassword1 = !showPassword1"
            required/>

            <v-text-field
            v-model="newPassword1"
            label="New password"
            hint="Min 8 characters, upper/lowercase, number and symbol"
            :rules="[rules.required, rules.min, rules.passMatch, different, passMatch, passStr]"
            :append-icon="showPassword2 ? 'fa-eye' : 'fa-eye-slash'"
            :type="showPassword2 ? 'text' : 'password'"
            @click:append="showPassword2 = !showPassword2"
            required/>

            <v-text-field
            v-model="newPassword2"
            label="Confirm new password"
            hint="Min 8 characters, upper/lowercase, number and symbol"
            :rules="[rules.required, rules.min, rules.passMatch, different, passMatch, passStr]"
            :append-icon="showPassword3 ? 'fa-eye' : 'fa-eye-slash'"
            :type="showPassword3 ? 'text' : 'password'"
            @click:append="showPassword3 = !showPassword3"
            required/>
            <b id="pass-str"><div>Password strength: </div><div>{{passScoreText}}</div></b>
            <password-meter :password="newPassword1" @score="onScore"/>

            <v-btn :disabled="!isFormValid" @click="changePassword">Change password</v-btn>
        </v-form>
    </div>
</template>

<script>
import passwordMeter from 'vue-simple-password-meter';
export default {
    components: {
        passwordMeter
    },
    data() {
        return {
            passScore: 0,
            passScoreText: '',
            isFormValid: false,
            oldPassword: "",
            newPassword1: "",
            newPassword2: "",
            showPassword1: false,
            showPassword2: false,
            showPassword3: false,
            rules: {
                required: v => !!v || "Required",
                min: v => v.length >= 8 || "Min 8 characters",
            },
        }
    },
    methods: {
        onScore: function({score, strength}) {
            // console.log("Password score: " + strength);
            this.passScore = score;
            this.passScoreText = strength;
        },

        changePassword: function() {
            let changeRequest = {
                oldPassword: this.oldPassword,
                newPassword: this.newPassword1,
            };
            this.axios.post("/email/change", changeRequest)
                .then(r => {
                    console.log(r);
                    this.$router.push("/");
                })
                .catch(r => {
                    console.log(r);
                });
        }
    },
    computed: {
        passStr: function() {
            console.log("pass score: " + this.passScore);
            return this.passScore > 3 || "Use a stronger password!";
        },

        passMatch: function() {
            return this.newPassword1 == this.newPassword2 || "Passwords must match"
        },

        different: function() {
            return (this.newPassword1 != this.oldPassword && this.newPassword2 != this.oldPassword) || "Cannot use old password"
        }
    }
}
</script>

<style scoped>
    #reset-main {
        display: grid;
        place-items: center;
        height: 100vh;
        background: #fafafa;
    }

    #reset-container {
        display: flex;
        flex-direction: column;
        justify-content: center;
        min-width: 25rem;
        background: #fff;
        border: 1px solid #eee;
        padding: 1rem 2rem;
    }

    #pass-str {
        display: flex;
        flex-direction: column;
        justify-content: space-between;
    }
</style>