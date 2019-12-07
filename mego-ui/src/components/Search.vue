<template>
  <div>

    <div class="content-section implementation">

      <div>
        <h3>Required Attendees:</h3>
        <span class="p-fluid">
        <AutoComplete :multiple="true" v-model="selectedReqAttendees" :suggestions="filteredReqAttendees"
                      @complete="searchReqAttendees($event)" field="name">
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

      <div>
        <h3>Optional Attendees:</h3>
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
  </div>
</template>

<script>

    import AttendeesService from '../services/attendees'

    export default {
        name: "Search",
        data: function () {
            return {
                comp: this,
                selectedReqAttendees: [],
                filteredReqAttendees: null,
                selectedOptAttendees: [],
                filteredOptAttendees: null,
            }
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
