(function(e){function t(t){for(var a,i,o=t[0],l=t[1],c=t[2],u=0,p=[];u<o.length;u++)i=o[u],r[i]&&p.push(r[i][0]),r[i]=0;for(a in l)Object.prototype.hasOwnProperty.call(l,a)&&(e[a]=l[a]);d&&d(t);while(p.length)p.shift()();return n.push.apply(n,c||[]),s()}function s(){for(var e,t=0;t<n.length;t++){for(var s=n[t],a=!0,o=1;o<s.length;o++){var l=s[o];0!==r[l]&&(a=!1)}a&&(n.splice(t--,1),e=i(i.s=s[0]))}return e}var a={},r={app:0},n=[];function i(t){if(a[t])return a[t].exports;var s=a[t]={i:t,l:!1,exports:{}};return e[t].call(s.exports,s,s.exports,i),s.l=!0,s.exports}i.m=e,i.c=a,i.d=function(e,t,s){i.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:s})},i.r=function(e){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},i.t=function(e,t){if(1&t&&(e=i(e)),8&t)return e;if(4&t&&"object"===typeof e&&e&&e.__esModule)return e;var s=Object.create(null);if(i.r(s),Object.defineProperty(s,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var a in e)i.d(s,a,function(t){return e[t]}.bind(null,a));return s},i.n=function(e){var t=e&&e.__esModule?function(){return e["default"]}:function(){return e};return i.d(t,"a",t),t},i.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},i.p="/";var o=window["webpackJsonp"]=window["webpackJsonp"]||[],l=o.push.bind(o);o.push=t,o=o.slice();for(var c=0;c<o.length;c++)t(o[c]);var d=l;n.push([0,"chunk-vendors"]),s()})({0:function(e,t,s){e.exports=s("56d7")},"0d5c":function(e,t,s){},"1d79":function(e,t,s){"use strict";var a=s("72db"),r=s.n(a);r.a},2594:function(e,t,s){"use strict";var a=s("6244"),r=s.n(a);r.a},"4cbb":function(e,t,s){},"53ce":function(e,t,s){"use strict";var a=s("e9e1"),r=s.n(a);r.a},"56d7":function(e,t,s){"use strict";s.r(t);s("cadf"),s("551c"),s("097d");var a=s("2b0e"),r=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"container-fluid",attrs:{id:"app"}},[s("h1",{staticClass:"title"},[e._v("Peerster UI")]),s("transition",{attrs:{appear:"",name:"fade",mode:"out-in"}},[e.started?e._e():s("div",{staticClass:"col-md-8 col-sm-12 container-fluid",attrs:{id:"dashboard"}},[s("PeerForm",{on:{"create-peer":e.startPeer}})],1),e.started?s("div",{staticClass:"container-fluid",attrs:{id:"dashboard"}},[s("p",[e._v("Peer: "+e._s(e.name))]),s("p",[e._v("Address: "+e._s(e.address))]),s("button",{staticClass:"btn btn-outline-danger",attrs:{type:"button"},on:{click:e.deletePeer}},[e._v("Delete this Peer")]),s("div",{staticClass:"row"},[s("div",{staticClass:"col-md-6 col-sm-12"},[s("NewMessageInput",{attrs:{msg:"Welcome to Your Vue.js App"},on:{"new-message":e.onNewMessage}})],1),s("div",{staticClass:"col-md-6 col-sm-12"},[s("NewPeerInput",{attrs:{msg:"Welcome to Your Vue.js App"},on:{"new-peer":e.onNewPeer}})],1)]),s("div",{staticClass:"row"},[s("div",{staticClass:"col-md-6 col-sm-12"},[s("MessageList",{attrs:{messages:e.messages,title:"Messages"}})],1),s("div",{staticClass:"col-md-6 col-sm-12"},[s("PeerList",{attrs:{peers:e.nodes,title:"Peers connected"}})],1)])]):e._e()]),s("transition",{attrs:{name:"fade"}},[e.alerting?s("div",{staticClass:"alert alert-danger",attrs:{role:"alert"}},[e._v("\n      "+e._s(e.alertText)+"\n    ")]):e._e()])],1)},n=[],i=(s("7f7f"),function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"hello"},[s("h1",[e._v(e._s(e.msg))])])}),o=[],l={name:"HelloWorld",props:{msg:String}},c=l,d=(s("65ef"),s("2877")),u=Object(d["a"])(c,i,o,!1,null,"e44aa914",null);u.options.__file="HelloWorld.vue";var p=u.exports,m=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"card"},[s("div",{staticClass:"card-body"},[s("div",{staticClass:"input-group mb-3"},[s("input",{directives:[{name:"model",rawName:"v-model",value:e.newMessage,expression:"newMessage"}],staticClass:"form-control",attrs:{type:"text",placeholder:"New Message","aria-label":"Recipient's username","aria-describedby":"button-addon2"},domProps:{value:e.newMessage},on:{input:function(t){t.target.composing||(e.newMessage=t.target.value)}}}),s("div",{staticClass:"input-group-append"},[s("button",{staticClass:"btn btn-outline-secondary",attrs:{type:"button",id:"button-addon2"},on:{click:e.addMessage}},[e._v("Send")])])])])])},v=[],g={name:"NewMessageInput",props:{msg:String},methods:{addMessage:function(){this.newMessage&&(this.$emit("new-message",this.newMessage),this.newMessage="")}},data:function(){return{newMessage:""}}},f=g,h=(s("8718"),Object(d["a"])(f,m,v,!1,null,"28c5b157",null));h.options.__file="NewMessageInput.vue";var b=h.exports,w=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"card"},[s("div",{staticClass:"card-body"},[s("div",{staticClass:"input-group mb-3"},[s("input",{directives:[{name:"model",rawName:"v-model",value:e.peerAdress,expression:"peerAdress"}],staticClass:"form-control",attrs:{type:"text",placeholder:"New peer name","aria-label":"New peer address","aria-describedby":"button-addon2"},domProps:{value:e.peerAdress},on:{input:function(t){t.target.composing||(e.peerAdress=t.target.value)}}}),s("div",{staticClass:"input-group-append"},[s("button",{staticClass:"btn btn-primary",attrs:{type:"button",id:"button-addon2"},on:{click:e.addPeer}},[e._v("Add")])])])])])},_=[],P={name:"NewPeerInput",props:{msg:String},methods:{addPeer:function(){this.peerAdress&&(this.$emit("new-peer",this.peerAdress),this.peerAdress="")}},data:function(){return{peerAdress:""}}},C=P,y=(s("f1a8"),Object(d["a"])(C,w,_,!1,null,"f7f7909c",null));y.options.__file="NewPeerInput.vue";var x=y.exports,A=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"card"},[s("div",{staticClass:"card-header"},[e._v("\n      "+e._s(e.title)+"\n    ")]),s("div",{staticClass:"card-body"},[s("ul",{staticClass:"list-group"},e._l(e.peers,function(t,a){return s("li",{key:a,staticClass:"list-group-item"},[e._v("\n        "+e._s(t)+"\n        "),s("button",{directives:[{name:"show",rawName:"v-show",value:e.removable,expression:"removable"}],staticClass:"btn btn-outline-danger",attrs:{type:"button"},on:{click:function(t){e.removePeer(a)}}},[e._v("x")])])}))])])},N=[],M={name:"PeerList",props:{title:String,peers:Array,removable:{type:Boolean,default:!1}},methods:{removePeer:function(e){this.$emit("remove-peer",e)}}},I=M,j=(s("940c"),Object(d["a"])(I,A,N,!1,null,"c5ed597a",null));j.options.__file="PeerList.vue";var O=j.exports,T=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"card"},[s("div",{staticClass:"card-header"},[e._v("\n      "+e._s(e.title)+"\n    ")]),s("div",{staticClass:"card-body"},[s("ul",{staticClass:"list-group list-group-flush"},e._l(e.messages,function(t,a){return s("li",{key:a,staticClass:"list-group-item"},[s("div",{staticClass:"card"},[s("div",{staticClass:"card-header bg-dark text-white"},[e._v("\n          "+e._s(t.origin)+"\n          "),s("span",{staticClass:"badge badge-secondary"},[e._v("ID: "+e._s(t.id))])]),s("div",{staticClass:"card-body"},[e._v("\n           "+e._s(t.text)+"\n          ")])])])}))])])},$=[],k={name:"MessageList",props:{title:String,messages:Array}},S=k,L=(s("1d79"),Object(d["a"])(S,T,$,!1,null,"0692954a",null));L.options.__file="MessageList.vue";var E=L.exports,D=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"jumbotron centered"},[s("h3",{staticClass:"display-5"},[e._v("Let's create a new peer!")]),s("div",{staticClass:"input-group mb-3"},[e._m(0),s("input",{directives:[{name:"model",rawName:"v-model",value:e.peerName,expression:"peerName"}],staticClass:"form-control",attrs:{type:"text",placeholder:"Peer Name","aria-label":"PeerName","aria-describedby":"basic-addon1"},domProps:{value:e.peerName},on:{input:function(t){t.target.composing||(e.peerName=t.target.value)}}})]),s("div",{staticClass:"input-group mb-3"},[e._m(1),s("input",{directives:[{name:"model",rawName:"v-model",value:e.peerAddress,expression:"peerAddress"}],staticClass:"form-control",attrs:{type:"text",placeholder:"Peer Address","aria-label":"PeerAddress","aria-describedby":"basic-addon1"},domProps:{value:e.peerAddress},on:{input:function(t){t.target.composing||(e.peerAddress=t.target.value)}}})]),s("hr",{staticClass:"my-4"}),s("p",[e._v("This will create a gossipper peer that will start gossiping with all the nodes in the network")]),s("PeerList",{directives:[{name:"show",rawName:"v-show",value:e.newPeers.length>0,expression:"newPeers.length>0"}],attrs:{removable:"",peers:e.newPeers,title:"Peers to connect"},on:{"remove-peer":e.removePeer}}),s("NewPeerInput",{attrs:{msg:"Welcome to Your Vue.js App"},on:{"new-peer":e.onNewPeer}}),s("button",{staticClass:"btn btn-primary btn-lg",attrs:{disabled:e.isValid,type:"button"},on:{click:e.createPeer}},[e._v("Create")])],1)},W=[function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"input-group-prepend"},[s("span",{staticClass:"input-group-text text-white bg-primary"},[e._v("Name")])])},function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"input-group-prepend"},[s("span",{staticClass:"input-group-text text-white bg-primary"},[e._v("Address")])])}],V={name:"PeerForm",props:{msg:String},components:{NewPeerInput:x,PeerList:O},data:function(){return{newPeers:[],peerAddress:"127.0.0.1:5000",peerName:"nodeA"}},computed:{isValid:function(){return!this.peerAddress||!this.peerName}},methods:{createPeer:function(){this.$emit("create-peer",{name:this.peerName,address:this.peerAddress,peers:this.newPeers})},onNewPeer:function(e){this.newPeers.push(e)},removePeer:function(e){this.newPeers.splice(e,1)}}},F=V,G=(s("2594"),Object(d["a"])(F,D,W,!1,null,"5f914ae5",null));G.options.__file="PeerForm.vue";var H=G.exports,Y=s("bc3a"),J=s.n(Y),B="http://localhost:8080",R={name:"app",components:{HelloWorld:p,NewMessageInput:b,NewPeerInput:x,PeerList:O,MessageList:E,PeerForm:H},data:function(){return{name:"nodeA",address:"0.0.0.0:0000",started:!1,messages:[],nodes:[],loading:!1,alerting:!1,alertText:"",messageTimerID:"",peerTimerID:""}},methods:{onNewPeer:function(e){var t=this,s={name:this.name,node:e};J.a.post(B+"/node",s).then(function(e){t.loading=!1,t.peers=e.data},function(e){t.loading=!1,t.showAlert(e.message)}),this.peers.push(e)},onNewMessage:function(e){var t=this,s={name:this.name,msg:e};this.loading=!0,J.a.post(B+"/message",s).then(function(e){t.loading=!1,t.messages=e.data},function(e){t.loading=!1,t.showAlert(e.message)})},startPeer:function(e){var t=this;this.loading=!0;var s={name:e.name,address:e.address};e.peers&&e.peers.length>0&&(s.peers=e.peers.join(",")),J.a.post(B+"/start",s).then(function(s){t.loading=!1,t.peers=e.peers,t.name=s.data.name,t.address=s.data.address,t.started=!0},function(e){t.loading=!1,t.showAlert(e.message)}),this.startGettingMessages(),this.startGettingPeers()},deletePeer:function(){var e=this;this.loading=!0;var t={name:this.name};J.a.post(B+"/delete",t).then(function(){e.exit()},function(t){e.loading=!1,e.showAlert(t.message)})},showAlert:function(e){var t=this;this.alerting=!0,this.alertText=e,setTimeout(function(){t.alerting=!1},4e3)},startGettingMessages:function(){var e=this;this.messageTimerID=setInterval(function(){e.getMessages()},3e3)},startGettingPeers:function(){var e=this;this.peerTimerID=setInterval(function(){e.getPeers()},3e3)},getMessages:function(){var e=this,t={name:this.name};this.loading=!0,J.a.get(B+"/message",{params:t}).then(function(t){e.loading=!1,e.messages=t.data},function(t){e.loading=!1,e.exit(),e.showAlert(t.message)})},getPeers:function(){var e=this,t={name:this.name};this.loading=!0,J.a.get(B+"/node",{params:t}).then(function(t){e.loading=!1,e.nodes=t.data},function(t){e.loading=!1,e.exit(),e.showAlert(t.message)})},exit:function(){clearInterval(this.messageTimerID),clearInterval(this.peerTimerID),this.peers=[],this.messages=[],this.name="",this.address="",this.started=!1}}},U=R,q=(s("53ce"),Object(d["a"])(U,r,n,!1,null,"1db4a9bc",null));q.options.__file="App.vue";var z=q.exports;a["a"].config.productionTip=!1,new a["a"]({render:function(e){return e(z)}}).$mount("#app")},6244:function(e,t,s){},"65ef":function(e,t,s){"use strict";var a=s("4cbb"),r=s.n(a);r.a},"72db":function(e,t,s){},8718:function(e,t,s){"use strict";var a=s("90cd"),r=s.n(a);r.a},"882d":function(e,t,s){},"90cd":function(e,t,s){},"940c":function(e,t,s){"use strict";var a=s("0d5c"),r=s.n(a);r.a},e9e1:function(e,t,s){},f1a8:function(e,t,s){"use strict";var a=s("882d"),r=s.n(a);r.a}});
//# sourceMappingURL=app.06841c2b.js.map