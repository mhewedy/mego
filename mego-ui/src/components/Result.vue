<template>

  <div id="result-panel">

  </div>

</template>

<script>
    export default {
        name: "Result",
        props: {
            searchInput: null
        },
        mounted() {
            console.log("searchInput", this.searchInput)
        }
    }
</script>

<style scoped>

</style>
