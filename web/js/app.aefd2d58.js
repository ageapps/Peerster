(function(e){function t(t){for(var n,r,o=t[0],l=t[1],c=t[2],u=0,p=[];u<o.length;u++)r=o[u],i[r]&&p.push(i[r][0]),i[r]=0;for(n in l)Object.prototype.hasOwnProperty.call(l,n)&&(e[n]=l[n]);d&&d(t);while(p.length)p.shift()();return a.push.apply(a,c||[]),s()}function s(){for(var e,t=0;t<a.length;t++){for(var s=a[t],n=!0,o=1;o<s.length;o++){var l=s[o];0!==i[l]&&(n=!1)}n&&(a.splice(t--,1),e=r(r.s=s[0]))}return e}var n={},i={app:0},a=[];function r(t){if(n[t])return n[t].exports;var s=n[t]={i:t,l:!1,exports:{}};return e[t].call(s.exports,s,s.exports,r),s.l=!0,s.exports}r.m=e,r.c=n,r.d=function(e,t,s){r.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:s})},r.r=function(e){"undefined"!==typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},r.t=function(e,t){if(1&t&&(e=r(e)),8&t)return e;if(4&t&&"object"===typeof e&&e&&e.__esModule)return e;var s=Object.create(null);if(r.r(s),Object.defineProperty(s,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var n in e)r.d(s,n,function(t){return e[t]}.bind(null,n));return s},r.n=function(e){var t=e&&e.__esModule?function(){return e["default"]}:function(){return e};return r.d(t,"a",t),t},r.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},r.p="/";var o=window["webpackJsonp"]=window["webpackJsonp"]||[],l=o.push.bind(o);o.push=t,o=o.slice();for(var c=0;c<o.length;c++)t(o[c]);var d=l;a.push([0,"chunk-vendors"]),s()})({0:function(e,t,s){e.exports=s("56d7")},"0d31":function(e,t,s){"use strict";var n=s("ad9b"),i=s.n(n);i.a},"0d5c":function(e,t,s){},"157d":function(e,t,s){"use strict";var n=s("2652"),i=s.n(n);i.a},"1cf0":function(e,t,s){"use strict";var n=s("614f"),i=s.n(n);i.a},"1d79":function(e,t,s){"use strict";var n=s("72db"),i=s.n(n);i.a},2652:function(e,t,s){},"38c3":function(e,t,s){},"3ea3":function(e,t,s){"use strict";var n=s("b21f"),i=s.n(n);i.a},"4cbb":function(e,t,s){},"4f3f":function(e,t,s){},"51aa":function(e,t,s){"use strict";var n=s("5e9c"),i=s.n(n);i.a},"52a4":function(e,t,s){"use strict";var n=s("d20f"),i=s.n(n);i.a},"56d7":function(e,t,s){"use strict";s.r(t);s("cadf"),s("551c"),s("097d");var n=s("2b0e"),i=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"container-fluid",attrs:{id:"app"}},[s("h1",{staticClass:"title"},[e._v("Peerster UI")]),s("transition",{attrs:{appear:"",name:"fade",mode:"out-in"}},[e.started?e._e():s("div",{staticClass:"col-md-8 col-sm-12 container-fluid",attrs:{id:"dashboard"}},[s("PeerForm",{on:{"create-peer":e.startPeer}})],1),e.started?s("div",{staticClass:"container-fluid",attrs:{id:"dashboard"}},[s("p",[e._v("Peer: "+e._s(e.name))]),s("p",[e._v("Address: "+e._s(e.address))]),s("button",{staticClass:"btn btn-outline-danger",attrs:{type:"button"},on:{click:e.deletePeer}},[e._v("Delete this Peer")]),s("div",{staticClass:"row"},[s("div",{staticClass:"col-md-6 col-sm-12"},[s("NewMessageInput",{attrs:{title:"Send Message"},on:{"new-message":e.onNewMessage}}),s("NewMessageInput",{attrs:{title:"Send Private Message",peers:e.hops,isprivate:""},on:{"new-message":e.onNewMessage}}),s("NewRequestInput",{attrs:{title:"Request MetaHash",peers:e.hops},on:{"new-request":e.onNewRequest}}),s("UploadInput",{attrs:{title:"Upload File"},on:{"upload-file":e.onUpload}}),s("PeerList",{attrs:{peers:e.nodes,title:"Peers connected"}}),s("MessageList",{attrs:{messages:e.messages,title:"Messages"}})],1),s("div",{staticClass:"col-md-6 col-sm-12"},[s("SimpleInput",{attrs:{title:"Add new peer",submittext:"Add"},on:{"new-input":e.onNewPeer}}),s("SimpleInput",{attrs:{title:"Keyboards: txt,file...",submittext:"Search"},on:{"new-input":e.onNewSearch}}),s("FileList",{attrs:{files:e.files,baseurl:e.baseurl,title:"Files"}}),s("HopsList",{attrs:{peers:e.hops,title:"Hops"}}),s("PrivateList",{attrs:{messages:e.privateMesages,title:"Private Messages"}})],1)])]):e._e()]),s("transition",{attrs:{name:"fade"}},[e.alerting?s("div",{staticClass:"alert alert-danger",attrs:{role:"alert"}},[e._v("\n      "+e._s(e.alertText)+"\n    ")]):e._e()])],1)},a=[],r=(s("7f7f"),function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"hello"},[s("h1",[e._v(e._s(e.msg))])])}),o=[],l={name:"HelloWorld",props:{msg:String}},c=l,d=(s("65ef"),s("2877")),u=Object(d["a"])(c,r,o,!1,null,"e44aa914",null);u.options.__file="HelloWorld.vue";var p=u.exports,f=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"card"},[s("div",{staticClass:"card-body"},[s("div",{staticClass:"input-group mb-3"},[s("div",{directives:[{name:"show",rawName:"v-show",value:e.isprivate,expression:"isprivate"}],staticClass:"input-group-prepend"},[s("button",{staticClass:"btn btn-outline-secondary dropdown-toggle",attrs:{type:"button","data-toggle":"dropdown","aria-haspopup":"true","aria-expanded":"false"},on:{click:e.toggleDropdown}},[e._v(e._s(e.selectedHop))]),s("div",{staticClass:"dropdown-menu",style:{display:e.dropdownDisplay}},e._l(e.peers,function(t,n){return s("a",{key:n,staticClass:"dropdown-item",on:{click:function(t){e.selectNode(n)}}},[e._v("\n          "+e._s(n)+"\n          ")])}))]),s("input",{directives:[{name:"model",rawName:"v-model",value:e.newMessage,expression:"newMessage"}],staticClass:"form-control",attrs:{type:"text",placeholder:e.title,"aria-label":"Recipient's username","aria-describedby":"button-addon2"},domProps:{value:e.newMessage},on:{keyup:function(t){return"button"in t||!e._k(t.keyCode,"enter",13,t.key,"Enter")?e.trigger(t):null},input:function(t){t.target.composing||(e.newMessage=t.target.value)}}}),s("div",{staticClass:"input-group-append"},[s("button",{ref:"sendButton",staticClass:"btn btn-outline-secondary",attrs:{type:"button",id:"button-addon2"},on:{click:e.addMessage}},[e._v("Send")])])])])])},m=[],v={name:"NewMessageInput",props:{title:String,peers:Object,isprivate:{type:Boolean,default:!1}},methods:{selectNode:function(e){this.selectedHop=e,this.destination=e,this.toggleDropdown()},toggleDropdown:function(){this.dropdownDisplay="none"==this.dropdownDisplay?"block":"none"},trigger:function(){this.$refs.sendButton.click()},addMessage:function(){if(this.newMessage){var e={message:this.newMessage};this.isprivate&&(e.destination=this.destination),this.$emit("new-message",e),this.newMessage="",this.destination=""}}},data:function(){return{newMessage:"",destination:"",dropdownDisplay:"none",selectedHop:"Peers"}}},h=v,g=(s("51aa"),Object(d["a"])(h,f,m,!1,null,"17d7f8f0",null));g.options.__file="NewMessageInput.vue";var b=g.exports,_=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"card"},[s("div",{staticClass:"card-body"},[s("div",{staticClass:"input-group mb-3"},[s("div",{staticClass:"input-group"},[s("button",{staticClass:"btn btn-outline-secondary dropdown-toggle",attrs:{type:"button","data-toggle":"dropdown","aria-haspopup":"true","aria-expanded":"false"},on:{click:e.toggleDropdown}},[e._v(e._s(e.selectedHop))]),s("div",{staticClass:"dropdown-menu",style:{display:e.dropdownDisplay}},e._l(e.peers,function(t,n){return s("a",{key:n,staticClass:"dropdown-item",on:{click:function(t){e.selectNode(n)}}},[e._v("\n          "+e._s(n)+"\n          ")])}))]),s("input",{directives:[{name:"model",rawName:"v-model",value:e.requestHash,expression:"requestHash"}],staticClass:"form-control",attrs:{type:"text",placeholder:e.title,"aria-label":"Recipient's username","aria-describedby":"button-addon2"},domProps:{value:e.requestHash},on:{input:function(t){t.target.composing||(e.requestHash=t.target.value)}}}),s("div",{staticClass:"input-group"},[s("input",{directives:[{name:"model",rawName:"v-model",value:e.fileName,expression:"fileName"}],staticClass:"form-control",attrs:{type:"text",placeholder:"File name","aria-label":"File name","aria-describedby":"button-addon2"},domProps:{value:e.fileName},on:{input:function(t){t.target.composing||(e.fileName=t.target.value)}}})]),s("button",{staticClass:"btn btn-outline-secondary",attrs:{type:"button",id:"button-addon2"},on:{click:e.sendRequest}},[e._v("Send")])])])])},w=[],y={name:"NewRequestInput",props:{title:String,peers:Object},methods:{selectNode:function(e){this.selectedHop=e,this.destination=e,this.toggleDropdown()},toggleDropdown:function(){this.dropdownDisplay="none"==this.dropdownDisplay?"block":"none"},sendRequest:function(){if(this.fileName&&this.requestHash&&this.destination){var e={fileName:this.fileName,requestHash:this.requestHash,destination:this.destination};this.$emit("new-request",e),this.requestHash="",this.fileName="",this.destination=""}}},data:function(){return{fileName:"",requestHash:"",destination:"",dropdownDisplay:"none",selectedHop:"Peers"}}},C=y,P=(s("1cf0"),Object(d["a"])(C,_,w,!1,null,"77ee6f48",null));P.options.__file="NewRequestInput.vue";var x=P.exports,N=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"card"},[s("div",{staticClass:"card-body"},[s("div",{staticClass:"input-group mb-3"},[s("div",{staticClass:"input-group-prepend"},[s("button",{staticClass:"btn btn-outline-secondary",attrs:{type:"button",id:"button-addon2"},on:{click:function(t){e.toggleUpload()}}},[e._v("Select File")])]),s("input",{ref:"file",staticClass:"form-control",attrs:{hidden:"",type:"file",id:"file"},on:{change:function(t){e.handleFileUpload()}}}),s("input",{directives:[{name:"model",rawName:"v-model",value:e.fileName,expression:"fileName"}],staticClass:"form-control",attrs:{disabled:"",type:"text",placeholder:e.title},domProps:{value:e.fileName},on:{input:function(t){t.target.composing||(e.fileName=t.target.value)}}}),s("div",{staticClass:"input-group-append"},[s("button",{staticClass:"btn btn-outline-secondary",attrs:{type:"button",id:"button-addon2"},on:{click:function(t){e.submitFile()}}},[e._v("Upload")])])])])])},I=[],M={name:"UploadInput",props:{title:String},methods:{selectNode:function(e){this.selectedHop=e,this.destination=e,this.toggleDropdown()},toggleDropdown:function(){this.dropdownDisplay="none"==this.dropdownDisplay?"block":"none"},addMessage:function(){if(this.newMessage){var e={message:this.newMessage};this.isprivate&&(e.destination=this.destination),this.$emit("new-message",e),this.newMessage="",this.destination=""}},handleFileUpload:function(){this.file=this.$refs.file.files[0],this.fileName=this.file.name},toggleUpload:function(){var e=this.$refs.file;e.click()},submitFile:function(){var e={file:this.file};this.$emit("upload-file",e)}},data:function(){return{file:"",fileName:""}}},D=M,k=(s("951e"),Object(d["a"])(D,N,I,!1,null,"f1c36484",null));k.options.__file="UploadInput.vue";var S=k.exports,T=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"card"},[s("div",{staticClass:"card-body"},[s("div",{staticClass:"input-group mb-3"},[s("input",{directives:[{name:"model",rawName:"v-model",value:e.searchText,expression:"searchText"}],staticClass:"form-control",attrs:{type:"text",placeholder:e.title,"aria-label":"New peer address","aria-describedby":"button-addon2"},domProps:{value:e.searchText},on:{keyup:function(t){return"button"in t||!e._k(t.keyCode,"enter",13,t.key,"Enter")?e.trigger(t):null},input:function(t){t.target.composing||(e.searchText=t.target.value)}}}),s("div",{staticClass:"input-group-append"},[s("button",{ref:"sendButton",staticClass:"btn btn-primary",attrs:{type:"button",id:"button-addon2"},on:{click:e.addPeer}},[e._v(e._s(e.submittext))])])])])])},H=[],A={name:"SimpleInput",props:{title:String,submittext:String},methods:{addPeer:function(){this.searchText&&(this.$emit("new-input",this.searchText),this.searchText="")},trigger:function(){this.$refs.sendButton.click()}},data:function(){return{searchText:""}}},j=A,O=(s("157d"),Object(d["a"])(j,T,H,!1,null,"7f478132",null));O.options.__file="SimpleInput.vue";var $=O.exports,F=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"card"},[s("div",{staticClass:"card-header"},[e._v("\n      "+e._s(e.title)+"\n    ")]),s("div",{staticClass:"card-body"},[s("ul",{staticClass:"list-group"},e._l(e.peers,function(t,n){return s("li",{key:n,staticClass:"list-group-item"},[e._v("\n        "+e._s(t)+"\n        "),s("button",{directives:[{name:"show",rawName:"v-show",value:e.removable,expression:"removable"}],staticClass:"btn btn-outline-danger",attrs:{type:"button"},on:{click:function(t){e.removePeer(n)}}},[e._v("x")])])}))])])},L=[],q={name:"PeerList",props:{title:String,peers:Array,removable:{type:Boolean,default:!1}},methods:{removePeer:function(e){this.$emit("remove-peer",e)}}},E=q,U=(s("940c"),Object(d["a"])(E,F,L,!1,null,"c5ed597a",null));U.options.__file="PeerList.vue";var R=U.exports,G=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"card"},[s("div",{staticClass:"card-header"},[e._v("\n      "+e._s(e.title)+"\n    ")]),s("div",{staticClass:"card-body"},[s("ul",{staticClass:"list-group"},e._l(e.peers,function(t,n){return s("li",{key:n,staticClass:"list-group-item"},[e._v("\n          "+e._s(n)+" -> "+e._s(t.IP)+":"+e._s(t.Port)+"\n        ")])}))])])},B=[],W={name:"HopsList",props:{title:String,peers:Object}},V=W,J=(s("3ea3"),Object(d["a"])(V,G,B,!1,null,"3ab38e11",null));J.options.__file="HopsList.vue";var K=J.exports,Y=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"card"},[s("div",{staticClass:"card-header"},[e._v("\n      "+e._s(e.title)+"\n    ")]),s("div",{staticClass:"card-body"},[s("ul",{staticClass:"list-group"},e._l(e.files,function(t,n){return s("li",{key:n,staticClass:"list-group-item"},[e._v("\n          ["+e._s(n)+"]\n          "),s("a",{attrs:{target:"blanc",href:e.getUrl(t)}},[e._v(e._s(t))])])}))])])},z=[],Q={name:"FilesList",props:{title:String,files:Object,baseurl:String},methods:{getUrl:function(e){return this.baseurl+"/files/"+e}}},X=Q,Z=(s("0d31"),Object(d["a"])(X,Y,z,!1,null,"79e98fe6",null));Z.options.__file="FileList.vue";var ee=Z.exports,te=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"card"},[s("div",{staticClass:"card-header"},[e._v("\n      "+e._s(e.title)+"\n    ")]),s("div",{staticClass:"card-body"},[s("ul",{staticClass:"list-group"},e._l(e.messages,function(t,n){return s("li",{key:n,staticClass:"list-group-item"},[e._v("\n          Origin: "+e._s(n)+"\n                  "),s("ul",{staticClass:"list-group"},e._l(t,function(t,n){return s("li",{key:n,staticClass:"list-group-item"},[e._v("\n                        "+e._s(t.Destination)+" -> ["+e._s(t.ID)+"] "+e._s(t.Text)+"\n                    ")])}))])}))])])},se=[],ne={name:"HopsList",props:{title:String,messages:Object}},ie=ne,ae=(s("52a4"),Object(d["a"])(ie,te,se,!1,null,"4a92ac02",null));ae.options.__file="PrivateList.vue";var re=ae.exports,oe=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"card"},[s("div",{staticClass:"card-header"},[e._v("\n      "+e._s(e.title)+"\n    ")]),s("div",{staticClass:"card-body"},[s("ul",{staticClass:"list-group list-group-flush"},e._l(e.messages,function(t,n){return s("li",{key:n,staticClass:"list-group-item"},[s("div",{staticClass:"card"},[s("div",{staticClass:"card-header bg-dark text-white"},[e._v("\n          "+e._s(t.origin)+"\n          "),s("span",{staticClass:"badge badge-secondary"},[e._v("ID: "+e._s(t.id))])]),s("div",{staticClass:"card-body"},[e._v("\n           "+e._s(t.text)+"\n          ")])])])}))])])},le=[],ce={name:"MessageList",props:{title:String,messages:Array}},de=ce,ue=(s("1d79"),Object(d["a"])(de,oe,le,!1,null,"0692954a",null));ue.options.__file="MessageList.vue";var pe=ue.exports,fe=function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"jumbotron centered"},[s("h3",{staticClass:"display-5"},[e._v("Let's create a new peer!")]),s("div",{staticClass:"input-group mb-3"},[e._m(0),s("input",{directives:[{name:"model",rawName:"v-model",value:e.peerName,expression:"peerName"}],staticClass:"form-control",attrs:{type:"text",placeholder:"Peer Name","aria-label":"PeerName","aria-describedby":"basic-addon1"},domProps:{value:e.peerName},on:{input:function(t){t.target.composing||(e.peerName=t.target.value)}}})]),s("div",{staticClass:"input-group mb-3"},[e._m(1),s("input",{directives:[{name:"model",rawName:"v-model",value:e.peerAddress,expression:"peerAddress"}],staticClass:"form-control",attrs:{type:"text",placeholder:"Peer Address","aria-label":"PeerAddress","aria-describedby":"basic-addon1"},domProps:{value:e.peerAddress},on:{input:function(t){t.target.composing||(e.peerAddress=t.target.value)}}})]),s("hr",{staticClass:"my-4"}),s("p",[e._v("This will create a gossipper peer that will start gossiping with all the nodes in the network")]),s("PeerList",{directives:[{name:"show",rawName:"v-show",value:e.newPeers.length>0,expression:"newPeers.length>0"}],attrs:{removable:"",peers:e.newPeers,title:"Peers to connect"},on:{"remove-peer":e.removePeer}}),s("SimpleInput",{attrs:{title:"New peer",submittext:"Add",msg:"Welcome to Your Vue.js App"},on:{"new-input":e.onNewPeer}}),s("button",{staticClass:"btn btn-primary btn-lg",attrs:{disabled:e.isValid,type:"button"},on:{click:e.createPeer}},[e._v("Create")])],1)},me=[function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"input-group-prepend"},[s("span",{staticClass:"input-group-text text-white bg-primary"},[e._v("Name")])])},function(){var e=this,t=e.$createElement,s=e._self._c||t;return s("div",{staticClass:"input-group-prepend"},[s("span",{staticClass:"input-group-text text-white bg-primary"},[e._v("Address")])])}],ve={name:"PeerForm",props:{msg:String},components:{SimpleInput:$,PeerList:R},data:function(){return{newPeers:[],peerAddress:"127.0.0.1:5000",peerName:"nodeA"}},computed:{isValid:function(){return!this.peerAddress||!this.peerName}},methods:{createPeer:function(){this.$emit("create-peer",{name:this.peerName,address:this.peerAddress,peers:this.newPeers})},onNewPeer:function(e){this.newPeers.push(e)},removePeer:function(e){this.newPeers.splice(e,1)}}},he=ve,ge=(s("5fe6"),Object(d["a"])(he,fe,me,!1,null,"1b0eef1f",null));ge.options.__file="PeerForm.vue";var be=ge.exports,_e=s("bc3a"),we=s.n(_e),ye="http://localhost:8080",Ce=5,Pe={name:"app",components:{HelloWorld:p,NewMessageInput:b,NewRequestInput:x,UploadInput:S,SimpleInput:$,PeerList:R,HopsList:K,FileList:ee,PrivateList:re,MessageList:pe,PeerForm:be},data:function(){return{name:"nodeA",address:"0.0.0.0:0000",started:!1,messages:[],nodes:[],files:{},hops:{},privateMesages:{},loading:!1,alerting:!1,alertText:"",messageTimerID:"",peerTimerID:"",hopsTimerID:"",privateTimerID:"",filesTimerID:"",baseurl:ye}},methods:{onNewPeer:function(e){var t=this,s={name:this.name,node:e};we.a.post(ye+"/node",s).then(function(e){t.loading=!1,t.peers=e.data},function(e){t.loading=!1,t.showAlert(e.message)}),this.peers.push(e)},onNewSearch:function(e){var t=this,s={name:this.name,search:e};we.a.post(ye+"/search",s).then(function(){t.loading=!1},function(e){t.loading=!1,t.showAlert(e.message)})},onUpload:function(e){var t=this,s=new FormData;s.append("file",e.file),s.append("name",this.name),we.a.post(ye+"/upload",s,{headers:{"Content-Type":"multipart/form-data"}}).then(function(){t.showAlert("File uploaded!")}).catch(function(){t.showAlert("Error uploading the file!")})},onNewMessage:function(e){var t=this,s=ye+"/message",n={name:this.name,msg:e.message};e.destination&&(n.destination=e.destination,s=ye+"/private"),this.loading=!0,we.a.post(s,n).then(function(e){t.loading=!1,t.messages=e.data},function(e){t.loading=!1,t.showAlert(e.message)})},onNewRequest:function(e){var t=this,s=ye+"/request",n={name:this.name,file:e.fileName,destination:e.destination,hash:e.requestHash};this.loading=!0,we.a.post(s,n).then(function(){t.loading=!1},function(e){t.loading=!1,t.showAlert(e.message)})},startPeer:function(e){var t=this;this.loading=!0;var s={name:e.name,address:e.address};e.peers&&e.peers.length>0&&(s.peers=e.peers.join(",")),we.a.post(ye+"/start",s).then(function(s){t.loading=!1,t.peers=e.peers,t.name=s.data.name,t.address=s.data.address,t.started=!0},function(e){t.loading=!1,t.showAlert(e.message)}),this.startGettingMessages(),this.startGettingPrivateMessages(),this.startGettingPeers(),this.startGettingHops(),this.startGettingFiles()},deletePeer:function(){var e=this;this.loading=!0;var t={name:this.name};we.a.post(ye+"/delete",t).then(function(){e.exit()},function(t){e.loading=!1,e.showAlert(t.message)})},showAlert:function(e){var t=this;this.alerting=!0,this.alertText=e,setTimeout(function(){t.alerting=!1},4e3)},startGettingMessages:function(){var e=this;this.messageTimerID=setInterval(function(){e.getMessages()},1e3*Ce)},startGettingFiles:function(){var e=this;this.filesTimerID=setInterval(function(){e.getFiles()},1e3*Ce)},startGettingPrivateMessages:function(){var e=this;this.privateTimerID=setInterval(function(){e.getPrivateMessages()},1e3*Ce)},startGettingPeers:function(){var e=this;this.peerTimerID=setInterval(function(){e.getPeers()},1e3*Ce)},startGettingHops:function(){var e=this;this.hopsTimerID=setInterval(function(){e.getHops()},1e3*Ce)},getMessages:function(){var e=this,t=ye+"/message";this.getData(t,function(t){t&&(e.messages=t)})},getPrivateMessages:function(){var e=this,t=ye+"/private";this.getData(t,function(t){t&&(e.privateMesages=t)})},getPeers:function(){var e=this,t=ye+"/node";this.getData(t,function(t){t&&(e.nodes=t)})},getHops:function(){var e=this,t=ye+"/routes";this.getData(t,function(t){t&&(e.hops=t)})},getFiles:function(){var e=this,t=ye+"/files";this.getData(t,function(t){t&&(e.files=t)})},getData:function(e,t){var s=this,n={name:this.name};this.loading=!0,we.a.get(e,{params:n}).then(function(e){s.loading=!1,t(e.data)},function(e){s.loading=!1,s.exit(),s.showAlert(e.message),t(null)})},exit:function(){clearInterval(this.messageTimerID),clearInterval(this.peerTimerID),clearInterval(this.hopsTimerID),clearInterval(this.privateTimerID),clearInterval(this.filesTimerID),this.nodes=[],this.files={},this.messages=[],this.hops={},this.privateMesages={},this.name="",this.address="",this.started=!1}}},xe=Pe,Ne=(s("92c7"),Object(d["a"])(xe,i,a,!1,null,"6035f337",null));Ne.options.__file="App.vue";var Ie=Ne.exports;n["a"].config.productionTip=!1,new n["a"]({render:function(e){return e(Ie)}}).$mount("#app")},"5e9c":function(e,t,s){},"5fe6":function(e,t,s){"use strict";var n=s("38c3"),i=s.n(n);i.a},"614f":function(e,t,s){},"65ef":function(e,t,s){"use strict";var n=s("4cbb"),i=s.n(n);i.a},"72db":function(e,t,s){},"92c7":function(e,t,s){"use strict";var n=s("f158"),i=s.n(n);i.a},"940c":function(e,t,s){"use strict";var n=s("0d5c"),i=s.n(n);i.a},"951e":function(e,t,s){"use strict";var n=s("4f3f"),i=s.n(n);i.a},ad9b:function(e,t,s){},b21f:function(e,t,s){},d20f:function(e,t,s){},f158:function(e,t,s){}});
//# sourceMappingURL=app.aefd2d58.js.map