Vue.component('peerster-message',{
  props: ['title','body'],

  template: `

  <article v-show='isVisible' class="message">
    <div class="message-header">

      <p>
      <slot name="header"></slot>
      </p>
    </div>
    <div class="message-body">
    <slot></slot>
    </div>
  </article>
  `,
  data(){
    return {
      isVisible: true
    }
  }
})


new Vue({
  el: '#app',
  data: {
    messages: [
      'Hello Vue.js!',
      'Hello 1',
      'Hello 2',
      'Hello 3'
    ],
    newMessage: '',
    isLoading: false,
    title: 'VUE'
  },
  methods: {
    addMessage(){
      this.messages.push(this.newMessage)
      this.newMessage = ''
      this.isLoading = !this.isLoading
    }
  },
  computed: {
    reverseTitle(){
      return this.title.split('').reverse().join('')
    }
  }
})
