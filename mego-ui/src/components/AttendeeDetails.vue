<template>

  <div>
    <div v-if="!attendeeDetails">
      <ProgressSpinner style="width:20px;height:20px" strokeWidth="3" animationDuration=".5s"/>
    </div>

    <div v-if="attendeeDetails">
      <div class="p-grid nested-grid">
        <div class="p-col-2">
          <img v-if="attendeeDetails.image && attendeeDetails.image !== 'NA'"
               alt="" :src="'data:image/png;base64,' + attendeeDetails.image" style="width: 50px;"/>
        </div>
        <div class="p-col-10">
          <div class="p-grid">
            <div class="p-col-12" style="font-weight: bold">
              {{attendeeDetails.display_name}}
            </div>
            <div class="p-col-12" v-if="attendeeDetails.title" style="margin-top: -5px">
              {{attendeeDetails.title}}
            </div>
            <div class="p-col-12" v-if="attendeeDetails.department" style="margin-top: -5px">
              {{attendeeDetails.department}}
            </div>
          </div>
        </div>
        <div class="p-col-12" style="margin-top: -5px">
        </div>

        <div class="p-col-12" style="margin-top: -5px">
        </div>
        <div class="p-col-6" v-if="attendeeDetails.business_phone_numbers" style="margin-top: -5px">
          Phone: {{attendeeDetails.business_phone_numbers}}
        </div>
        <div class="p-col-6" v-if="attendeeDetails.mobile_phone" style="margin-top: -5px">
          Mobile: {{attendeeDetails.mobile_phone}}
        </div>
        <div class="p-col-12" style="margin-top: -5px">
        </div>
        <div class="p-col-12" v-if="attendeeDetails.office_location" style="margin-top: -5px">
          {{attendeeDetails.office_location}}
        </div>
      </div>
    </div>
  </div>

</template>

<script>
    export default {
        name: "AttendeeDetails",
        props: {
            attendeeDetails: null
        }
    }
</script>

<style scoped>

</style>
