<template>

  <div id="result-panel">

    <div v-if="loadingResult" class="p-grid">
      <div class="p-col-5"></div>
      <div class="p-col-2" style="text-align: center;">
        <ProgressSpinner mode="indeterminate"/>
      </div>
      <div class="p-col-5"></div>
    </div>

    <div v-if="!loadingResult" id="result" style="padding-bottom: 50px">

      <div id="result-duration" class="p-grid" style="padding-bottom: 20px">
        <div class="p-col-12">
          <div style="float: right">
            <span style="font-weight: bold">Duration</span>
            <Spinner v-model="duration" :step="30" :min="30" :readonly="true" style="padding-left: 10px"/>
          </div>
        </div>
      </div>

      <div class="p-grid row">
        <span class="p-col-1">
        </span>
        <span v-for="t in timeSlotCount" :key="t" :style="{width: (90/timeSlotCount) + '%'}">
          {{buildHeaderTime(t)}}
        </span>
      </div>

      <div v-for="r in rowsCount" :key="r" class="p-grid row"> <!--  for each input.rooms-->
        <span class="p-col-1" :class="'row-'+r">
        </span>
        <span v-for="t in timeSlotCount" :key="t"
              :id="'slot-'+ r +'-'+t" :ref="'slot-'+ r +'-'+t"
              class="slot" :class="'slot-' + r" :style="{width: 90/timeSlotCount + '%'}"
              :data-slot-from="buildSlotData(t)">

        </span>
      </div>

      <div style="padding-top: 40px" class="p-grid">

        <span class="Busy p-col-1 legend-block"></span>
        <span class="p-col-1 legend-key">Busy</span>

        <span class="RoomBusy p-col-1 legend-block"></span>
        <span class="p-col-1 legend-key">Busy Room</span>

        <span class="Tentative p-col-1 legend-block"></span>
        <span class="p-col-1 legend-key">Tentative</span>

        <span class="error p-col-1 legend-block"></span>
        <span class="p-col-1 legend-key">Error</span>
      </div>

    </div>

    <Event :eventDetails="eventDetails"></Event>

  </div>

</template>

