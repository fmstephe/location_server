// Global Constants
var host = "178.79.176.206";
var maxPower = 200;
var minPower = 0;
var initPower = 100;
var powerInc = 2;
var gravity = 12;
var turretLength = 20;
var rotationSpeed = Math.PI/50;
var frameRate = 30;
var framePause = Math.floor(1000/frameRate);
var expLife = 0.1*frameRate;
var expRadius = 50;

document.onkeydown = captureKeydown
document.onkeyup = captureKeyup

// Location
var lat;
var lng;
// Environment
var fgCtxt;
var terrainCtxt;
var bgCtxt;
var terrain;
var canvasHeight;
var canvasWidth;
// Game entity lists
var localPlayer;
var keyBindingList;
var launchList;
var playerList;
var missileList;
var explosionList;
// Current player whose turn it is
var gameOver = false;
// Frame rate tracking 
var lastCycle;
var thisCycle;
// Display toggle for nerdy info
var devMode;

//
var id;
// Location Service
var locService;
// Message Service
var msgService;

function getLocalCoords() {
	if (navigator.geolocation) {
		navigator.geolocation.getCurrentPosition(init,function(error) { console.log(JSON.stringify(error)), init({"coords": {"latitude":1, "longitude":1}}) });
	} else {
		alert("Get a real browser");
	}
}

function init(position) {
	lat = position.coords.latitude;
	lng = position.coords.longitude;
	id = getId();
	console.log("Id provided: " + id);
	initMsgService();
}

function initMsgService() {
	addMsg = new Add(id);
	msgService = new WSClient("Message", "ws://"+host+":8003/msg", handleMsg, function(){ 
		this.jsonsend(addMsg); 
		initLocService(); }, 
		function() {});
	msgService.connect();
}

function initLocService() {
	addMsg = new Add(id);
	initMsg = new InitLoc(lat, lng);
	locService = new WSClient("Location", "ws://"+host+":8002/loc", handleLoc, function(){ 
		this.jsonsend(addMsg); 
		this.jsonsend(initMsg);
       		initGame();}, 
		function() {});
	locService.connect();
}

function initGame() {
	devMode = false;
	lastCycle = new Date().getTime();
	thisCycle = new Date().getTime();
	var fgCanvas = document.getElementById("foreground");
	var terrainCanvas = document.getElementById("terrain");
	var bgCanvas = document.getElementById("background");
	fgCtxt = fgCanvas.getContext("2d");
	terrainCtxt = terrainCanvas.getContext("2d");
	bgCtxt = bgCanvas.getContext("2d");
	canvasHeight = fgCanvas.height;
	canvasWidth = fgCanvas.width;
	terrain = new Terrain(canvasWidth, canvasHeight);
	var kb1 = new KeyBindings(87,83,65,68,70);
	localPlayer = new Player(r(canvasWidth), "Player1", turretLength, initPower, minPower, maxPower, powerInc, expRadius, kb1);
	explosionList = new LinkedList();
	missileList = new LinkedList();
	launchList = new LinkedList();
	playerList = new LinkedList();
	playerList.append(localPlayer);
	keyBindingList = new LinkedList();
	keyBindingList.append(kb1);
	initRender();
	setInterval(loop, framePause);
}

function initRender() {
	fgCtxt.fillStyle = "rgba(255,30,40,1.0)";
	fgCtxt.strokeStyle = "rgba(255,255,255,1.0)";
	fgCtxt.lineWidth = 5;
	terrainCtxt.fillStyle = "rgba(100,100,100,1.0)";
	bgCtxt.fillStyle = "rgba(0,0,0,1.0)";
	bgCtxt.fillRect(0,0,canvasWidth,canvasHeight);
	terrain.render(terrainCtxt, canvasHeight);
	terrain.clearMods();
}

function loop() {
	manageInfo();
	logInfo();
	if (!gameOver) {
		// Clear out each of the last frame's positions
		playerList.forEach(function(p) {p.setClear(fgCtxt, canvasHeight);});
		missileList.forEach(function(m) {m.setClear(fgCtxt, canvasHeight);});
		explosionList.forEach(function(e) {e.setClear(fgCtxt, canvasHeight);});
		// Filter removable elements from entity lists
		playerList.filter(function(p) {return p.shouldRemove();});
		missileList.filter(function(m) {return m.shouldRemove();});
		explosionList.filter(function(e) {return e.shouldRemove();});
		// Manage game entities
		playerList.forEach(function(p) {updatePlayer(p);});
		missileList.forEach(function(m) {updateMissile(m);});
		explosionList.forEach(function(e) {updateExplosion(e);});
		if (launchList.size == playerList.size) {
			launchList.forEach(function(p) {launchMissile(p);});
			launchList.clear();
			keyBindingList.forEach(function(kb) {kb.reset();});
		}
		// If an explosion has caused the terrain to change clear out the affected region
		terrain.setClear(terrainCtxt, canvasHeight);
		// Render game entities
		terrain.render(terrainCtxt, canvasHeight);
		playerList.forEach(function(p){p.render(fgCtxt, canvasHeight)});
		missileList.forEach(function(m){m.render(fgCtxt, canvasHeight)});
		explosionList.forEach(function(e){e.render(fgCtxt, canvasHeight)});
		terrain.clearMods();
		//playerList.filter(function(p) {return p.health <= 0});
		if (playerList.length() < 1) {
			gameOver = true;
		}
	} else {
		missileList.clear();
		explosionList.clear();
		bgCtxt.font = "100pt Calibri-bold";
		var name = playerList.getFirst().name;
		fgCtxt.clearRect(0,0,canvasWidth,canvasHeight);
		fgCtxt.fillText(name + " wins!", canvasWidth/2, canvasHeight/2);
	}
}

