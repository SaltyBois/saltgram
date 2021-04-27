<template>
    <div id="user-main">
        <div id="user-header">
            <div id="user-icon-logout">
                <i id="user-icon" class="fa fa-user"></i>
                <div id="logout-settings">
                    <v-btn @click="logout" class="accent">Logout</v-btn>
                    <v-btn @click="showSettings=!showSettings">
                        <i class="fa fa-cog" aria-hidden="true"></i>
                    </v-btn>
                </div>
                <div v-if="showSettings" id="settings-container">
                    <v-form id="settings-password" v-model="isFormValid">
                        <h2>Change password:</h2>
                        <b id="err">{{err}}</b>
                        <v-text-field
                        v-model="oldPassword"
                        label="Old password"
                        :rules="[rules.required, rules.min, different]"
                        :append-icon="showPassword2 ? 'fa-eye' : 'fa-eye-slash'"
                        :type="showPassword2 ? 'text' : 'password'"
                        @click:append="showPassword2 = !showPassword2"
                        required/>
                        <v-text-field
                        v-model="newPassword1"
                        label="New password"
                        hint="Min 8 characters, upper/lowercase, number and symbol"
                        :rules="[rules.required, rules.min, rules.passMatch, passMatch, passStr]"
                        :append-icon="showPassword2 ? 'fa-eye' : 'fa-eye-slash'"
                        :type="showPassword2 ? 'text' : 'password'"
                        @click:append="showPassword2 = !showPassword2"
                        required/>

                        <v-text-field
                        v-model="newPassword2"
                        label="Confirm new password"
                        hint="Min 8 characters, upper/lowercase, number and symbol"
                        :rules="[rules.required, rules.min, rules.passMatch, passMatch, passStr]"
                        :append-icon="showPassword3 ? 'fa-eye' : 'fa-eye-slash'"
                        :type="showPassword3 ? 'text' : 'password'"
                        @click:append="showPassword3 = !showPassword3"
                        required/>
                        <b id="pass-str"><div>Password strength: </div><div>{{passScoreText}}</div></b>
                        <password-meter :password="newPassword1" @score="onScore"/>
                        <v-btn :disabled="!isFormValid" @click="changePassword">Change password</v-btn>
                    </v-form>
                </div>
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
import passwordMeter from 'vue-simple-password-meter';
export default {
    components: {
        passwordMeter,
    },
    data: function() {
        return {
            oldPassword: '',
            newPassword1: '',
            newPassword2: '',
            passScore: 0,
            passScoreText: '',
            isFormValid: false,
            showSettings: false,
            user: {},
            rules: {
                required: v => !!v || "Required",
                min: v => v.length >= 8 || "Min 8 characters",
            },
        }
    },

    methods: {
        changePassword: function() {
            this.err = "";

            this.refreshToken()
                .then((r) => {
                    this.$store.state.jws = r.data;
                    console.log("Stored new jws")
                    let changeRequest = {
                        oldPassword: this.oldPassword,
                        newPassword: this.newPassword1,
                    }
                    console.log("Sending new jws")
                    this.axios.post("users/changepass", changeRequest, {headers: {"Authorization": "Bearer " + this.$store.state.jws}})
                        .then(r => {
                            console.log(r);
                            this.showSettings = false;
                        })
                        .catch(r => {
                            console.log(r);
                            this.err = "Invalid old password!";
                        });
                })
                .catch(r => {
                    console.log(r);
                    this.$router.push("/")
                });
        },

        onScore: function({score, strength}) {
            // console.log("Password score: " + strength);
            this.passScore = score;
            this.passScoreText = strength;
        },

        logout: function() {
            this.$store.state.jws = "";
            this.$router.go();
        },

        refreshToken: async function() {
            let jws = this.$store.state.jws
            if (!jws) {
                this.$router.push("/")
            }

            return this.axios.get("auth/refresh", {headers: {"Authorization": "Bearer " + jws}})
        },

        getUserInfo: function() {
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
                    this.this.refreshToken()
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
        this.getUserInfo();
    },
    computed: {
        different: function() {
            return (this.newPassword1 != this.oldPassword && this.newPassword2 != this.oldPassword) || "Cannot use old password"
        },

        passStr: function() {
            console.log("pass score: " + this.passScore);
            return this.passScore > 3 || "Use a stronger password!";
        },

        passMatch: function() {
            return this.newPassword1 == this.newPassword2 || "Passwords must match"
        },
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

    #logout-settings {
        display: flex;
        flex-direction: row;
        justify-content: center;
    }

    #settings-password {
        display: flex;
        flex-direction: column;
        justify-content: center;
        min-width: 25rem;
    }

    #err {
        color: #f00;
    }

    #pass-str {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
    }

    #user-stories {
        text-align: center;
    }

    #user-media {
        text-align: center;
    }
</style>