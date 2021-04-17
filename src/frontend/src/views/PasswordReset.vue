<template>
    <div id="reset-main">
        <v-text-field
        v-model="oldPassword"
        label="Old password"
        hint="At least 8 characters"
        :rules="[rules.required, rules.min]"
        :append-icon="showPassword1 ? 'fa-eye' : 'fa-eye-slash'"
        :type="showPassword1 ? 'text' : 'password'"
        @click:append="showPassword1 = !showPassword1"
        required/>

        <v-text-field
        v-model="newPassword1"
        label="New password"
        hint="At least 8 characters"
        :rules="[rules.required, rules.min, rules.passMatch]"
        :append-icon="showPassword2 ? 'fa-eye' : 'fa-eye-slash'"
        :type="showPassword2 ? 'text' : 'password'"
        @click:append="showPassword2 = !showPassword2"
        required/>

        <v-text-field
        v-model="newPassword2"
        label="Confirm new password"
        hint="At least 8 characters"
        :rules="[rules.required, rules.min, rules.passMatch]"
        :append-icon="showPassword3 ? 'fa-eye' : 'fa-eye-slash'"
        :type="showPassword3 ? 'text' : 'password'"
        @click:append="showPassword3 = !showPassword3"
        required/>

        <v-btn @click="changePassword">Change password</v-btn>
    </div>
</template>

<script>
export default {
    data() {
        return {
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
        changePassword: function() {
            let changeRequest = {
                oldPassword: this.oldPassword,
                newPassword: this.newPassword1,
            };
            this.axios.post("http://localhost:8081/email/change", changeRequest)
                .then(r => {
                    console.log(r);
                })
                .catch(r => {
                    console.log(r);
                });
        }
    },
    computed: {
        passMatch: function() {
            return this.newPassword1 == this.newPassword2 || "Passwords must match"
        },
    }
}
</script>

<style scoped>
    #reset-main {
        display: flex;
        flex-direction: column;
    }
</style>