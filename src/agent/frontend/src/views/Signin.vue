<template>
    <div id="signin-main">
        <div id="signin-container">
            <h2>SaltAgent signin</h2>
            <p>Enter your SaltAgent credentials</p>
            <v-form id="signin-form">
                <v-text-field
                v-model="username"
                label="Username"/>
                <v-text-field
                v-model="password"
                label="Password"
                :type="showPassword ? 'text' : 'password'"
                />
            </v-form>
        </div>
    </div>
</template>

<script>
export default {
    data: () => ({
        
    })
}
</script>