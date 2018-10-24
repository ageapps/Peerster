<template>
    <div class="card">
      <div class="card-body">
        <div class="input-group mb-3">
          <input type="text" class="form-control" placeholder="New Message" aria-label="Recipient's username" aria-describedby="button-addon2"  v-model="newMessage">
          <div class="input-group-append">
            <button class="btn btn-outline-secondary" @click="addMessage" type="button" id="button-addon2">Send</button>
          </div>
        </div>
      </div>
    </div>
</template>

<script>
export default {
  name: 'NewMessageInput',
  props: {
    msg: String
  },
  methods: {
    addMessage(){
      if(this.newMessage){
        this.$emit('new-message',this.newMessage)
        this.newMessage = ''
      }
    }
  },
  data(){
    return {
      newMessage: ''
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
