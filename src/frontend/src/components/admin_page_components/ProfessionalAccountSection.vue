<template>
  <div style="height: 100%">
    <div class="title">
      <h2>Professional Account Applications Section</h2>
    </div>
    <div style="display: inline-flex; height: 90%; width: 100%">
      <div class="dropdown-list">
        <div class="inner-layout">

        </div>
      </div>
      <div class="application-view">

      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "ProfessionalAccountSection"
}
</script>

<style scoped>

.title {
  height: 10%;
  text-align: center;
}

.dropdown-list {
  height: 100%;
  width: 40%;
  background-color: red;
  overflow-y: auto;
}

.inner-layout {
  margin: 5px;

  border: 1px solid black;
  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;
}

.application-view {
  height: 100%;
  width: 60%;
  background-color: #36c400;
}

</style>