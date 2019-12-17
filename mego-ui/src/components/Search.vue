<template>
  <div id="search-panel" style="padding-top: 25px">

    <div class="p-grid">
      <div class="p-col-1">Attendees</div>
      <div class="p-col-11">
      <span class="p-fluid">
      <AutoComplete :multiple="true" v-model="selectedReqAttendees" :suggestions="filteredReqAttendees"
                    @complete="searchReqAttendees($event)" @select="getAttendeeDetails()" field="name">
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
      <div class="p-col-1">Rooms</div>
      <div class="p-col-11">
      <span class="p-fluid">
      <Tree :value="roomsTree" selectionMode="checkbox" :selectionKeys.sync="selectedRooms"></Tree>
    </span>
      </div>
    </div>

    <div class="p-grid">
      <div class="p-col-1">Start time</div>
      <div class="p-col-3">
      <span class="p-fluid">
      <Calendar v-model="startTime" :showTime="true" hourFormat="12" :showIcon="true"
                :showButtonBar="true" :stepMinute="30" :manualInput="false"/>
      </span>
      </div>
      <div class="p-col-6"></div>
    </div>

    <div class="p-grid">
      <div class="p-col-3"></div>
      <div class="p-col-6 p-fluid" style="margin: 40px 0 20px 0;">
        <Button label="Search" icon="pi pi-search"
                class="p-button-raised p-button-rounded p-button-info" @click="search()"
                :disabled="isResultLoading"
        />
      </div>
      <div class="p-col-3"></div>
    </div>

  </div>

</template>

<script>

    import AttendeesService from '../services/attendees'
    import RoomsService from '../services/rooms'
    import MessageService from '../services/messages'

    export default {
        name: "Search",
        props: {
            isResultLoading: null
        },
        data: function () {
            return {
                selectedReqAttendees: [],
                filteredReqAttendees: null,
                roomsTree: null,
                roomsList: null,
                selectedRooms: null,
                startTime: this.getNextMeetingTime(),
                input: null
            }
        },
        mounted() {
            RoomsService.tree(it => this.roomsTree = it);
            RoomsService.list(it => this.roomsList = it);
        },
        methods: {
            searchReqAttendees: function (event) {
                let toExclude = this.selectedReqAttendees.map(it => it.email_address);
                AttendeesService.search(event.query, toExclude,
                    (data) => {
                        this.filteredReqAttendees = data.map(it => {
                            it["name"] = it.email_address;
                            return it
                        });

                    }, function (err) {
                        // eslint-disable-next-line
                        console.log(err)
                    });
            },
            getAttendeeDetails: function () {
                setTimeout(() => {
                    let names = document.getElementsByClassName("p-autocomplete-token-label");
                    for (let n of names) {
                        n.style.cursor = 'pointer';
                        n.onclick = (e) => {
                            AttendeesService.getDetails(e.target.innerText,
                                (data) => {
                                    console.log(data);
                                },
                                (err) => {
                                    MessageService.error(err);
                                });
                        };
                    }
                }, 10)
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
                let emails = [];
                let input = {
                    rooms: rooms,
                    emails: emails,
                    from: this.startTime.toISOString()
                };

                for (const key in this.selectedRooms) {
                    let value = this.selectedRooms[key];
                    if (value.checked) {
                        if (this.roomsList.indexOf(key) > -1) {
                            rooms.push(key)
                        }
                    }
                }

                rooms.sort((r1, r2) => this.roomsList.indexOf(r1) - this.roomsList.indexOf(r2));

                emails.push(...this.selectedReqAttendees.map(it => it.email_address));

                if (this.validate(input)) {
                    this.input = input;
                    this.$emit("search", input);
                }
            },
            validate: function (input) {
                if (input.emails.length === 0) {
                    MessageService.error('One attendee is required at least');
                    return false
                }
                if (input.rooms.length === 0) {
                    MessageService.error('Choose one meeting room at least');
                    return false
                }
                return true
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


  .p-treenode-content.p-highlight {
    background-color: #ffffff !important;
    color: #333333 !important;
  }

  body .p-tree .p-tree-container .p-treenode .p-treenode-content.p-highlight .p-tree-toggler {
    color: #848484 !important;
  }

  .p-treenode-children {
    padding-left: 40px !important;
  }
</style>
