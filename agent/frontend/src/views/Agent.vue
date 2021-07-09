<template>
    <div>
        <main-navigation></main-navigation>
        <div id="agent-main">
            <div id="agent-container">
                <v-tabs
                v-model="tab"
                background-color="#0c162dee"
                color="accent"
                flat>
                    <v-tab>Campaigns</v-tab>
                    <v-tab>Influencers</v-tab>
                    <v-tab>Products</v-tab>
                </v-tabs>
                <v-tabs-items v-model="tab" >
                    <v-tab-item>
                        <div id="campaign-container">
                            <div v-for="campaign in campaigns" :key="campaign.id" class="campaign">
                                <b>{{campaign.productName}}</b>
                                <b>{{campaign.type}}</b>
                            </div>
                        </div>
                    </v-tab-item>
                    <v-tab-item>
                        Jos nesto
                    </v-tab-item>
                    <v-tab-item>
                        Opet nesto
                    </v-tab-item>
                </v-tabs-items>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    data: () => ({
        tab: null,
    }),
}
</script>

<style scoped>
    #agent-main {
        display: grid;
        place-items: center;
        min-height: 100vh;
        background: var(--v-primary-base);
        margin-top: 8vh;
    }

    .v-tabs {
        border: solid 1px var(--v-secondary-lighten1);
        margin-bottom: 20px;
        backdrop-filter: blur( 8.0px );
        -webkit-backdrop-filter: blur( 8.0px );border-radius: 5px;
    }

    .v-tab:not(.v-tab--active) {
        color: var(--v-secondary-lighten4) !important;
    }

    .v-tabs-items {
        border: solid 1px var(--v-secondary-lighten1);
        backdrop-filter: blur( 8.0px );
        background: #0c162dee !important;
        -webkit-backdrop-filter: blur( 8.0px );border-radius: 5px;
        color: var(--v-secondary-lighten4);
    }


</style>