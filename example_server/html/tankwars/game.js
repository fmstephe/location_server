var terrainHeight = 640 * 2;
var terrainWidth = 960 * 2;

function mkTankGame() {
	// Global Constants
	var host = "178.79.176.206";
	var maxPower = 100;
	var minPower = 0;
	var initPower = 50;
	var powerInc = 1;
	var health = 100;
	var gravity = 4;
	var turretLength = 20;
	var rotateInc = Math.PI/50;
	var frameRate = 30;
	var framePause = Math.floor(1000/frameRate);
	var expDuration = 0.1*frameRate;
	var expRadius = 50;

	// Canvas elements
	var tankCtxt;
	var missileCtxt;
	var terrainCtxt;
	var bgCtxt;

	var tankCanvas;
	var missileCanvas;
	var terrainCanvas;
	var bgCanvas;

	var height;
	var width;
	var oldHeight;
	var oldWidth;

	// Network services
	var connect;

	// Location
	var lat;
	var lng;
	// terrain
	var terrain;
	// Game entity lists
	var playerMe, playerYou;
	var keybindings;
	var launchList;
	var playerList;
	var missileList;
	var explosionList;
	//
	var turnQ;
	// 
	var gameOver = false;
	//
	var gameOverFun;
	// Frame rate tracking 
	var lastCycle;
	var thisCycle;
	var loopId;
	// Display toggle for nerdy info
	var devMode = false;

	return {
		init : function(idMe, idYou, nickMe, nickYou, xPosMe, xPosYou, cnct, divs, turnQHandler, goFun) {
			       console.log(idMe);
			       console.log(nickMe);
			       console.log(idYou);
			       console.log(nickYou);
			       document.onkeydown = captureKeydown;
			       document.onkeyup = captureKeyup;
			       window.onresize = scaleCanvas;
			       turnQ = turnQHandler.q;
			       lastCycle = new Date().getTime();
			       thisCycle = new Date().getTime();
			       tankCanvas = document.getElementById("tank");
			       missileCanvas = document.getElementById("missile");
			       terrainCanvas = document.getElementById("terrain");
			       bgCanvas = document.getElementById("background");
			       tankCtxt = tankCanvas.getContext("2d");
			       missileCtxt = missileCanvas.getContext("2d");
			       terrainCtxt = terrainCanvas.getContext("2d");
			       bgCtxt = bgCanvas.getContext("2d");
			       connect = cnct;
			       gameOverFun = goFun;
			       terrain = new Terrain(terrainWidth, terrainHeight, divs);
			       keybindings = new KeyBindings(87,83,65,68,70);
			       playerMe = new Player(idMe, xPosMe, nickMe, turretLength, initPower, minPower, maxPower, powerInc, rotateInc, health, false);
			       playerYou = new Player(idYou, xPosYou, nickYou, turretLength, initPower, minPower, maxPower, powerInc, rotateInc, health, true);
			       explosionList = new LinkedList();
			       missileList = new LinkedList();
			       launchList = new LinkedList();
			       playerList = new LinkedList();
			       playerList.append(playerMe);
			       playerList.append(playerYou);
			       initRender();
			       loopId = setInterval(loop, framePause);
		       },
		     kill : function() {
				    killPrivate();
			    }
	}

	function killPrivate() {
		clearInterval(loopId);
		tankCtxt.clearRect(0, 0, terrainWidth, terrainHeight);
		missileCtxt.clearRect(0, 0, terrainWidth, terrainHeight);
		terrainCtxt.clearRect(0, 0, terrainWidth, terrainHeight);
		bgCtxt.clearRect(0, 0, terrainWidth, terrainHeight);
		gameOverFun();
	}

	function initRender() {
		tankCanvas.width = terrainWidth;
		tankCanvas.height = terrainHeight;
		missileCanvas.width = terrainWidth;
		missileCanvas.height = terrainHeight;
		terrainCanvas.width = terrainWidth;
		terrainCanvas.height = terrainHeight;
		bgCanvas.width = terrainWidth;
		bgCanvas.height = terrainHeight;
		scaleCanvas();
		tankCtxt.fillStyle = "rgba(255,30,40,1.0)";
		tankCtxt.strokeStyle = "rgba(255,255,255,1.0)";
		tankCtxt.lineWidth = 5;
		missileCtxt.strokeStyle = "rgba(255,255,255,1.0)";
		missileCtxt.fillStyle = "rgba(255,0,0,0.8)";
		missileCtxt.lineWidth = 5;
		terrainCtxt.fillStyle = "rgba(100,100,100,1.0)";
		bgCtxt.fillStyle = "rgba(0,0,0,1.0)";
		bgCtxt.fillRect(0, 0, terrainWidth, terrainHeight);
		terrain.render(terrainCtxt);
		terrain.clearMods();
	}

	function loop() {
		manageInfo();
		logInfo();
		if (!gameOver) {
			// Clear out each of the last frame's positions
			playerList.forEach(function(p) {p.setClear(tankCtxt, terrainHeight);});
			missileList.forEach(function(m) {m.setClear(missileCtxt, terrainHeight);});
			explosionList.forEach(function(e) {e.setClear(missileCtxt, terrainHeight);});
			// Filter removable elements from entity lists
			playerList.filter(function(p) {return p.shouldRemove();});
			missileList.filter(function(m) {return m.shouldRemove();});
			explosionList.filter(function(e) {return e.shouldRemove();});
			// Manage game entities
			playerList.forEach(function(p) {updatePlayer(p);});
			missileList.forEach(function(m) {updateMissile(m);});
			explosionList.forEach(function(e) {updateExplosion(e);});
			// Check the turnQ to see if new turn messages have arrived
			turnQ.forEach(function(msg) {if (msg.From == playerYou.id) launchList.append(msg.Content.player);});
			turnQ.clear();
			if (launchList.length() == playerList.size) {
				launchList.forEach(function(p) {launchMissile(p);});
				launchList.clear();
			}
			// If an explosion has caused the terrain to change clear out the affected region
			terrain.setClear(terrainCtxt);
			// Render game entities
			terrain.render(terrainCtxt);
			playerList.forEach(function(p){p.render(tankCtxt, terrainHeight)});
			missileList.forEach(function(m){m.render(missileCtxt, terrainHeight)});
			explosionList.forEach(function(e){e.render(missileCtxt, terrainHeight)});
			terrain.clearMods();
			if (playerList.length() < 2) {
				console.log(this);
				setTimeout(killPrivate, 3000);
				gameOver = true;
			}
		} else {
			missileList.clear();
			explosionList.clear();
			bgCtxt.font = "100pt Calibri-bold";
			var name = playerList.getFirst().name;
			tankCtxt.clearRect(0,0,terrainWidth,terrainHeight);
			tankCtxt.fillText(name + " wins!", terrainWidth/2, terrainHeight/2);
		}
	}

	function updatePlayer(player) {
		hr = terrain.heightArray;
		player.y = hr[player.x];
		if (player.y <= 0) {
			player.health = 0;
		}
		if (player === playerMe && missileList.length() == 0 && launchList.satAll(function(e){return player.id != e.id})) {
			player.animate = true;
			if (keybindings.left) {
				player.rotateLeft();
			}
			if (keybindings.right) {
				player.rotateRight();
			}
			if (keybindings.up) {
				player.incPower();
			}
			if (keybindings.down) {
				player.decPower();
			}
			if (keybindings.firing) {
				launchList.append(player);
				player.animate = false;
				connect.sendMsg(idYou, new PlayerMsg(player));
			}
		}
	}

	function launchMissile(playerMsg) {
		var power = playerMsg.power;
		var arc = playerMsg.arc;
		var x = playerMsg.x+(playerMsg.turretLength*Math.sin(arc));
		var y = playerMsg.y+(playerMsg.turretLength*Math.cos(arc));
		missileList.append(new Missile(power, arc, x, y, gravity));
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
		if (missile.x > terrainWidth || missile.x < 0) {
			missile.remove();
		}
		if (missile.y <= 0) {
			explodeMissile(missile,endX,0);
		}
	}

	function explodeMissile(missile, x, y) {
		exp = new Explosion(x, y, expDuration, expRadius);
		explode(exp);
		explosionList.append(exp); 
		missile.remove();
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
			if (x+i < terrainWidth && bottom < hr[x+i]) {
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
			//console.log("Frame Rate: " + Math.floor(frameRate), "\tPlayers: " + playerList.length(), "\tMissiles: " + missileList.length(), "\tExplosions: " + explosionList.length());
			console.log(launchList.length());
		}
	}

	function captureKeydown(e) {
		var keyCode = e.keyCode;
		if (keyCode == 48) {
			devMode = !devMode;
			return;
		}
		keydown(keyCode, keybindings);
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
		keyup(keyCode, keybindings);
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

	function scaleCanvas() {
		var viewWidth = window.innerWidth - 10;
		var viewHeight = window.innerHeight - 10;

		tankCanvas.style.width = viewWidth;
		tankCanvas.style.height = viewHeight;
		missileCanvas.style.width = viewWidth;
		missileCanvas.style.height = viewHeight;
		terrainCanvas.style.width = viewWidth;
		terrainCanvas.style.height = viewHeight;
		bgCanvas.style.width = viewWidth;
		bgCanvas.style.height = viewHeight;
	}
}
