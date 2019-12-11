<template>

  <div id="result-panel">

    <div v-if="loadingResult" class="p-grid">
      <div class="p-col-5"></div>
      <div class="p-col-2" style="text-align: center;">
        <ProgressSpinner mode="indeterminate"/>
      </div>
      <div class="p-col-5"></div>
    </div>

    <div v-if="!loadingResult" id="result">

      <div v-for="(t, i) in timeSlotCount" :key="t" class="p-grid">

        <span v-for="tt in t" :key="tt" :id="'slot-'+i+'-'+tt" :ref="'slot-'+i+'-'+tt"
              :style="{width: 100/t + '%'}" style="border: 1px groove #2c3e50; min-height: 60px"
              :data-slot="buildSlotData(i, tt)"

              @click="clickMe('slot-'+i+'-'+tt)"
        >

        </span>

      </div>

    </div>

  </div>

</template>

<script>
    import EventService from '../services/events'
    import MessageService from '../services/messages'

    const slotsIntervalInMinutes = 15;

    export default {
        name: "Result",
        props: {
            searchInput: null
        },
        data() {
            return {
                loadingResult: false,
                start: null,
                end: null,
                timeSlotCount: []
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
                this.$emit("resultLoad", true);

                EventService.search(input, function (data) {
                    that.draw(input, data);
                    that.loadingResult = false;
                    that.$emit("resultLoad", false);
                }, function (err) {
                    MessageService.error(err);
                    console.log('error:', err);
                    that.loadingResult = false;
                    that.$emit("resultLoad", false);
                });
            },
            draw(input, result) {
                let index = 0;  // todo for loop
                result = result[index];

                let getTo = function () {
                    if (result.busy.length > 0 && result.free.length > 0) {
                        return Math.max(
                            new Date(result.busy[result.busy.length - 1].end).getTime(),
                            new Date(result.free[result.free.length - 1].end).getTime(),
                        )
                    } else if (result.busy.length > 0) {
                        return result.busy[result.busy.length - 1].end
                    } else if (result.free.length > 0) {
                        return result.free[result.free.length - 1].end
                    }
                };

                let from = new Date(input.from);
                let to = new Date(getTo());
                this.timeSlotCount[index] =
                    Math.ceil(Math.floor((Math.abs(to - from) / 1000) / 60) / slotsIntervalInMinutes);

                this.start = from;
                this.end = to;
            },
            buildSlotData: function (rowId, slotId) {
                console.log(rowId, slotId, (slotId - 1) * slotsIntervalInMinutes);
                let from = new Date(this.start);
                from.setMinutes(from.getMinutes() + ((slotId - 1) * slotsIntervalInMinutes));
                from.setSeconds(0);
                return 'from=' + from.toLocaleTimeString('en-US', {hour12: false}) + ','
            },
            clickMe: function (ref) {
                console.log(this.$refs[ref][0].getAttribute("data-slot"))

            }
        }
    }
</script>

<style scoped>

</style>
