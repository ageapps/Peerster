<template>
  <div id="app" class="container-fluid">
  <h1>Peerster UI</h1>
  <p>Peer: {{name}}</p>
  <div class="row">
    <div class="col-md-6 col-sm-12">
      <MessageForm @new-message="onNewMessage" msg="Welcome to Your Vue.js App"/>
    </div>
    <div class="col-md-6 col-sm-12">
      <PeerForm @new-peer="onNewPeer" msg="Welcome to Your Vue.js App"/>
    </div>
    </div>

    <div class="row">
      <div class="col-md-6 col-sm-12">
        <MessageList v-bind:messages="messages"  title="Messages"/>
      </div>
      <div class="col-md-6 col-sm-12">
        <PeerList v-bind:peers="peers"  title="Peers connected"/>
      </div>
    </div>
  </div>
</template>

<script>
import HelloWorld from './components/HelloWorld.vue'
import MessageForm from './components/MessageForm.vue'
import PeerForm from './components/PeerForm.vue'
import PeerList from './components/PeerList.vue'
import MessageList from './components/MessageList.vue'

export default {
  name: 'app',
  components: {
    HelloWorld,
    MessageForm,
    PeerForm,
    PeerList,
    MessageList
  },
  data() {
    return {
      name: 'nodeA',
      messages: [
        '172.0.0.1:4000',
        '172.0.0.1:4002',
        '172.0.0.1:4003',
      ],
      peers: [
        '172.0.0.1:4000',
        '172.0.0.1:4002',
        '172.0.0.1:4003',
      ]
    }
  },
  methods: {
    onNewPeer(data){
      // console.log("New Peer: " + data)
      this.peers.push(data)
    },
    onNewMessage(data){
      // console.log("New Message: " + data)
      this.messages.push(data)
    }
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
