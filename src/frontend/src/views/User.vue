<template>
    <div id="user-main">
        <div id="user-header">
            <div id="user-icon-logout">
                <i id="user-icon" class="fa fa-user"></i>
                <v-btn @click="logout" class="accent">Logout</v-btn>
            </div>
            <div id="user-info">
                <h2 id="username">{{user.username}}</h2>
                <p>{{user.fullName}}</p>
                <p>{{user.email}}</p>
            </div>
        </div>
        <div id="user-stories">
            Stories
        </div>
        <div id="user-media">
            Media
        </div>
    </div>
</template>

<script>
export default {
    data: function() {
        return {
            user: {},
        }
    },

    methods: {
        logout: function() {
            this.$store.state.jws = "";
            this.$router.go();
        },

        sendJWS: function() {
            let jws = this.$store.state.jws
            if (!jws) {
                this.$router.push("/");
            }

            this.axios.get("users", {headers:{"Authorization": "Bearer " + jws}})
                .then(r => {
                    console.log(r);
                    this.user = r.data;
                })
                .catch(r => {
                    console.log(r);
                    // NOTE(Jovan): Try to refresh
                    // TODO(Jovan): Maybe send existing jwt, just change exp date
                    this.axios.get("http://localhost:8081/auth/refresh", {headers: {"Authorization": "Bearer " + jws}})
                        .then(r => {
                            console.log(r);
                            this.$store.state.jws = r.data;
                            this.$router.go()
                        })
                        .catch(r => {
                            console.log(r);
                            this.$router.push("/");
                        });
                });
        },
    },
    mounted() {
        this.sendJWS();
    },
}
</script>

<style scoped>
    #user-main {
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        align-content: center;
        /* text-align: center; */
        background: #fafafa;
        height: 100vh;
    }

    #user-header {
        display: flex;
        flex-direction: row;
        justify-content: center;
    }

    #user-icon-logout {
        display: flex;
        flex-direction: column;
        justify-content: center;
    }

    #user-icon {
        text-align: center;
        font-size: 8rem;
        margin: 1rem 2rem;
    }

    #username {
        font-weight: 400;
        font-size: 2rem;
        font-family: sans-serif;
    }

    #user-info {
        display: flex;
        flex-direction: column;
        justify-content: flex-start;
        margin: 1rem 2rem;
        /* text-align:left; */
        /* padding: 1rem 2rem;
        background: #fff;
        border: 1px solid #eee; */
    }

    #user-stories {
        text-align: center;
    }

    #user-media {
        text-align: center;
    }
</style>