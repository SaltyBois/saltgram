<template>
    <div>
        <main-navigation>
        </main-navigation>
        <div id="signin-main">
        <div id="signin-container">
            <v-form id="signin-form" v-model="isFormValid">
                <p>Enter your SaltAgent credentials</p>
                <v-text-field
                dark
                v-model="username"
                :rules="[rules.required]"
                label="Username"/>
                <v-text-field
                dark
                v-model="password"
                label="Password"
                :type="showPassword ? 'text' : 'password'"
                :append-icon="showPassword ? 'fa-eye' : 'fa-eye-slash'"
                @click:append="showPassword = !showPassword"
                :rules="[rules.required, rules.password]"
                />
                <v-spacer></v-spacer>
                <v-btn plain color="accent" :disabled="!isFormValid">Sign in</v-btn>
            </v-form>
            <div class="signup">
                <p>Don't have an account? <router-link to="/signup">Sign up!</router-link></p>
            </div>
        </div>
    </div>
    </div>
</template>

<script>
export default {
    data: () => ({
        isFormValid: false,
        username: '',
        password: '',
        rules: {
            required: value => !!value || "Required",
            min: value => value.length >= 8 || "Min 8 characters",
        },
        showPassword: false,
    })
}
</script>

<style scoped>
    #signin-main {
        display: grid;
        place-items: center;
        height: 100vh;
        background: var(--v-primary-base);
    }

    #signin-container {
        display: flex;
        flex-direction: column;
    }

    #signin-form {
        text-align: center;
        display: flex;
        flex-direction: column;
        margin-top: 20px;
        margin-bottom: 20px;
        padding: 10px;
        border: solid 1px var(--v-secondary-lighten1);
        background: #0c162dee;
        /* box-shadow: 0 8px 32px 0 rgba( 31, 38, 135, 0.37 ); */
        backdrop-filter: blur( 8.0px );
        -webkit-backdrop-filter: blur( 8.0px );
        border-radius: 10px;
        color: var(--v-secondary-lighten3);
        height: 30vh;
    }

    #signin-form p {
        font-family: monospace;
    }

    .signup {
        display: grid;
        place-items: center;
        color: #b6b9c0;
        border: solid 1px var(--v-secondary-lighten1);
        background: #0c162dee;
        /* box-shadow: 0 8px 32px 0 rgba( 31, 38, 135, 0.37 ); */
        backdrop-filter: blur( 8.0px );
        -webkit-backdrop-filter: blur( 8.0px );border-radius: 10px;
    }

    .signup p {
        padding: 10px;
        margin: 0;
        font-family: monospace;
    }

    .signup a {
        color: #fff;
        font-family: monospace;
        font-weight: 500;
    }
</style>