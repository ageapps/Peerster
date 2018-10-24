<template>
  <div class="jumbotron centered">
    <h3 class="display-5">Let's create a new peer!</h3>
    <div class="input-group mb-3">
      <div class="input-group-prepend">
        <span class="input-group-text text-white bg-primary">Name</span>
      </div>
      <input type="text" class="form-control" placeholder="Peer Name" aria-label="PeerName" aria-describedby="basic-addon1"  v-model="peerName">
    </div>
    <div class="input-group mb-3">
      <div class="input-group-prepend">
        <span class="input-group-text text-white bg-primary" >Address</span>
      </div>
      <input type="text" class="form-control" placeholder="Peer Address" aria-label="PeerAddress" aria-describedby="basic-addon1" v-model="peerAddress">
    </div>
    <hr class="my-4">
    <p>This will create a gossipper peer that will start gossiping with all the nodes in the network</p>
    <PeerList v-show="newPeers.length>0" removable @remove-peer="removePeer" v-bind:peers="newPeers"  title="Peers to connect"/>
    <NewPeerInput @new-peer="onNewPeer" msg="Welcome to Your Vue.js App"/>
    <button :disabled="isValid" class="btn btn-primary btn-lg" @click="createPeer" type="button">Create</button>
  </div>
</template>

<script>
import PeerList from './PeerList.vue'
import NewPeerInput from './NewPeerInput.vue'

export default {
  name: 'PeerForm',
  props: {
    msg: String
  },
  components: {
    NewPeerInput,
    PeerList,
  },
  data(){
    return {
      newPeers: [],
      peerAddress: '127.0.0.1:5000',
      peerName:'nodeA',
    }
  },
  computed: {
    isValid() {
      // evaluate whatever you need to determine disabled here...
      return !this.peerAddress || !this.peerName;
    }
  },
  methods: {
    createPeer(){
      
      this.$emit('create-peer',{
        name: this.peerName,
        address: this.peerAddress,
        peers: this.newPeers
      });
      // this.peerAddress = '';
      // this.peerName = '';
      // this.newPeers = [];
    },
    onNewPeer(data){
      this.newPeers.push(data)
    },
    removePeer(data){
      this.newPeers.splice(data,1)
    },
  },
  
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.centered{
    text-align: center;
}
.btn{
  margin-top: 20px;
}
</style>
