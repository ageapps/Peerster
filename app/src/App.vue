<template>
  <div id="app" class="container-fluid">
    <h1 class="title">Peerster UI</h1>
    <transition appear name="fade" mode="out-in">
      <div v-if="!started" id="dashboard" class="col-md-8 col-sm-12 container-fluid">
        <PeerForm @create-peer="startPeer"/>
      </div>
      <div v-if="started" id="dashboard" class="container-fluid">
        <p>Peer: {{name}}</p>
        <p>Address: {{address}}</p>
        <button type="button" class="btn btn-outline-danger" @click="deletePeer" >Delete this Peer</button>
        <div class="row">
          <div class="col-md-6 col-sm-12">
            <NewMessageInput @new-message="onNewMessage" msg="Welcome to Your Vue.js App"/>
          </div>
          <div class="col-md-6 col-sm-12">
            <NewPeerInput @new-peer="onNewPeer" msg="Welcome to Your Vue.js App"/>
          </div>
        </div>

        <div class="row">
          <div class="col-md-6 col-sm-12">
            <MessageList v-bind:messages="messages"  title="Messages"/>
          </div>
          <div class="col-md-6 col-sm-12">
            <PeerList v-bind:peers="nodes"  title="Peers connected"/>
          </div>
        </div>
      </div>
    </transition>
    <transition name="fade">
      <div v-if="alerting" class="alert alert-danger" role="alert">
        {{alertText}}
      </div>
    </transition>
  </div>
</template>

<script>
import HelloWorld from './components/HelloWorld.vue'
import NewMessageInput from './components/NewMessageInput.vue'
import NewPeerInput from './components/NewPeerInput.vue'
import PeerList from './components/PeerList.vue'
import MessageList from './components/MessageList.vue'
import PeerForm from './components/PeerForm.vue'
import axios from 'axios';


var BACKEND_URL = "http://localhost:8080"
export default {
  name: 'app',
  components: {
    HelloWorld,
    NewMessageInput,
    NewPeerInput,
    PeerList,
    MessageList,
    PeerForm,
  },
  data() {
    return {
      name: 'nodeA',
      address: '0.0.0.0:0000',
      started: false,
      messages: [],
      nodes: [],
      loading: false,
      alerting: false,
      alertText: "",
      messageTimerID: "",
      peerTimerID: "",
    }
  },
  methods: {
    onNewPeer(data){
      // console.log("New Peer: " + data)
       var params = {
        name: this.name,
        node: data
      }
      axios.post(BACKEND_URL+"/node", params)
      .then((response)  =>  {
        this.loading = false;
        this.peers = response.data;
      }, (error)  =>  {
        this.loading = false;
        this.showAlert(error.message);
        return
      })      
      this.peers.push(data)
    },
    onNewMessage(data){
      // console.log("New Message: " + data)
      var params = {
        name: this.name,
        msg: data
      }
      this.loading = true;
      axios.post(BACKEND_URL+"/message", params)
      .then((response)  =>  {
        this.loading = false;
        this.messages = response.data;
      }, (error)  =>  {
        this.loading = false;
        this.showAlert(error.message);
        return
      })      
    },
    startPeer(data){
      this.loading = true;
      var params = {
        name: data.name,
        address: data.address
      }
      if(data.peers && data.peers.length > 0){
        params.peers=data.peers.join(',')
      }
      axios.post(BACKEND_URL+"/start", params)
      .then((response)  =>  {
        this.loading = false;
        this.peers = data.peers
        this.name = response.data.name
        this.address = response.data.address
        this.started = true
      }, (error)  =>  {
        this.loading = false;
        this.showAlert(error.message);
        return
      })
      this.startGettingMessages();    
      this.startGettingPeers();    
    },
    deletePeer(){
      this.loading = true;
      clearInterval(this.messageTimerID);
      clearInterval(this.peerTimerID);
      var params = {
        name: this.name,
      }
      console.log("deleting...")
      axios.post(BACKEND_URL+"/delete", params)
      .then((response)  =>  {
        console.log(response)
        this.peers = []
        this.name = ""
        this.address = ""
        this.started = false
      }, (error)  =>  {
        console.log(error)
        this.loading = false;
        this.showAlert(error.message);
        return
      })
    },
    showAlert(text){
      this.alerting = true;
      this.alertText = text;
      setTimeout(() => {  
        this.alerting = false;
      },4 * 1000);
    },
    startGettingMessages(){
      this.messageTimerID = setInterval(() => {  
        this.getMessages()
      },3 * 1000);
    },
    startGettingPeers(){
      this.peerTimerID = setInterval(() => {  
        this.getPeers()
      },3 * 1000);
    },
    getMessages(){
      var params = {
      name: this.name
      }
      this.loading = true;
      axios.get(BACKEND_URL+"/message", {
        params:params
      }).then((response)  =>  {
        this.loading = false;
        this.messages = response.data;
      }, (error)  =>  {
        this.loading = false;
        this.showAlert(error.message);
        return
      });    
    },
    getPeers(){
      var params = {
      name: this.name
      }
      this.loading = true;
      axios.get(BACKEND_URL+"/node", {
        params:params
      }).then((response)  =>  {
        console.log(response)
        this.loading = false;
        this.nodes = response.data;
      }, (error)  =>  {
        this.loading = false;
        this.showAlert(error.message);
        return
      });    
    }
  }
}
</script>

<style scoped>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  margin-top: 60px;
}
.title{
  text-align: center;
}
.btn{
  margin-bottom: 20px;
}
.fade-enter-active, .fade-leave-active {
  transition: opacity .5s;
}
.fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
  opacity: 0;
}
.alert{
  position: absolute;
  bottom: 0;
  left: 20px;
}
</style>
