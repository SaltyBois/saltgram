<template>
    <div id="user-main">
        {{user}}
    </div>
</template>

<script>
export default {
    data: function() {
        return {
            user: {}
        }
    },

    methods: {
        sendJWS: function() {
            let jws = localStorage.getItem("jws");
            if (!jws) {
                this.$router.push("/");
            }

            this.axios.get("/users", {headers:{"Authorization": "Bearer " + jws}})
                .then(r => {
                    console.log(r);
                    this.user = r.data;
                })
                .catch(r => {
                    console.log(r);
                    this.$router.push("/");
                });
        },
    },
    mounted() {
        this.sendJWS();
    },
}
</script>