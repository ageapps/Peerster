<template>
  <div class="card">
      <div class="card-header">
        {{title}}
      </div>
      <div class="card-body">    
        <ul class="list-group list-group-flush">
          <li class="list-group-item" v-for="(message, key) in messages" :key = "key">
          <div class="card">
            <div class="card-header bg-dark text-white">
            {{message.origin}}
            <span class="badge badge-secondary">ID: {{message.id}}</span>
            </div>
            <div class="card-body">
             {{message.text}}
            </div>
          </div>
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
