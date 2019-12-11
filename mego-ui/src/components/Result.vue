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

      <div v-for="r in rowsCount" :key="r" class="p-grid"> <!--  for each input.rooms-->

        <span v-for="t in timeSlotCount" :key="t"
              :id="'slot-'+ r +'-'+t" :ref="'slot-'+ r +'-'+t"
              class="slot" :class="'slot-' + r" :style="{width: 100/timeSlotCount + '%'}"
              :data-slot-from="buildSlotData(t)"
              @click="clickMe('slot-'+ r +'-'+t)">

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
                rowsCount: null,
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
                return from;
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
                let that = this;
                this.rowsCount = result.length;

                this.start = new Date(input.from);
                this.end = new Date(this.start);
                this.end.setHours(18);    // TODO read from server
                this.end.setMinutes(0);

                this.timeSlotCount =
                    Math.ceil(Math.floor((Math.abs(this.end - this.start) / 1000) / 60) / slotIntervalInMinutes);


                setTimeout(() => {
                    // set busy
                    for (let rowId = 0; rowId < result.length; rowId++) {

                        let roomResult = result[rowId];
                        let busyDetails = roomResult.busy_details;
                        for (let key in busyDetails) {

                            let detail = busyDetails[key];
                            detail.forEach(event => {
                                let slotIds = that.getSlotsIdsByEvent(event, rowId + 1);

                                for (let slotId of slotIds) {
                                    let div = document.createElement("div");
                                    div.classList.add(event.busy_type);
                                    document.getElementById(slotId).append(div)
                                }
                            });
                        }
                    }

                    // set style
                    let slots = document.getElementsByClassName("slot");
                    for (let slot of slots) {

                        if (new Date(slot.getAttribute("data-slot-from")).getMinutes() === 0){
                            slot.classList.add("slot-left");
                        }

                        let divs = slot.getElementsByTagName("div");
                        for (let i = 0; i < divs.length; i++) {
                            let div = divs[i];
                            div.setAttribute("style", "height: " + 100 / divs.length + "%")
                        }
                    }

                }, 10);
            },
            getSlotsIdsByEvent(event, rowId) {
                let eventStart = new Date(event.start);
                let eventEnd = new Date(event.end);

                let slotsIds = [];

                let list = document.getElementsByClassName("slot-" + rowId);
                for (let i = 0; i < list.length; i++) {
                    let slot = list[i];

                    let slotFrom = new Date(slot.getAttribute("data-slot-from"));
                    let slotTo = new Date(slot.getAttribute("data-slot-from"));
                    slotTo.setMinutes(slotTo.getMinutes() + slotIntervalInMinutes);

                    if (eventStart <= slotFrom && eventEnd >= slotTo) {
                        slotsIds.push(slot.getAttribute("id"))
                    }
                }

                return slotsIds;
            },
            clickMe: function (ref) {
                console.log(this.$refs[ref][0].getAttribute("data-slot-from"))

            }
        }
    }
</script>

<style scoped>

  .slot {
    border-top: 1px groove #2c3e50;
    border-bottom: 1px groove #2c3e50;
    height: 100px
  }

</style>
