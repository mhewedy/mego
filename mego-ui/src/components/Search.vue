<template>
  <div id="search-panel">

    <div class="p-grid">
      <div class="p-col-2">Required Attendees</div>
      <div class="p-col-10">
          <span class="p-fluid">
            <AutoComplete :multiple="true" v-model="selectedReqAttendees" :suggestions="filteredReqAttendees"
                          @complete="searchReqAttendees($event)" field="name">
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
      <div class="p-col-2">Optional Attendees</div>
      <div class="p-col-10">
          <span class="p-fluid">
            <AutoComplete :multiple="true" v-model="selectedOptAttendees" :suggestions="filteredOptAttendees"
                          @complete="searchOptAttendees($event)" field="name">
               <template #item="slotProps">
                    <div class="p-clearfix p-autocomplete-brand-item">
                      <img v-if="slotProps.item.image" alt="" :src="'data:image/png;base64,' + slotProps.item.image"/>
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
      <div class="p-col-2">Rooms</div>
      <div class="p-col-10">
          <span class="p-fluid">
            <Tree :value="roomsTree" selectionMode="checkbox" :selectionKeys.sync="selectedRooms"></Tree>
        </span>
      </div>
    </div>

    <div class="p-grid">
      <div class="p-col-2">Start time</div>
      <div class="p-col-3">
          <span class="p-fluid">
            <Calendar v-model="startTime" :showTime="true" hourFormat="12" :showIcon="true"
                      :showButtonBar="true" :stepMinute="30" :manualInput="false"/>
          </span>
      </div>
      <div class="p-col-2"></div>
      <div class="p-col-1">Duration</div>
      <div class="p-col-3">
          <span class="p-fluid">
            <Spinner v-model="duration" :step="30" :min="30"/>
          </span>
      </div>
    </div>

    <div class="p-grid">
      <div class="p-col-3"></div>
      <div class="p-col-6 p-fluid">
        <Button label="Search" icon="pi pi-search"
                class="p-button-raised p-button-rounded p-button-info" @click="search()"/>
      </div>
      <div class="p-col-3"></div>
    </div>

  </div>

</template>

<script>

    import AttendeesService from '../services/attendees'
    import RoomsService from '../services/rooms'

    export default {
        name: "Search",
        data: function () {
            return {
                selectedReqAttendees: [],
                filteredReqAttendees: null,
                selectedOptAttendees: [],
                filteredOptAttendees: null,
                roomsTree: null,
                roomsList: null,
                selectedRooms: null,
                startTime: this.getNextMeetingTime(),
                duration: 30
            }
        },
        mounted() {
            RoomsService.tree(it => this.roomsTree = it);
            RoomsService.list(it => this.roomsList = it);
        },
        methods: {
            searchReqAttendees: function (event) {
                const that = this;
                AttendeesService.search(event.query,
                    function (data) {
                        that.filteredReqAttendees = data.map(it => {
                            it["name"] = it.email_address;
                            return it
                        });

                        that.filteredReqAttendees.map(it => {
                            if (!it.image) {
                                AttendeesService.getPhoto(it.email_address, function (data) {
                                    it.image = data.base64
                                })
                            }
                        })

                    }, function (err) {
                        console.log(err)
                    })
            },
            searchOptAttendees: function (event) {
                const that = this;
                AttendeesService.search(event.query,
                    function (data) {
                        that.filteredOptAttendees = data.map(it => {
                            it["name"] = it.email_address;
                            return it
                        });

                        that.filteredOptAttendees.map(it => {
                            if (!it.image) {
                                AttendeesService.getPhoto(it.email_address, function (data) {
                                    it.image = data.base64
                                })
                            }
                        })

                    }, function (err) {
                        console.log(err)
                    })
            },
            getNextMeetingTime: function () {
                let date = new Date();
                if (date.getMinutes() > 0 && date.getMinutes() < 30) {
                    date.setMinutes(30);
                } else if (date.getMinutes() > 30) {
                    date.setMinutes(0);
                    date.setHours(date.getHours() + 1);
                }
                return date
            },
            search: function () {

                let rooms = [];
                let mails = [];
                let input = {
                    rooms: rooms,
                    mails: mails,
                    from: this.startTime.toISOString(),
                    duration: this.duration
                };

                for (const key in this.selectedRooms) {
                    let value = this.selectedRooms[key];
                    if (value.checked) {
                        if (this.roomsList.indexOf(key) > -1) {
                            rooms.push(key)
                        }
                    }
                }

                mails.push(...this.selectedReqAttendees.map(it => it.email_address));
                mails.push(...this.selectedOptAttendees.map(it => it.email_address));

                console.log(input)
            }
        }
    }
</script>

<style lang="scss">
  .p-autocomplete-brand-item {
    img {
      width: 32px;
      display: inline-block;
      float: right;
      margin: 5px 0 2px 5px;
    }

    span {
      font-size: 16px;
      margin: 10px 10px 0 0;
    }
  }
</style>
