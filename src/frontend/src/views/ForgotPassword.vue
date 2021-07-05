<template>
    <div id="forgot-main">
        <v-form v-if="!emailSent" id="forgot-container" v-model="isFormValid">
            <v-text-field 
            v-model="email"
            label="Email"
            :rules="[rules.required, rules.email]"
            required/>
            <v-btn :disabled="!isFormValid" @click="requestReset">Request reset</v-btn>
        </v-form>
        <div v-else id="forgot-container">
            <p id="email-icon"><i class="fa fa-envelope-o"></i></p>
            <p>Reset email sent!</p>
        </div>
    </div>
</template>

<script>
export default {
    data() {
        return {
            isFormValid: false,
            emailSent: false,
            email: "",
            rules: {
                required: v => !!v || "Required",
                email: v => !v || /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(v) || 'E-mail must be valid',
            },
        }
    },
    methods: {
        requestReset: function() {
            this.axios.post("email/forgot", this.email, {headers: {withCredentials: true}})
                .then(r => {
                    console.log(r);
                    this.emailSent = true;
                })
                .catch(r => {
                    console.log(r);
                });
        },
    },
}
</script>

<style scoped>
    #forgot-main {
        display: grid;
        place-items: center;
        height: 100vh;
        background: #fafafa;
    }

    #forgot-container {
        display: flex;
        flex-direction: column;
        justify-content: center;
        background: #fff;
        border: 1px solid #eee;
        padding: 1rem 2rem;
        text-align: center;
    }

    #email-icon {
        font-size: 3rem;
    }
</style>