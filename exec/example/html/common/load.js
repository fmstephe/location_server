function $() {
	return document.getElementById(arguments[0]);
}

function sendMessages(cnct, id, times, cnt, max, pause) {
	if (cnt >= max) {
		formatResults(times, id);
		times.forEach(function(itm) {console.log(id + " " + itm)});
		return;
	}
	var content = new Object();
	content.sent = new Date().getTime();
	cnct.sendMsg(id, content);
	setTimeout(function() {sendMessages(cnct,id,times,cnt+1, max, pause)},pause);
}

function formatResults(times, id) {
	var str = "";
	times.forEach(function(itm) {str+=id+" ... "+itm+"\<br\>"});
	console.log(str);
	$("results").innerHTML += str;
}

function makeConnection(id, iterations, pause) {
	var locHandlers = new LinkedList();
	var msgHandlers = new LinkedList();
	var times = new LinkedList();
	msgHandlers.append(handleMsgFunc(times));
	var cnct = new Connect(id, msgHandlers, locHandlers, function(){}, function(){});
	setTimeout(function() {sendMessages(cnct, id, times, 0, iterations, pause)}, 1);
}

function handleMsgFunc(times) {
	return {
		handleMsg: function(msg) {
				   var sent = msg.content.sent;
				   var now = new Date().getTime();
				   times.append(now-sent);
			   }
	}
}

