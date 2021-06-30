<template>
    <div id="main">
        <div id="container">
            <v-text-field @keyup="search" id="search-input" placeholder="Enter location" v-model="query" />
            <div v-if="results.length > 0" id="results">
                <div
                v-for="r in results"
                :key="r.place_id"
                @click="select(r)"
                class="result">{{r.display_name}}</div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    data: () => ({
        timeout: null,
        results: [],
        selected: null,
        query: '',
    }),

    methods: {
        select: function(r) {
            this.selected = r;
            this.results = [];
            this.$emit('selected', r);
        },
        search: function() {
            clearTimeout(this.timeout);
            this.results = [];
            let self = this;
            this.timeout = setTimeout(function() {
                self.axios.get('https://eu1.locationiq.com/v1/search.php?key=8d176a2a842649&format=json&addressdetails=1&limit=5&q=' + self.query,
                {withCredentials: false})
                    .then(r => {self.results = r.data; console.log(r.data)})
                    .catch(r => console.log(r));
            }, 1000);
        },
    },
}
</script>

<style scoped>
    #container {
        position: relative;
        display: inline-block;
    }

    #results {
        position: absolute;
        top: 100%;
        left: 0;
        right: 0;
        z-index: 99;
        border: solid 1px #eee;
    }

    #container input {
        margin-bottom: 0;
    }

    .result {
        cursor: pointer;
        border: solid 1px #eee;
        background: #fff;
    }

    .result:hover {
        background: #eee;
    }
</style>