<script>
    import EventService from '../services/events'
    import MessageService from '../services/messages'
    import Event from "./Event";

    const slotIntervalInMinutes = 15;

    export default {
        name: "Result",
        components: {Event},
        props: {
            searchInput: null
        },
        data() {
            return {
                loadingResult: false,
                start: null,
                end: null,
                rowsCount: null,
                timeSlotCount: null,
                eventDetails: null,
                duration: 30,
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
                this.loadingResult = true;
                this.$emit("resultLoad", true);

                this.start = this.calcStartDate();
                this.end = this.searchInput.to;

                EventService.search({
                    rooms: this.searchInput.rooms,
                    emails: this.searchInput.emails,
                    from: this.start,
                    to: this.end
                }, data => {
                    this.draw(data);
                    this.loadingResult = false;
                    this.$emit("resultLoad", false);
                }, err => {
                    MessageService.error(err);
                    this.loadingResult = false;
                    this.$emit("resultLoad", false);
                });
            },
            draw(result) {
                this.rowsCount = result.length;

                this.timeSlotCount =
                    Math.ceil(Math.floor((Math.abs(this.end - this.start) / 1000) / 60) / slotIntervalInMinutes);

                setTimeout(() => {
                    // set busy
                    for (let rowId = 0; rowId < result.length; rowId++) {

                        let roomResult = result[rowId];
                        document.getElementsByClassName("row-" + (rowId + 1))[0].innerText = roomResult.room_name;

                        if (roomResult.error) {
                            // set style
                            let slots = document.getElementsByClassName("slot-" + (rowId + 1));
                            for (let i = 0; i < slots.length; i++) {
                                let slot = slots[i];

                                // style slot
                                if (i === 0) slot.classList.add("slot-left");
                                if (i === slots.length - 1) slot.classList.add("slot-right");

                                slot.classList.add("error");

                                slot.classList.add("tooltip");

                                let spanTooltip = document.createElement("span");
                                spanTooltip.innerText = roomResult.error;
                                spanTooltip.classList.add("tooltiptext");
                                spanTooltip.classList.add("wide");
                                slot.append(spanTooltip);
                            }

                        } else {
                            let busyDetails = roomResult.busy_details;
                            for (let key in busyDetails) {

                                let detail = busyDetails[key];
                                detail.forEach(event => {
                                    let slotIds = this.getSlotsIdsByEvent(event, rowId + 1);

                                    for (let slotId of slotIds) {
                                        let div = document.createElement("div");

                                        let tooltip = key + "(" + event.busy_type + ")";
                                        let busyType = event.busy_type;

                                        if (key === roomResult.room && event.busy_type === "Busy") {
                                            tooltip = key + "(Busy Room)";
                                            busyType = 'RoomBusy'
                                        }

                                        div.classList.add(busyType);
                                        div.classList.add("tooltip");

                                        let spanTooltip = document.createElement("span");
                                        spanTooltip.innerText = tooltip;
                                        spanTooltip.classList.add("tooltiptext");
                                        div.append(spanTooltip);

                                        document.getElementById(slotId).append(div)

                                    }
                                });
                            }

                            // set style
                            let slots = document.getElementsByClassName("slot-" + (rowId + 1));
                            for (let i = 0; i < slots.length; i++) {
                                let slot = slots[i];

                                this.setSlotBGColor(slot);

                                // style slot
                                if (i === 0) slot.classList.add("slot-left");
                                if (i === slots.length - 1) slot.classList.add("slot-right");

                                if (new Date(slot.getAttribute("data-slot-from")).getMinutes() === 0) {
                                    slot.classList.add("slot-left");
                                }

                                this.handleEvents(slots, i);

                                // style divs, set height
                                let divs = slot.getElementsByTagName("div");
                                for (let j = 0; j < divs.length; j++) {
                                    let div = divs[j];
                                    div.setAttribute("style", "height: " + 100 / (Math.max(6, divs.length))
                                        + "%; border-bottom: white 3px solid;")
                                }
                            }
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
                    let eventDetails = {};

                    let rowId = /-([0-9]+)-/.exec(evt.target.getAttribute("id"))[1];

                    eventDetails.duration = this.duration;
                    eventDetails.emails = this.searchInput.emails;
                    eventDetails.start = new Date(evt.target.getAttribute("data-slot-from"));
                    eventDetails.room = this.searchInput.rooms[rowId - 1];

                    this.eventDetails = eventDetails;
                };

                slot.addEventListener("mousemove", () => {
                    if (divs.length === 0) {    // no sub divs added (no events)
                        let numSlots = this.duration / slotIntervalInMinutes;

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
                            this.setSlotBGColor(s);
                            s.style.cursor = "not-allowed";
                        }
                        slot.removeEventListener("click", clickHandler);
                    }
                });
            },
            setSlotBGColor: function (slot) {
                let slotTo  = new Date(slot.getAttribute("data-slot-from"));
                slotTo.setMinutes(slotTo.getMinutes() + slotIntervalInMinutes);
                if (slotTo <= new Date(this.searchInput.from)) {
                    slot.style.background = 'rgba(222,224,227,0.4)'
                }else {
                    slot.style.background = "transparent";
                }
            },
            calcStartDate: function () {
                let start = new Date(this.searchInput.from);
                start.setHours(8);
                start.setMinutes(0);
                return start;
            }
        }
    }
</script>

<style scoped>

  .slot {
    border-top: 1px groove #2c3e50;
    border-bottom: 1px groove #2c3e50;
    height: 100px;
    cursor: not-allowed;
  }

  .row {
    padding-bottom: 7px;
  }

</style>

<!--not scoped-->
<style>

  /**/
  .Busy {
    color: white;
    background-color: #0057e7;
  }

  .Tentative {
    color: white;
    /*background: repeating-linear-gradient(45deg, #606dbc, #606dbc 10px, #deedf8 10px, #deedf8 20px);*/
    background-color: #008744;
  }

  .RoomBusy {
    background-color: #ffa700;
  }

  .error {
    background-color: #d62d20;
  }

  /**/

  .slot-left {
    border-left: 1px groove #2c3e50;
  }

  .slot-right {
    border-right: 1px groove #2c3e50;
  }

  .tooltip {
    position: relative;
    /*display: inline-block;*/
    border-bottom: 1px dotted black;
  }

  .tooltip .tooltiptext {
    /*visibility: hidden;*/
    /*width: 120px;*/
    background-color: black;
    color: #fff;
    display: none;
    /*text-align: center;*/
    border-radius: 4px;
    padding: 5px 5px;

    /* Position the tooltip */
    position: absolute;
    z-index: 1;
    margin: 20px 0 0 40px;
  }

  .wide {
    min-width: 500px;
  }

  .tooltip:hover .tooltiptext {
    /*visibility: visible;*/
    display: inline-block;
  }

  #result-duration .p-spinner-input {
    font-weight: bold;
    width: 70px !important;
  }

  .legend-key {
    margin-top: -7px;
    width: 7% !important;
  }

  .legend-block {
    width: 20px !important;
    height: 20px !important;
  }

</style>
