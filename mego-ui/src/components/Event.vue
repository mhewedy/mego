<template>

  <div>

    <Dialog header="Meeting Request" :visible.sync="display" :modal="true">

      <div class="p-grid" id="req-autocomplete">
        <div class="p-col-2">Required Attendees</div>
        <div class="p-col-10">
      <span class="p-fluid">
      <AutoComplete :multiple="true" v-model="selectedReqAttendees" :disabled="true"></AutoComplete>
      </span>
        </div>
      </div>

      <div class="p-grid">
        <div class="p-col-2">Optional Attendees</div>
        <div class="p-col-10">
        <span class="p-fluid">
        <AutoComplete :multiple="true" v-model="selectedOptAttendees" :suggestions="filteredOptAttendees"
                      @complete="searchOptAttendees($event)" field="name">
           <template #item="slotProps">
            <div class="p-clearfix p-autocomplete-brand-item">
              <img v-if="slotProps.item.image && slotProps.item.image !== 'NA'"
                   alt="" :src="'data:image/png;base64,' + slotProps.item.image"/>
              <div style="display: flex">
              <span>{{slotProps.item.email_address}}</span>
              <span><b>{{slotProps.item.display_name}}</b></span>
              <span v-if="slotProps.item.title">({{slotProps.item.title}})</span>
              </div>
            </div>
          </template>
        </AutoComplete>
        </span>
        </div>
      </div>

      <div class="p-grid">
        <div class="p-col-2">Room</div>
        <div class="p-col-10">
      <span class="p-fluid">
      <AutoComplete v-model="selectedRooms" :disabled="true"></AutoComplete>
      </span>
        </div>
      </div>

      <div class="p-grid">
        <div class="p-col-2">Start time</div>
        <div class="p-col-3">
      <span class="p-fluid">
      <Calendar v-model="start" :showTime="true" hourFormat="12" :showIcon="false" :disabled="true"/>
      </span>
        </div>
        <div class="p-col-2"></div>
        <div class="p-col-1">Duration</div>
        <div class="p-col-3">
      <span class="p-fluid">
      <Spinner v-model="duration" :disabled="true"/>
      </span>
        </div>
      </div>

    </Dialog>

  </div>

</template>

<script>
    export default {
        name: "Event",
        props: {
            eventDetails: null
        },
        watch: {
            eventDetails: function () {
                this.display = this.eventDetails != null;

                if (this.eventDetails != null) {
                    this.start = this.eventDetails.start;
                    this.duration = this.eventDetails.duration;
                    this.selectedReqAttendees = this.eventDetails.emails;
                    this.selectedRooms = [];
                    this.selectedRooms.push(this.eventDetails.room);

                    setTimeout(() => {
                        let reqEmailsTimes = document.getElementById("req-autocomplete")
                            .getElementsByClassName("pi-times");
                        for (let reqEmailTimes of reqEmailsTimes) {
                            reqEmailTimes.setAttribute("style", "display: none;")
                        }

                        let reqEmailsLabel = document.getElementById("req-autocomplete")
                            .getElementsByClassName("p-autocomplete-token-label");
                        for (let reqEmailLabel of reqEmailsLabel) {
                            reqEmailLabel.setAttribute("style", "margin-right: 0em !important;")
                        }

                    }, 10)
                }
            }
        },
        data() {
            return {
                display: this.eventDetails != null,
                selectedReqAttendees: [],
                selectedRooms: [],
                selectedOptAttendees: [],
                filteredOptAttendees: null,
                start: null,
                duration: null
            }
        },
        methods: {
            searchOptAttendees: function (event) {
                console.log(event);
            }
        }
    }
</script>

<style scoped>

  .p-dialog {
    width: 65%;
  }

</style>
