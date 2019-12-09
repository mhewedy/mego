<template>

  <div id="result-panel">

    <div v-if="loadingResult" class="p-grid">
      <div class="p-col-5"></div>
      <div class="p-col-2" style="text-align: center;">
        <ProgressSpinner mode="indeterminate"/>
      </div>
      <div class="p-col-5"></div>
    </div>

  </div>

</template>

<script>
    import EventService from '../services/events'

    export default {
        name: "Result",
        props: {
            searchInput: null
        },
        data() {
            return {
                loadingResult: false
            }
        },
        watch: {
            searchInput: function (newSearchInput) {
                const that = this;
                this.loadingResult = true;

                EventService.search(newSearchInput, function (data) {
                    console.log(data);
                    that.loadingResult = false;
                }, function (err) {
                    console.log('error:', err);
                    that.loadingResult = false;
                });
            }
        }
    }
</script>

<style scoped>

</style>
