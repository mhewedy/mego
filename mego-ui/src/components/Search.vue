<template>
  <div class="container">

    <hr/>
    <h3>Meeting Organizer</h3>
    <hr/>

    <div class="row">
      <div class="input-group mb-3">
        <div class="input-group-prepend">
          <span class="input-group-text" id="basic-addon1">Required Attendees</span>
        </div>
        <textarea class="form-control" aria-label="With textarea"></textarea>
      </div>
    </div>
    <div class="row">
      <div class="input-group mb-3">
        <div class="input-group-prepend">
          <span class="input-group-text" id="basic-addon2">Optional Attendees</span>
        </div>
        <textarea class="form-control" aria-label="With textarea"></textarea>
      </div>
    </div>

    <div class="row">
      <div class="input-group mb-3">
        <div class="input-group-prepend">
          <span class="input-group-text" id="basic-addon3">Rooms</span>
        </div>
        <textarea class="form-control" aria-label="With textarea" placeholder="List of candidate rooms"
                  :disabled="allRooms"></textarea>
        <div class="input-group-append">
          <div class="input-group-text">
            <div class="form-check">
              <input class="form-check-input adjusted-checkbox" type="checkbox" id="check1"
                     v-on:change="toggleAllRooms()" :checked="allRooms">
              <label class="form-check-label" for="check1">
                All Rooms
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="row">
      <div class="input-group mb-3">
        <div class="input-group-prepend">
          <span class="input-group-text">Duration</span>
        </div>
        <input type="number" class="form-control" aria-label="Duration" value="30">
        <div class="input-group-append">
          <span class="input-group-text">Minutes</span>
        </div>
      </div>
    </div>

    <div class="row">
      <div class="input-group mb-3">
        <div class="input-group-prepend">
          <span class="input-group-text" id="basic-addon4">Start time</span>
        </div>
        <input type="datetime-local" class="form-control" :disabled="startsNow">

        <div class="input-group-append">
          <div class="input-group-text">
            <div class="form-check">
              <input class="form-check-input adjusted-checkbox" type="checkbox" id="check2"
                     v-on:change="toggleStartsNow()" :checked="startsNow">
              <label class="form-check-label" for="check1">
                Starts immediately
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="row">
      <button type="button" class="btn btn-secondary btn-lg btn-block" v-on:click="search">Search</button>
    </div>

  </div>
</template>

<script>
    import axios from 'axios'

    export default {
        name: "Search",
        data: function () {
            return {
                allRooms: true,
                startsNow: true
            }
        },
        methods: {
            toggleAllRooms: function () {
                this.allRooms = !this.allRooms
            },
            toggleStartsNow: function () {
                this.startsNow = !this.startsNow
            },
            search: function () {
                axios.get("/api/test")
                    .then(resp => alert(resp.data["key"]))
            }
        }
    }
</script>

<style scoped>

  textarea {
    height: 38px !important;
  }

  .input-group-text {
    width: 170px;
  }

  .adjusted-checkbox {
    margin-top: 5px !important;
  }
</style>
