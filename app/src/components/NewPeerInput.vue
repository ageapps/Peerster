<template>
    <div class="card">
      <div class="card-body">
        <div class="input-group mb-3">
          <input type="text" class="form-control" placeholder="New peer name" aria-label="New peer address" aria-describedby="button-addon2" v-model="peerAdress">
          <div class="input-group-append">
            <button class="btn btn-primary" @click="addPeer" type="button" id="button-addon2">Add</button>
          </div>
        </div>
      </div>
    </div>
</template>

<script>
export default {
  name: 'NewPeerInput',
  props: {
    msg: String
  },
  methods: {
    addPeer(){
      if(this.peerAdress){
        // console.log(this.peerAdress)
        this.$emit('new-peer',this.peerAdress)
        this.peerAdress = ''
      }
    }
  },
  data(){
    return {
      peerAdress: ''
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.input-group{
  margin-bottom: 0px !important;
}
</style>
