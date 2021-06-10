<template>
    <div id="register-container">
        <div v-if="emailSent" class="register-and-logo">
            <p id="email-icon"><i class="fa fa-envelope-o"></i></p>
            <p>Your activation email has been sent!</p>
        </div>
        <div v-else>
            <div v-if="processing" class="register-and-logo">
                <v-progress-circular
                indeterminate
                color="primary"/>
            </div>
            <div v-else class="register-and-logo">
                <h1 id="title">Saltgram</h1>
                <v-form id="register" v-model="isFormValid">
                    <v-text-field 
                    v-model="email"
                    label="Email"
                    :rules="[rules.required, rules.email]"
                    required/>
                    <v-text-field 
                    v-model="fullName"
                    label="Full name"
                    :rules="[rules.required]"
                    required/>
                    <v-text-field
                    v-model="username"
                    label="Username"
                    :rules="[rules.required]"
                    required/>
                    <v-text-field
                    v-model="password1"
                    label="Password"
                    hint="Min 8 characters, upper/lowercase, number and symbol"
                    :rules="[rules.required, rules.min, passMatch, passStr]"
                    :append-icon="showPassword1 ? 'fa-eye' : 'fa-eye-slash'"
                    :type="showPassword1 ? 'text' : 'password'"
                    @click:append="showPassword1 = !showPassword1"
                    required/>
                    <v-text-field
                    v-model="password2"
                    label="Confirm password"
                    hint="Min 8 characters, upper/lowercase, number and symbol"
                    :rules="[rules.required, rules.min, passMatch, passStr]"
                    :append-icon="showPassword2 ? 'fa-eye' : 'fa-eye-slash'"
                    :type="showPassword2 ? 'text' : 'password'"
                    @click:append="showPassword2 = !showPassword2"
                    required/>
                    <b id="pass-str"><div>Password strength: </div><div>{{passScoreText}}</div></b>
                    <password-meter :password="password1" @score="onScore"/>
                    <vue-recaptcha
                    ref="recaptcha"
                    @verify="onCaptchaVerified"
                    @expired="onCaptchaExpired"
                    size="invisible"
                    :sitekey="sitekey">
                    </vue-recaptcha>
                    <v-text-field
                    label="Add Description"
                    hint="start with + please"
                    v-model="description"/>
                    <v-text-field
                    label="Phone Number"
                    hint="start with + please"/>
                    <v-select
                    :items="genderRules"
                    label="Gender"/>
                    <v-date-picker />
                    <v-text-field
                    label="Web Site"
                    hint="You can enter it without 'www.'"/>
                    <v-btn class="mt-3 mb-3" v-bind:class="!privateProfile ? 'primary' : ''" @click="privateProfile = false"><i class="fa fa-unlock mr-1"/>Public profile</v-btn>
                    <v-btn class="mb-10" v-bind:class="privateProfile ? 'primary' : ''" @click="privateProfile = true"><i class="fa fa-lock mr-1"/>Private profile</v-btn>
                    <v-btn :disabled="!isFormValid" class="accent" @click="registerUser">Sign up</v-btn>
                </v-form>
            </div>
        </div>
    </div>
</template>

<script>
import passwordMeter from 'vue-simple-password-meter';
export default {
    name: "Register",
    components: {
        passwordMeter
    },
    data: function() {
        return {
            passScore: 0,
            passScoreText: '',
            emailSent: false,
            processing: false,
            isFormValid: false,
            reCaptchaStatus: "",
            fullName: "",
            username: "",
            email: "",
            password1: "",
            password2: "",
            showPassword1: false,
            showPassword2: false,
            rules: {
                required: v => !!v || "Required",
                min: v => v.length >= 8 || "Min 8 characters",
                email: v => !v || /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(v) || 'E-mail must be valid',
            },
            genderRules: [ 'Male', 'Female' ],
            privateProfile: false,
            description: ''
        }
    },

    methods: {
        onScore: function({score, strength}) {
            // console.log("Password score: " + strength);
            this.passScore = score;
            this.passScoreText = strength;
        },

        registerUser: function() {
            this.$refs.recaptcha.execute();
        },

        onCaptchaVerified: function(token) {
            // TODO(Jovan): Validate
            this.processing = true;
            this.emailSent = false;
            this.reCaptchaStatus = "submitting";
            let user = {
                username: this.username,
                fullName: this.fullName,
                email: this.email,
                password: this.password1,
                reCaptcha: {
                    token: token,
                    action: "register",
                }
            }
            this.axios.post("users/register", user)
                .then(response => {
                    console.log("Registered");
                    console.log(response);
                    this.emailSent = true;
                })
                .catch(response => {
                    console.log("Not registered");
                    console.log(response);
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

        passStr: function() {
            console.log("pass score: " + this.passScore);
            return this.passScore > 3 || "Use a stronger password!";
        },

        sitekey: function() {
            return process.env.VUE_APP_RECAPTCHA_SITE_KEY;
        },

        passMatch: function() {
            return this.password1 == this.password2 || "Passwords must match"
        },
    },
}
</script>

<style scoped>
    #email-icon {
        font-size: 3rem;
    }
    #title {
        font-size: 2.5rem;
        font-family: "Lucida Handwriting", cursive;
    }

    #register-container {
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        text-align: center;
    }

    .register-and-logo {
        min-width: 25rem;
        border: 1px solid #eee;
        padding: 1rem 2rem;
        background: #fff;
    }

    #register {
        display: flex;
        flex-direction: column;
        align-content: center;
    }
    
    #pass-str {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
    }

</style>
