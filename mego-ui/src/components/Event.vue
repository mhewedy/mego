<template>

  <div>

    <Dialog header="Meeting Request" :visible.sync="display" :modal="true">

      <div v-if="isSending">
        <ProgressBar mode="indeterminate" style="height: .5em; margin: -14px 0 15px 0;"/>
      </div>

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

      <div class="p-grid">
        <div class="p-col-2">Subject</div>
        <div class="p-col-10">
          <span class="p-fluid">
          <InputText v-model="subject"></InputText>
          </span>
        </div>
      </div>

      <div class="p-grid">
        <div class="p-col-12">
          <span class="p-fluid">
          <Textarea v-model="body" :autoResize="true" rows="10"></Textarea>
          </span>
        </div>
      </div>

      <template #footer>
        <Button label="Send" icon="pi pi-external-link"
                class="p-button-raised p-button-info" @click="createEvent()"
                :disabled="isSending"/>
      </template>

    </Dialog>

  </div>

</template>

<script>

    import EventService from '../services/events'
    import MessagesService from '../services/messages'
    import AttendeesService from '../services/attendees'

    export default {
        name: "Event",
        props: {
            eventDetails: null
        },
        data() {
            return {
                display: this.eventDetails != null,
                selectedReqAttendees: [],
                selectedRooms: [],
                selectedOptAttendees: [],
                filteredOptAttendees: null,
                start: null,
                duration: null,
                subject: null,
                body: null,
                isSending: false
            }
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
        methods: {
            searchOptAttendees: function (event) {

                let toExclude = this.selectedOptAttendees.map(it => it.email_address);
                toExclude.push(...this.selectedReqAttendees);

                AttendeesService.search(event.query, toExclude,
                    data => {
                        this.filteredOptAttendees = data.map(it => {
                            it["name"] = it.email_address;
                            return it
                        });

                        this.filteredOptAttendees.map(it => {
                            if (!it.image) {
                                AttendeesService.getPhoto(it.email_address, function (data) {
                                    it.image = data.base64
                                })
                            }
                        })

                    }, function (err) {
                        // eslint-disable-next-line
                        console.log(err)
                    })
            },
            createEvent: function () {

                let opts = [];
                opts.push(...this.selectedOptAttendees.map(it => it.email_address));

                let input = {
                    to: this.selectedReqAttendees,
                    optional: opts,
                    subject: this.subject,
                    body: this.body,
                    room: this.selectedRooms[0],
                    from: this.start,
                    duration: this.duration
                };

                this.isSending = true;

                EventService.create(input, () => {
                    this.isSending = false;
                    this.display = false;
                    MessagesService.success("Event created successfully!");
                    window.scrollTo(0, 0);
                    this.resetInput();
                    setTimeout(() => this.$parent.search(), 500);
                }, (error) => {
                    console.log(error);
                    this.isSending = false;
                    MessagesService.error(error.response.data.error);
                    window.scrollTo(0, 0);
                    this.resetInput();
                })
            },
            resetInput: function () {
                this.subject = null;
                this.body = null;
                this.selectedOptAttendees = [];
                this.filteredOptAttendees = null;
            }
        }
    }
</script>

<style scoped>

  .p-dialog {
    width: 65%;
  }

</style>
