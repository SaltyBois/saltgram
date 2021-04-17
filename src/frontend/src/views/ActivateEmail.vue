<template>
    <div id="activate-main"></div>
</template>

<script>
export default {
    methods: {
        activateEmail: function() {
            let token = this.$route.params.token
            this.axios.get("http://localhost:8081/email/activate/" + token)
                .then(() => {
                    this.$router.push("/");
                })
                .catch(r => {
                    console.log("Not activated");
                    console.log(r);
                });
        },
    },
    mounted() {
        this.activateEmail();
    },
}
</script>