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

      <div class="p-grid row">
        <span class="p-col-2">
        </span>
        <span v-for="t in timeSlotCount" :key="t" :style="{width: (80/timeSlotCount) + '%'}">
          {{buildHeaderTime(t)}}
        </span>
      </div>

      <div v-for="r in rowsCount" :key="r" class="p-grid row"> <!--  for each input.rooms-->
        <span class="p-col-2" :class="'row-'+r">
        </span>
        <span v-for="t in timeSlotCount" :key="t"
              :id="'slot-'+ r +'-'+t" :ref="'slot-'+ r +'-'+t"
              class="slot" :class="'slot-' + r" :style="{width: 80/timeSlotCount + '%'}"
              :data-slot-from="buildSlotData(t)">

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
            this.search()
        },
        watch: {
            searchInput: function () {
                this.search()
            }
        },
        methods: {
            buildSlotData: function (slotId) {
                let from = new Date(this.start);
                from.setMinutes(from.getMinutes() + ((slotId - 1) * slotIntervalInMinutes));
                from.setSeconds(0);
                return from;
            },
            buildHeaderTime: function (slotId) {
                let slotDate = this.buildSlotData(slotId);
                return slotDate.getMinutes() === 0 ? slotDate.toLocaleTimeString("en-US", {
                    hour: "2-digit",
                    minute: "2-digit"
                }) : ""
            },
            search: function () {
                const that = this;
                this.loadingResult = true;
                this.$emit("resultLoad", true);

                EventService.search(this.searchInput, function (data) {
                    that.draw(data);
                    that.loadingResult = false;
                    that.$emit("resultLoad", false);
                }, function (err) {
                    MessageService.error(err);
                    console.log('error:', err);
                    that.loadingResult = false;
                    that.$emit("resultLoad", false);
                });
            },
            draw(data) {
                let that = this;
                let result = data.room_events;
                this.rowsCount = result.length;

                this.start = new Date(this.searchInput.from);
                this.end = new Date(this.start);
                this.end.setHours(data.end_of_day_hours);
                this.end.setMinutes(0);

                this.timeSlotCount =
                    Math.ceil(Math.floor((Math.abs(this.end - this.start) / 1000) / 60) / slotIntervalInMinutes);


                setTimeout(() => {
                    // set busy
                    for (let rowId = 0; rowId < result.length; rowId++) {

                        let roomResult = result[rowId];
                        document.getElementsByClassName("row-" + (rowId + 1))[0].innerText = roomResult.room_name;

                        let busyDetails = roomResult.busy_details;
                        for (let key in busyDetails) {

                            let detail = busyDetails[key];
                            detail.forEach(event => {
                                let slotIds = that.getSlotsIdsByEvent(event, rowId + 1);

                                for (let slotId of slotIds) {
                                    let div = document.createElement("div");

                                    if (key === roomResult.room && event.busy_type === "Busy") {
                                        div.setAttribute("title", key + "(RoomBusy)");
                                        div.classList.add("RoomBusy");
                                    } else {
                                        div.setAttribute("title", key + "(" + event.busy_type + ")");
                                        div.classList.add(event.busy_type);
                                    }
                                    document.getElementById(slotId).append(div)
                                }
                            });
                        }
                    }
                    // set style
                    let slots = document.getElementsByClassName("slot");
                    for (let i = 0; i < slots.length; i++) {
                        let slot = slots[i];

                        // style slot
                        if (i === 0) slot.classList.add("slot-left");
                        if (i === slots.length - 1) slot.classList.add("slot-right");

                        if (new Date(slot.getAttribute("data-slot-from")).getMinutes() === 0) {
                            slot.classList.add("slot-left");
                        }

                        this.handleEvents(slots, i);

                        let divs = slot.getElementsByTagName("div");
                        // style divs, set height
                        for (let j = 0; j < divs.length; j++) {
                            let div = divs[j];
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
            handleEvents: function (slots, i) {
                let slot = slots[i];
                let divs = slot.getElementsByTagName("div");

                let clickHandler = (evt) => {
                    console.log(evt.target, this.searchInput)
                };

                slot.addEventListener("mousemove", () => {
                    if (divs.length === 0) {
                        let numSlots = this.searchInput.duration / slotIntervalInMinutes;

                        let truth = [];
                        for (let x = 0; x < numSlots; x++) {
                            if (i + x < slots.length &&
                                slots[i + x].getElementsByTagName("div").length === 0) {
                                truth.push(true)
                            } else {
                                truth.push(false)
                            }
                        }

                        if (truth.every(it => it === true)) {
                            for (let x = 0; x < numSlots; x++) {
                                slots[i + x].style.backgroundColor = "#ffcc00";
                                slots[i + x].style.cursor = "pointer";
                            }
                            slot.addEventListener("click", clickHandler)
                        }
                    }
                });
                slot.addEventListener("mouseout", () => {
                    if (divs.length === 0) {
                        for (let s of slots) {
                            s.style.backgroundColor = "transparent";
                            s.style.cursor = "not-allowed";
                        }
                        slot.removeEventListener("click", clickHandler);
                    }
                });
            }
        }
    }
</script>

<style scoped>

  .slot {
    border-top: 1px groove #2c3e50;
    border-bottom: 1px groove #2c3e50;
    height: 150px;
    cursor: not-allowed;
  }

  .row {
    padding-bottom: 7px;
  }

</style>
