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

      <div class="p-grid"> <!--  for each input.rooms-->

        <span v-for="t in timeSlotCount" :key="t"
              :id="'slot-'+t" :ref="'slot-'+t"
              class="slot" :style="{width: 100/timeSlotCount + '%'}"
              :data-slot-from="buildSlotData(t)"
              @click="clickMe('slot-'+t)">

        </span>
      </div>

    </div>

  </div>

</template>

<script>
    import EventService from '../services/events'
    import MessageService from '../services/messages'

    const slotIntervalInMinutes = 15;

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
                timeSlotCount: null
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
            buildSlotData: function (slotId) {
                let from = new Date(this.start);
                from.setMinutes(from.getMinutes() + ((slotId - 1) * slotIntervalInMinutes));
                from.setSeconds(0);
                return from.toLocaleTimeString('en-US', {hour12: false});
            },
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

                this.start = new Date(input.from);
                this.end = new Date(this.start);
                this.end.setHours(18);
                this.end.setMinutes(0);

                this.timeSlotCount =
                    Math.ceil(Math.floor((Math.abs(this.end - this.start) / 1000) / 60) / slotIntervalInMinutes);

                console.log(this.start)
                console.log(this.end)
                console.log(this.timeSlotCount)

                // start for loop

                let index = 0;  // todo for loop
                result = result[index];

                // TODO
                // call getSlotIdsByEvent for each result.busy_details object
                /*"example@mhewedy.onmicrosoft.com": [
                    {
                        "start": "2019-12-11T11:15:00+03:00",
                        "end": "2019-12-11T11:45:00+03:00",
                        "busy_type": "Busy"
                    },
                    {
                        "start": "2019-12-11T12:30:00+03:00",
                        "end": "2019-12-11T13:30:00+03:00",
                        "busy_type": "Busy"
                    }
                ],*/
            },
            getSlotsIdsByEvent(event) {
                let eventStart = new Date(event.start).toLocaleTimeString('en-US', {hour12: false});
                let eventEnd = new Date(event.end).toLocaleTimeString('en-US', {hour12: false});

                let slotsIds = [];

                document.querySelector(".slot").forEach(it => {
                    let slotFrom = it.getAttribute("data-slot-from");
                    let slotTo = new Date(slotFrom);
                    slotTo.setMinutes(slotTo.getMinutes() + slotIntervalInMinutes);

                    if (eventStart <= slotFrom && eventEnd >= slotTo ) {
                        slotsIds.push(it.getAttribute("id"))
                    }
                });

                return slotsIds;
            },
            clickMe: function (ref) {
                console.log(this.$refs[ref][0].getAttribute("data-slot-from"))

            }
        }
    }
</script>

<style scoped>

  .slot{
    border: 1px groove #2c3e50; min-height: 60px
  }
</style>
