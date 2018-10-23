<template>
  <div class="card">
      <div class="card-header">
        {{title}}
      </div>
      <div class="card-body">    
        <ul class="list-group list-group-flush">
          <li class="list-group-item" v-for="(message, key) in messages" :key = "key">
            {{message}}
          </li>
        </ul>
      </div>
  </div>
</template>

<script>
export default {
  name: 'MessageList',
  props: {
    title: String,
    messages: Array
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
