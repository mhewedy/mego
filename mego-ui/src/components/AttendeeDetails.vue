<template>

  <div>
    <div v-if="!attendeeDetails">
      <ProgressSpinner style="width:20px;height:20px" strokeWidth="3" animationDuration=".5s"/>
    </div>

    <div v-if="attendeeDetails">
      <div class="p-grid nested-grid">
        <div class="p-col-3">
          <img v-if="attendeeDetails.image && attendeeDetails.image !== 'NA'"
               alt="" :src="'data:image/png;base64,' + attendeeDetails.image"/>
        </div>
        <div class="p-col-9">
          <div class="p-grid">
            <div class="p-col-12" style="font-weight: bold">
              {{attendeeDetails.display_name}}
            </div>
            <div class="p-col-12" v-if="attendeeDetails.title">
              {{attendeeDetails.title}}
            </div>
            <div class="p-col-12" v-if="attendeeDetails.department">
              {{attendeeDetails.department}}
            </div>
          </div>
        </div>
        <div class="p-col-12">
        </div>

        <div class="p-col-12">
        </div>
        <div class="p-col-6" v-if="attendeeDetails.business_phone_numbers">
          Phone: {{attendeeDetails.business_phone_numbers}}
        </div>
        <div class="p-col-6" v-if="attendeeDetails.mobile_phone">
          Mobile: {{attendeeDetails.mobile_phone}}
        </div>
        <div class="p-col-12">
        </div>
        <div class="p-col-12" v-if="attendeeDetails.office_location">
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
