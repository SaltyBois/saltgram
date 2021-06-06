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
        <v-hover v-slot="{ hover }">
          <v-carousel class="post-content-media" :cycle="!hover" :interval="10000" :continuous="false">
            <v-carousel-item v-for="(item, index) in contentPlaceHolder.length" :key="index">
              <v-img contain
                     v-if="contentPlaceHolder[index].endsWith('.jpg') || contentPlaceHolder[index].endsWith('.png') || contentPlaceHolder[index].endsWith('.jpeg')"
                     :src="contentPlaceHolder[index]"/>
              <video autoplay
                     playsinline
                     v-else-if="contentPlaceHolder[index].endsWith('.mp4')"
                     :src="contentPlaceHolder[index]"/>
            </v-carousel-item>
            <div class="close-friends-div" :v-if="true">
              <h3>
                CLOSE FRIENDS
              </h3>
            </div>
          </v-carousel>
        </v-hover>
        <div class="viewers">
          <h3>
            Seen by 503
          </h3>
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
      contentPlaceHolder: [
        'https://skinnyms.com/wp-content/uploads/2015/04/9-Best-Grumpy-Cat-Memes-750x500.jpg',
        'https://i.kym-cdn.com/entries/icons/original/000/035/692/cover1.jpg',
        'https://www.thehonestkitchen.com/blog/wp-content/uploads/2019/07/CatMemes-copy-10.jpg',
        'https://i.ytimg.com/vi/KHa4OOvYLx0/maxresdefault.jpg',
        'https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1543715662l/43075028._SX318_.jpg',
        'https://www.w3schools.com/html/movie.mp4'
      ]
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
        // this.timerSeconds = 10
        // this.countDownTimer()
      }
    },
    // countDownTimer() {
    //   if(this.timerSeconds > 0) {
    //     setTimeout(() => {
    //       this.timerSeconds -= 1
    //       this.countDownTimer()
    //     }, 1000)
    //   }
    //   else {
    //     this.visible = false;
    //   }
    // }
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
  left: 33%;
  top: 93%;
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
  top: 2%;
  right: 2%;
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