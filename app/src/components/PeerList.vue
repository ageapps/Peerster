<template>
  <div class="card">
      <div class="card-header">
        {{title}}
      </div>
      <div class="card-body">    
        <ul class="list-group">
          <li class="list-group-item" v-for="(peer,key) in peers" :key = "key">
          {{peer}}
          </li>
        </ul>
      </div>
  </div>
</template>

<script>
export default {
  name: 'PeerList',
  props: {
    title: String,
    peers: Array
  },
  // data() {
  //   return {
  //   }
  // }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
