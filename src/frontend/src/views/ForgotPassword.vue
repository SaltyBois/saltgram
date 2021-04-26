<template>
    <div id="forgot-main">
        <v-text-field 
        v-model="email"
        label="Email"
        :rules="[rules.required, rules.email]"
        required/>
        <v-btn @click="requestReset">Request reset</v-btn>
    </div>
</template>

<script>
export default {
    data() {
        return {
            email: "",
            rules: {
                required: v => !!v || "Required",
                email: v => !v || /^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(v) || 'E-mail must be valid',
            },
        }
    },
    methods: {
        requestReset: function() {
            this.axios.post("/email/forgot", this.email, {headers: {withCredentials: true}})
                .then(r => {
                    console.log(r);
                })
                .catch(r => {
                    console.log(r);
                });
        },
    },
}
</script>