<template>
    <div>
        <main-navigation></main-navigation>
        <div id="signup-main">
            <div id="signup-container">
                <v-stepper id="step-container" v-model="steps" elevation="0">
                    <v-stepper-items>
                        <v-stepper-content step="1">
                            <div id="signup-role">
                                <v-btn plain color="lightsecondary" @click="steps++; agent=false" >Regular user ></v-btn>
                                <v-btn plain color="lightsecondary" @click="steps++; agent=true" >Agent user > </v-btn>
                            </div>
                        </v-stepper-content>
                        <v-stepper-content step="2">
                            <v-form id="signup-form" v-model="formValid">
                                <p>Create an account</p>
                                <v-text-field
                                dark
                                v-model="username"
                                :rules="[rules.required]"
                                label="Username"/>
                                <v-text-field 
                                dark
                                v-model="email"
                                label="Email"
                                :rules="[rules.required, rules.email]"
                                required/>
                                <v-text-field
                                dark
                                v-model="password1"
                                label="Password"
                                hint="Min 8 characters, upper/lowercase, number and symbol"
                                :rules="[rules.required, rules.min, passMatch]"
                                :append-icon="showPassword1 ? 'fa-eye' : 'fa-eye-slash'"
                                :type="showPassword1 ? 'text' : 'password'"
                                @click:append="showPassword1 = !showPassword1"
                                required/>
                                <v-text-field
                                dark
                                v-model="password2"
                                label="Confirm password"
                                hint="Min 8 characters, upper/lowercase, number and symbol"
                                :rules="[rules.required, rules.min, passMatch]"
                                :append-icon="showPassword2 ? 'fa-eye' : 'fa-eye-slash'"
                                :type="showPassword2 ? 'text' : 'password'"
                                @click:append="showPassword2 = !showPassword2"
                                required/>
                                <v-text-field
                                v-if="agent"
                                dark
                                v-model="token"
                                :rules="[rules.required]"
                                label="API Token"/>
                                <v-spacer></v-spacer>
                                <div class="d-flex flex-row">
                                    <v-btn plain color="lightsecondary" @click="steps--">&#60; Back</v-btn>
                                    <v-spacer></v-spacer>
                                    <v-btn :disabled="!formValid" plain color="accent" @click="signup">Sign up</v-btn>
                                </div>
                            </v-form>
                        </v-stepper-content>
                    </v-stepper-items>
                </v-stepper>
                <div class="signin">
                    <p>Have an account? <router-link to="/signin">Sign in!</router-link></p>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    data: () => ({
        formValid: false,
        steps: 1,
        username: '',
        password1: '',
        password2: '',
        email: '',
        agent: false,
        token: '',
        rules: {
            required: v => !!v || "Required",
            min: v => v.length >= 8 || "Min 8 characters",
            email: v => !v || /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(v) || 'E-mail must be valid',
        },
        showPassword1: false,
        showPassword2: false,
    }),

    methods: {
        signup: function() {
            let request = {
                username: this.username,
                password: this.password1,
                email: this.email,
                agent: this.agent,
            };

            this.axios.post('/agent/signup', request)
                .then(() => this.$router.push('/signin'))
                .catch(r => console.log(r));
        },
    },

    computed: {

        passMatch: function() {
            return this.password1 == this.password2 || "Passwords must match"
        },
    },
}
</script>

<style scoped>
    #signup-main {
        display: grid;
        place-items: center;
        height: 100vh;
        background: var(--v-primary-base);
    }

    #signup-container {
        display: flex;
        flex-direction: column;
    }

    #step-container {
        border: solid 1px #eee;
        background: #0c162dee;
        border: solid 1px var(--v-secondary-lighten1);
        backdrop-filter: blur( 8.0px );
        -webkit-backdrop-filter: blur( 8.0px );border-radius: 10px;
        margin-bottom: 20px;
    }

    #signup-role {
        display: flex;
        flex-direction: column;
    }

    #signup-form {
        text-align: center;
        display: flex;
        flex-direction: column;
        margin-top: 20px;
        padding: 10px;
        color: var(--v-secondary-lighten3);
        height: 480px;
        width: 300px;
    }

    #signup-form p {
        font-family: monospace;
    }

    .signin {
        display: grid;
        place-items: center;
        color: #b6b9c0;
        border: solid 1px #eee;
        background: #0c162dee;
        border: solid 1px var(--v-secondary-lighten1);
        background: #0c162dee;
        /* box-shadow: 0 8px 32px 0 rgba( 31, 38, 135, 0.37 ); */
        backdrop-filter: blur( 8.0px );
        -webkit-backdrop-filter: blur( 8.0px );border-radius: 10px;
    }

    .signin p {
        padding: 10px;
        margin: 0;
        font-family: monospace;
    }

    .signin a {
        color: #fff;
        font-family: monospace;
        font-weight: 500;
    }
</style>