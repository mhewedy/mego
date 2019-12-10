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
    import MessageService from '../services/messages'

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
        mounted() {
            this.search(this.searchInput)
        },
        watch: {
            searchInput: function (newSearchInput) {
                this.search(newSearchInput)
            }
        },
        methods: {
            search: function (input) {
                const that = this;
                this.loadingResult = true;

                EventService.search(input, function (data) {
                    console.log(data);
                    that.loadingResult = false;
                }, function (err) {
                    MessageService.error(err);
                    console.log('error:', err);
                    that.loadingResult = false;
                });
            }
        }
    }
</script>

<style scoped>

</style>
