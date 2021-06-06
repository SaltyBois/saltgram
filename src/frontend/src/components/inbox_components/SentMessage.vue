<template>
  <div class="inbox-message-div-right">
    <div class="inbox-message-sent">
      Sent message template
    </div>
    <div class="message-time-sent">
      15.4.2015. 17:46 time template
    </div>
  </div>
</template>

<script>
export default {
name: "SentMessage"
}
</script>

<style scoped>

  .inbox-message-div-right {
    display: inline-block;
    height: auto;
    width: 100%;
    align-content: center;
    margin: 3px 0;
    /*padding: 20px;*/
  }

  .inbox-message-sent {
    background-color: #016ddb;
    color: white;
    border: #323232 2px solid;
    float: right;
    max-width: 60%;
    overflow-wrap: anywhere;
    padding: 5px;
    margin-right: 5px;

    border-start-end-radius: 10px 10px;
    border-end-end-radius: 10px 10px;
    border-start-start-radius: 10px 10px;
    border-end-start-radius: 10px 10px;
  }

  .message-time-sent {
    float: right;
    text-align: right;
    font: 12px bold;
    padding-top: 10px;
    padding-right: 5px;
  }

</style>