<template>
    <div>
        <main-navigation></main-navigation>
        <div id="signup-main">
            <div id="signup-container">
                <v-form id="signup-form">
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
                    :rules="[rules.required, rules.min, passMatch, passStr]"
                    :append-icon="showPassword1 ? 'fa-eye' : 'fa-eye-slash'"
                    :type="showPassword1 ? 'text' : 'password'"
                    @click:append="showPassword1 = !showPassword1"
                    required/>
                    <v-text-field
                    dark
                    v-model="password2"
                    label="Confirm password"
                    hint="Min 8 characters, upper/lowercase, number and symbol"
                    :rules="[rules.required, rules.min, passMatch, passStr]"
                    :append-icon="showPassword2 ? 'fa-eye' : 'fa-eye-slash'"
                    :type="showPassword2 ? 'text' : 'password'"
                    @click:append="showPassword2 = !showPassword2"
                    required/>
                    <v-spacer></v-spacer>
                    <v-btn depressed color="accent">Sign up</v-btn>
                </v-form>
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
        username: '',
        password: '',
        rules: {
            required: v => !!v || "Required",
            min: v => v.length >= 8 || "Min 8 characters",
            email: v => !v || /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(v) || 'E-mail must be valid',
        },
        showPassword1: false,
        showPassword2: false,
    })
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

    #signup-form {
        text-align: center;
        display: flex;
        flex-direction: column;
        margin-top: 20px;
        margin-bottom: 20px;
        padding: 10px;
        border: solid 1px #eee;
        background: #0c162dee;
        border: solid 1px var(--v-secondary-lighten1);
        /* box-shadow: 0 8px 32px 0 rgba( 31, 38, 135, 0.37 ); */
        backdrop-filter: blur( 8.0px );
        -webkit-backdrop-filter: blur( 8.0px );border-radius: 10px;
        color: var(--v-secondary-lighten3);
        height: 42vh;
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