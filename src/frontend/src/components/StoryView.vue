<template>
  <div>
    <transition name="fade" appear>
      <v-layout class="modal-overlay"
                v-if="visible"
                @click="toggleView">
        <div class="top-left-corner">
          <h3 id="home-title">Saltgram</h3>
        </div>
      </v-layout>
    </transition>
    <transition name="fade" appear>
      <v-layout class="modal"
                v-if="visible"
                justify-center
                align-center
                column>
        <div class="circle">
          <b style="font-size: 19px">{{timerSeconds}}</b>
        </div>
        <div class="circle-top-right">
          <div class="mt-3">
            <b style="font-size: 14px" >1/1</b>
          </div>
        </div>
        <div class="left-btn">
          <i class="fa fa-sign-out ml-2 mt-3" style="transform: scale(1.4)"/>
        </div>
        <v-img class="image" :src="imageSrc" />
        <div class="viewers">
          <h3>
            Seen by 503
          </h3>
        </div>
        <div class="close-friends-div">
          <h3>
            CLOSE FRIENDS
          </h3>
        </div>
        <div class="right-btn">
          <i class="fa fa-sign-out mt-2 ml-1" style="transform: scale(1.4)" />
        </div>
      </v-layout>
    </transition>
  </div>
</template>

<script>
export default {
  name: "StoryView",
  data: function() {
    return {
      visible: false,
      timerSeconds: 10,
    }
  },
  props: {
    imageSrc: {
      type: String,
      required: true
    }
  },
  methods: {
    toggleView() {
      this.visible = !this.visible
      if (this.visible) {
        this.timerSeconds = 10
        this.countDownTimer()
      }
    },
    countDownTimer() {
      if(this.timerSeconds > 0) {
        setTimeout(() => {
          this.timerSeconds -= 1
          this.countDownTimer()
        }, 1000)
      }
      else {
        this.visible = false;
      }
    }
  },
  computed() {
    this.countDownTimer()
  },
}
</script>

<style scoped>

.modal {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 99;

  width: 90vh;
  height: 90vh;
  background-color: #FFF;
  border-radius: 16px;

  padding: 5px;

  overflow-y: auto;

  object-fit: cover;
}

.top-left-corner {
  position: relative;
  top: 30px;
  left: 0;
  width: 300px;
  height: 100px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 98;
  background-color: rgba(0, 0, 0, 0.8);
}

.image {
  object-fit: contain;
  width: 80%;
  height: 90%;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform .5s;
}

.slide-enter,
.slide-leave-to {
  transform: translateY(-50%) translateX(100vw);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter,
.fade-leave-to {
  opacity: 0;
}

#home-title {
  font-size: 30px;
  font-family: "Lucida Handwriting", cursive;
  text-transform: capitalize;
  color: #FFFFFF;
}

.circle {
  position: absolute;
  top: 15px;
  left: 15px;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  border: 1px black solid;
  background-color: black;
  color: #FFFFFF;
  text-align: center;
  justify-content: center;
}

.circle-top-right {
  position: absolute;
  top: 15px;
  right: 10px;
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: 1px black solid;
  background-color: black;
  color: #FFFFFF;
  text-align: center;
  justify-content: center;
}

.left-btn, .right-btn {
  position: absolute;
  top: 50%;
  left: 15px;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 1px black solid;
  background-color: black;
  color: #FFFFFF;
  text-align: center;
  justify-content: center;
  cursor: pointer;
  transform: rotate(180deg);

}

.right-btn {
  position: absolute;
  left: 92%;
  transform: rotate(0deg);
}

.viewers {
  position: absolute;
  left: 15%;
  top: 90%;
  width: 200px;
  height: auto;
  background-color: #FFFFFF;

  border: #000000 solid 1px;

  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;

  cursor: pointer;
  transition: 0.3s;
}

.viewers:hover {
  transform: scale(1.1);
  transition: 0.3s;
}

.close-friends-div {
  position: absolute;
  top: 5%;
  right: 12%;
  width: 200px;
  height: auto;
  background-color: #36c400;
  color: #FFFFFF;

  border: #000000 solid 1px;

  border-start-end-radius: 10px 10px;
  border-end-end-radius: 10px 10px;
  border-start-start-radius: 10px 10px;
  border-end-start-radius: 10px 10px;

}

</style>