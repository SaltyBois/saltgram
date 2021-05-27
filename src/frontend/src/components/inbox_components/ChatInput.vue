<template>
  <div class="inner-inbox-content-footer">
    <EmojiPicker @emoji="append" :search="search">
      <div
          class="emoji-invoker"
          slot="emoji-invoker"
          slot-scope="{ events: { click: clickEvent } }"
          @click.stop="clickEvent">
        <svg height="24" viewBox="0 0 24 24" width="24" style="margin-top: 15px" xmlns="http://www.w3.org/2000/svg">
          <path d="M0 0h24v24H0z" fill="none"/>
          <path d="M11.99 2C6.47 2 2 6.48 2 12s4.47 10 9.99 10C17.52 22 22 17.52 22 12S17.52 2 11.99 2zM12 20c-4.42 0-8-3.58-8-8s3.58-8 8-8 8 3.58 8 8-3.58 8-8 8zm3.5-9c.83 0 1.5-.67 1.5-1.5S16.33 8 15.5 8 14 8.67 14 9.5s.67 1.5 1.5 1.5zm-7 0c.83 0 1.5-.67 1.5-1.5S9.33 8 8.5 8 7 8.67 7 9.5 7.67 11 8.5 11zm3.5 6.5c2.33 0 4.31-1.46 5.11-3.5H6.89c.8 2.04 2.78 3.5 5.11 3.5z"/>
        </svg>
      </div>

      <div slot="emoji-picker" slot-scope="{ emojis, insert }" style="z-index: 10;">
        <div class="emoji-picker" >
          <div class="emoji-picker__search">
            <input type="text" v-model="search" v-focus>
          </div>
          <div>
            <div v-for="(emojiGroup, category) in emojis" :key="category">
              <h5>{{ category }}</h5>
              <div class="emojis">
                <span
                    v-for="(emoji, emojiName) in emojiGroup"
                    :key="emojiName"
                    @click="insert(emoji)"
                    :title="emojiName"
                >{{ emoji }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </EmojiPicker>
    <i class="fa fa-photo picture-button" @click="$refs.file.click()"/>
    <input type="file"
           ref="file"
           style="display: none"
           @change="onSelectedFile"
           accept="image/*">
    <v-text-field label="Enter a message" style="padding-left: 5px; margin-bottom: 5px; min-width: 100%" />

    <div style="float: right; height: available; display: inline-block;">
      <v-btn class="follow-button" style="margin: 10px 5px 5px;width: 75px">
        Send
      </v-btn>
    </div>
  </div>
</template>

<script>
import EmojiPicker from "vue-emoji-picker";

export default {
  name: "ChatInput",
  components: {
    EmojiPicker
  },
  data: function () {
    return {
      search: '',
    }
  },
  methods: {
    onSelectedFile(event) {
      console.log(event)
      this.profilePicture = event.target.files[0]
      console.log(this.profilePicture)
    },
  }
}
</script>

<style scoped>

  .inner-inbox-content-footer {
    height: 60px;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
  }

  .emoji-invoker {
    position: relative;
    top: 0;
    left: 0;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    cursor: pointer;
    transition: all 0.2s;
  }
  .emoji-invoker:hover {
    transform: scale(1.1);
  }
  .emoji-invoker > svg {
    fill: black;
  }

  .emoji-picker {
    position: relative;
    border: 1px solid #707070;
    width: 250px;
    height: 200px;
    overflow: scroll;
    padding: 0;
    box-sizing: border-box;
    border-radius: 5%;
    background: #fff;
    box-shadow: 1px 1px 8px #c7dbe6;
  }
  .emoji-picker__search {
    display: flex;
  }
  .emoji-picker__search > input {
    flex: 1;
    border-radius: 10rem;
    border: 1px solid #ccc;
    padding: 0.5rem 1rem;
    outline: none;
  }
  .emoji-picker h5 {
    margin-bottom: 0;
    color: #b1b1b1;
    text-transform: uppercase;
    font-size: 0.8rem;
    cursor: default;
  }
  .emoji-picker .emojis {
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
  }
  .emoji-picker .emojis:after {
    content: "";
    flex: auto;
  }
  .emoji-picker .emojis span {
    padding: 0.2rem;
    cursor: pointer;
    border-radius: 5px;
  }
  .emoji-picker .emojis span:hover {
    background: #ececec;
    cursor: pointer;
  }

  .picture-button {
    transform: scale(1.5);
    margin-top: 24px;
    margin-right: 10px;
    cursor: pointer;
    transition: 0.3s;
  }

  .picture-button:hover {
    transform: scale(1.6);
    transition: 0.3s;
    color: #016ddb;
  }

  .follow-button {
    margin: 10px 0;
    width: 100px;
    height: 50px;
    background-color: transparent;
    color: #016ddb;
    border-color: #016ddb;
    border-style: solid;
    border-width: 1px;
    text-align: -webkit-center;
  }

</style>