function updatePlayer(player) {
	hr = terrain.heightArray;
	player.y = hr[player.x];
	if (!launchList.contains(player) && player.keyBindings != null) {
		if (player.keyBindings.left) {
			player.arc -= rotationSpeed;
		}
		if (player.keyBindings.right) {
			player.arc += rotationSpeed;
		}
		if (player.keyBindings.up) {
			player.incPower();
		}
		if (player.keyBindings.down) {
			player.decPower();
		}
		if (player.keyBindings.firing) {
			launchList.append(player);
		}
	}
}

function launchMissile(player) {
	missileList.append(new Missile(player, gravity));
}

function updateMissile(missile) {
	hr = terrain.heightArray;
	missile.advance();
	var startX = Math.floor(missile.pX);
	var endX = Math.floor(missile.x);
	var yD = missile.y - missile.pY;
	var startY = missile.pY;
	if (startX < endX) { // Travelling from left to right
		for (x = startX; x <= endX; x++) {
			yy = startY + (yD*((x-startX)/(endX-startX)));
			if (hr[x] > yy) {
				explodeMissile(missile,x,yy);
				return;
			}
		}
	} else if (endX < startX) { // Travelling from right to left
		for (x = startX; x >= endX; x--) {
			yy = startY + (yD*((startX-x)/(startX-endX)));
			if (hr[x] > yy) {
				explodeMissile(missile,x,yy);
				return;
			}
		}
	} else { // Travelling straight up and down
		if (missile.y < hr[startX]) {
			explodeMissile(missile,startX,hr[startX]);
			return;

		}
	}
	if (missile.x > canvasWidth || missile.x < 0 || missile.y < 0) {
		missle.remove();
	}
}

function explodeMissile(missile, x, y) {
	exp = new Explosion(x, y, expLife, expRadius);
	explode(exp);
	explosionList.append(exp); 
	missle.remove();
	return;
}

function updateExplosion(explosion) {
	explosion.deplete();
}

function explode(explosion) {
	hr = terrain.heightArray;
	var x = Math.floor(explosion.x);
	var y = Math.floor(explosion.y);
	for (i = 0; i < expRadius; i++) {
		var sub = Math.sqrt((expRadius*expRadius)-(i*i));
		var bottom = y-sub;
		if (x+i < canvasWidth && bottom < hr[x+i]) {
			hr[x+i] -= Math.min(sub*2, hr[x+i]-bottom);
		}
		if (x-i >= 0 && i != 0 && bottom < hr[x-i]) {
			hr[x-i] -= Math.min(sub*2, hr[x-i]-bottom);
		}
	}
	terrain.notifyMod(x-expRadius, x+expRadius);
	playerList.forEach(function(p) {explodePlayer(explosion, p)});
}

function explodePlayer(explosion, player) {
	var xD = player.x - explosion.x;
	var yD = player.y - explosion.y;
	var dist = Math.sqrt((xD*xD) + (yD*yD));
	if (dist < explosion.radius) {
		player.health -= explosion.radius - dist;
	}
}

function manageInfo() {
	lastCycle = thisCycle;
	thisCycle = new Date().getTime();
}

function logInfo() {
	if (devMode) {
		elapsed = thisCycle - lastCycle;
		frameRate = 1000/elapsed;
		console.log("Frame Rate: " + Math.floor(frameRate), "\tPlayers: " + playerList.length(), "\tMissiles: " + missileList.length(), "\tExplosions: " + explosionList.length());
	}
}

function r(lim) {
	return Math.floor(Math.random()*lim)+1
}

function captureKeydown(e) {
	var keyCode = e.keyCode;
	if (keyCode == 48) {
		devMode = !devMode;
		return;
	}	       
	keyBindingList.forEach(function(kb) {keydown(keyCode, kb);});
}

function keydown(keyCode, keyBinding) {
	if (keyBinding.upKey == keyCode) {
		keyBinding.up = true;
	}
	if (keyBinding.downKey == keyCode) {
		keyBinding.down = true;
	}
	if (keyBinding.leftKey == keyCode) {
		keyBinding.left = true;
	}
	if (keyBinding.rightKey == keyCode) {
		keyBinding.right = true;
	}
	if (keyBinding.firingKey == keyCode) {
		keyBinding.firing = true;
	}
}

function captureKeyup(e) {
	var keyCode = e.keyCode;
	keyBindingList.forEach(function(kb) {keyup(keyCode, kb);});
}

function keyup(keyCode, keyBinding) {
	if (keyBinding.upKey == keyCode) {
		keyBinding.up = false;
	}
	if (keyBinding.downKey == keyCode) {
		keyBinding.down = false;
	}
	if (keyBinding.leftKey == keyCode) {
		keyBinding.left = false;
	}
	if (keyBinding.rightKey == keyCode) {
		keyBinding.right = false;
	}
	if (keyBinding.firingKey == keyCode) {
		keyBinding.firing = false;
	}
}

