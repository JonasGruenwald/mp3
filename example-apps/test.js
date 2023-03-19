console.log("App running, will log heartbeat every second")
setInterval(() => {
  console.log('heartbeat ', new Date().toLocaleString())
},1000)
