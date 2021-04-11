<template>
    <div id="register-main">
        <div id="registration-container">
            <h1 id="title">Saltgram</h1>
            <div id="registration">
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
                :rules="[rules.required, rules.min]"
                :append-icon="showPassword1 ? 'fa-eye' : 'fa-eye-slash'"
                :type="showPassword1 ? 'text' : 'password'"
                @click:append="showPassword1 = !showPassword1"
                required/>
                <v-text-field
                v-model="password2"
                label="Confirm password"
                hint="At least 8 characters"
                :rules="[rules.required, rules.min, rules.passMatch]"
                :append-icon="showPassword2 ? 'fa-eye' : 'fa-eye-slash'"
                :type="showPassword2 ? 'text' : 'password'"
                @click:append="showPassword2 = !showPassword2"
                required/>
                <v-btn class="accent" @click="registerUser">Sign up</v-btn>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: "Register",
    data: function() {
        return {
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
                passMatch: () => this.password1 == this.password2 || "Passwords must match",
            }
        }
    },

    methods: {
        registerUser: function() {
            // TODO(Jovan): Validate
            let user = {
                username: this.username,
                fullName: this.fullName,
                email: this.email,
                password: this.password1,
            }
            this.axios.post("/users", user)
                .then(response => {
                    console.log(response);
                })
                .catch(response => {
                    console.log(response);
                })
        }
    },
}
</script>

<style scoped>
    #title {
        font-family: "Lucida Handwriting", cursive;
    }

    #register-main {
        display: grid;
        place-items: center;
        height: 100vh;
    }

    #registration-container {
        display: flex;
        flex-direction: column;
        text-align: center;
    }

    #registration {
        display: flex;
        flex-direction: column;
        justify-content: center;
    }
</style>