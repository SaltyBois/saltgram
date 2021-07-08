<template>
    <v-layout row
        style="justify-content: center; margin-top: 80px">
        <v-layout column
                  align-center>
                  <h2>Profile not found.</h2>
                  <br>
                  <v-btn class="btn-style"
                    @click="goBack()">
                      Go to main page
                  </v-btn>
        </v-layout>
    </v-layout>
</template>

<script>
export default {
    methods: {
        goBack: function() {
            this.$router.push('/main');
        }
    }
}
</script>
<style scoped>
.btn-style{
  background-color: transparent;
  color: #016ddb;
  border-color: #016ddb;
  border-style: solid;
  border-width: 1px;
  text-align: -webkit-center;
}

</style>