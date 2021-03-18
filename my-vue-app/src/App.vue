<template>
  <div>
    <img alt="Vue logo" src="./assets/logo.png" />
    <div v-for="(message, index) in messages" :key="index">{{message.body}}</div>
    <div>
      <input type="text" v-model="msg">
      <button @click="sendMsg">Hit</button>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      socket: null,
      messages: [],
      msg: ''
    }
  },
  mounted() {
    this.socket = new WebSocket("ws://localhost:8080/ws");
    this.connect()
  },
  methods: {
    connect() {
      console.log("Attempting Connection...");

      this.socket.onopen = () => {
        console.log("Successfully Connected");
      };

      this.socket.onmessage = msg => {
        const message = JSON.parse(msg.data)
        console.log(message.body);
        this.messages.push(message)
      };

      this.socket.onclose = event => {
        console.log("Socket Closed Connection: ", event);
      };

      this.socket.onerror = error => {
        console.log("Socket Error: ", error);
      };
    },
    sendMsg() {
      this.socket.send(this.msg);
      this.msg = ''
    }
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>