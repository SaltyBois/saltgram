<template>
    <div id="register-container">
        <div v-if="emailSent" id="register-and-logo">
            <p id="email-icon"><i class="fa fa-envelope-o"></i></p>
            <p>Your activation email has been sent!</p>
        </div>
        <div v-else>
            <div v-if="processing" id="register-and-logo">
                <v-progress-circular
                indeterminate
                color="primary"/>
            </div>
            <div v-else id="register-and-logo">
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
                    hint="At least 8 characters"
                    :rules="[rules.required, rules.min, rules.passMatch]"
                    :append-icon="showPassword1 ? 'fa-eye' : 'fa-eye-slash'"
                    :type="showPassword1 ? 'text' : 'password'"
                    @click:append="showPassword1 = !showPassword1"
                    required/>
                    <v-text-field
                    v-model="password2"
                    label="Confirm password"
                    hint="At least 8 characters"
                    :rules="[rules.required, rules.min, passMatch]"
                    :append-icon="showPassword2 ? 'fa-eye' : 'fa-eye-slash'"
                    :type="showPassword2 ? 'text' : 'password'"
                    @click:append="showPassword2 = !showPassword2"
                    required/>
                    <vue-recaptcha
                    ref="recaptcha"
                    @verify="onCaptchaVerified"
                    @expired="onCaptchaExpired"
                    size="invisible"
                    :sitekey="sitekey">
                    </vue-recaptcha>
                    <v-btn :disabled="!isFormValid" class="accent" @click="registerUser">Sign up</v-btn>
                </v-form>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: "Register",
    data: function() {
        return {
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
            }
        }
    },

    methods: {
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
            this.axios.post("http://localhost:8081/users", user)
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

    #register-and-logo {
        border: 1px solid #eee;
        padding: 1rem 2rem;
        background: #fff;
    }

    #register {
        display: flex;
        flex-direction: column;
        align-content: center;
    }
</style>
