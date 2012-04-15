function mkTankGame() {
	// Global Constants
	var host = "178.79.176.206";
	var maxPower = 200;
	var minPower = 0;
	var initPower = 50;
	var powerInc = 1;
	var health = 50;
	var gravity = 8;
	var turretLength = 20;
	var rotateInc = Math.PI/50;
	var frameRate = 30;
	var framePause = Math.floor(1000/frameRate);
	var expDuration = 0.1*frameRate;
	var expRadius = 60;

	// Canvas elements
	var tankCtxt;
	var missileCtxt;
	var terrainCtxt;
	var bgCtxt;

	var tankCanvas;
	var missileCanvas;
	var terrainCanvas;
	var bgCanvas;

	var gameHeight;
	var gameWidth;

	// Network services
	var connect;

	// Location
	var lat;
	var lng;
	// terrain
	var terrain;
	// Game entities
	var playerMe, playerYou;
	var keybindings;
	var playerList;
	var missileList;
	var explosionList;
	var wind;
	var windDiff;
	//
	var turnHandler;
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

	var completeFun = function(turnMsgs) {
		var wind = turnMsgs.filter(function(m) {return (typeof m.Content.data) == "number"});
		var players = turnMsgs.filter(function(m) {return m.Content.data.isPlayer;});
		return wind.size == 1 && players.size == 2;
	}

	return {
		init : function(height, width, idMe, idYou, nickMe, nickYou, xPosMe, xPosYou, initWind, cnct, divs, goFun) {
			       console.log(idMe);
			       console.log(nickMe);
			       console.log(idYou);
			       console.log(nickYou);
			       document.onkeydown = captureKeydown;
			       document.onkeyup = captureKeyup;
			       window.onresize = scaleCanvas;
			       var completeInit = function() {
				       gameHeight = height;
				       gameWidth = width;
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
				       gameOverFun = goFun;
				       var tanks = new LinkedList();
				       tanks.append(xPosMe);
				       tanks.append(xPosYou);
				       terrain = new Terrain(gameWidth, gameHeight, divs, tanks);
				       keybindings = new KeyBindings(87,83,65,68,70);
				       playerMe = new Player(idMe, xPosMe, nickMe, turretLength, initPower, minPower, maxPower, powerInc, rotateInc, health, false);
				       playerYou = new Player(idYou, xPosYou, nickYou, turretLength, initPower, minPower, maxPower, powerInc, rotateInc, health, true);
				       explosionList = new LinkedList();
				       missileList = new LinkedList();
				       playerList = new LinkedList();
				       playerList.append(playerMe);
				       playerList.append(playerYou);
				       wind = new Wind(gameWidth, gameHeight, 20);
				       wind.wind = initWind;
				       initRender();
				       loopId = setInterval(loop, framePause);
			       };
			       turnHandler = new TurnHandler(completeFun, idMe, idYou);
			       connect = cnct;
			       connect.addMsgHandler(turnHandler);
			       connect.sync(idMe, idYou, completeInit);
		       },
		     kill : function() {
				    killPrivate();
			    }
	}

	function killPrivate() {
		clearInterval(loopId);
		tankCtxt.clearRect(0, 0, gameWidth, gameHeight);
		missileCtxt.clearRect(0, 0, gameWidth, gameHeight);
		terrainCtxt.clearRect(0, 0, gameWidth, gameHeight);
		tankCtxt.clearRect(0, 0, gameWidth, gameHeight);
		missileCtxt.clearRect(0, 0, gameWidth, gameHeight);
		terrainCtxt.clearRect(0, 0, gameWidth, gameHeight);
		bgCtxt.clearRect(0, 0, gameWidth, gameHeight);
		connect.rmvMsgHandler(turnHandler);
		gameOverFun();
	}

	function initRender() {
		tankCanvas.width = gameWidth;
		tankCanvas.height = gameHeight;
		missileCanvas.width = gameWidth;
		missileCanvas.height = gameHeight;
		terrainCanvas.width = gameWidth;
		terrainCanvas.height = gameHeight;
		bgCanvas.width = gameWidth;
		bgCanvas.height = gameHeight;
		scaleCanvas();
		tankCtxt.fillStyle = "rgba(255,30,40,1.0)";
		tankCtxt.strokeStyle = "rgba(255,255,255,1.0)";
		tankCtxt.lineWidth = 5;
		missileCtxt.strokeStyle = "rgba(255,255,255,0.1)";
		missileCtxt.fillStyle = "rgba(255,0,0,1.0)";
		missileCtxt.lineWidth = 3;
		terrainCtxt.fillStyle = "rgba(100,100,100,1.0)";
		bgCtxt.fillStyle = "rgba(0,0,0,1.0)";
		bgCtxt.fillRect(0, 0, gameWidth, gameHeight);
		bgCtxt.strokeStyle = "rgba(255,255,255,1.0)";
		bgCtxt.lineWidth = 3;
		bgCtxt.globalAlpha = 0.3;
		terrain.render(terrainCtxt);
		terrain.clearMods();
	}

	function loop() {
		manageInfo();
		logInfo();
		if (!gameOver) {
			// Clear out each of the last frame's positions
			playerList.forEach(function(p) {p.setClear(tankCtxt, gameHeight);});
			missileList.forEach(function(m) {m.setClear(missileCtxt, gameHeight);});
			explosionList.forEach(function(e) {e.setClear(missileCtxt, gameHeight);});
			wind.setClear(bgCtxt);
			// Filter removable elements from entity lists
			playerList.filter(function(p) {return p.shouldRemove();});
			var filtered = missileList.filter(function(m) {return m.shouldRemove();});
			explosionList.filter(function(e) {return e.shouldRemove();});
			// Manage game entities
			playerList.forEach(function(p) {updatePlayer(p);});
			missileList.forEach(function(m) {updateMissile(m);});
			explosionList.forEach(function(e) {updateExplosion(e);});
			if (turnHandler.isComplete()) {
				var turn = turnHandler.getTurn();
				turn.forEach(function(m) {if (m.Content.data.id == playerYou.id) playerYou.arc = m.Content.data.arc});
				turn.forEach(function(m) {if (m.Content.data.isPlayer) launchMissile(m.Content.data);});
				turn.forEach(function(m) {if ((typeof m.Content.data) == "number") windDiff = m.Content.data});
			}
			// If an explosion has caused the terrain to change clear out the affected region
			terrain.setClear(terrainCtxt);
			// Render game entities
			terrain.render(terrainCtxt);
			playerList.forEach(function(p){p.render(tankCtxt, gameHeight)});
			if (filtered.size > 0 && missileList.size == 0) {
				playerMe.active = true;
				missileCtxt.clearRect(0,0,missileCanvas.width,missileCanvas.height);
				wind.addWind(windDiff);
			} else {
				missileList.forEach(function(m){m.render(missileCtxt, gameHeight)});
			}
			wind.advance();
			wind.render(bgCtxt);
			explosionList.forEach(function(e){e.render(missileCtxt, gameHeight)});
			terrain.clearMods();
			if (playerList.length() < 2) {
				setTimeout(killPrivate, 3000);
				gameOver = true;
			}
		} else {
			missileList.clear();
			explosionList.clear();
			tankCtxt.font = "50px Calibri-bold";
			var name = playerList.getFirst().name;
			tankCtxt.clearRect(0,0,gameWidth,gameHeight);
			tankCtxt.fillText(name + " wins!", gameWidth/2, gameHeight/2);
		}
	}

	function updatePlayer(player) {
		hr = terrain.heightArray;
		player.y = hr[player.x];
		if (player.y <= 0) {
			player.health = 0;
		}
		if (player.active) {
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
				if (playerMe.id > playerYou.id) {
					var windMsg = new TurnMsg(turnHandler.turnCount, wind.windChange());
					connect.sendMsg(idMe, windMsg);
					connect.sendMsg(idYou, windMsg);
				}
				player.active = false;
				var playerMsg = new TurnMsg(turnHandler.turnCount, player);
				connect.sendMsg(idMe, playerMsg);
				connect.sendMsg(idYou, playerMsg);
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
		if (missile.finalFrame) {
			missile.removed = true;
			return;
		}
		hr = terrain.heightArray;
		missile.advance(wind.wind);
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
		if (missile.x > gameWidth || missile.x < 0) {
			missile.finalFrame = true;
		}
		if (missile.y <= 0) {
			explodeMissile(missile,endX,0);
		}
	}

	function explodeMissile(missile, x, y) {
		exp = new Explosion(x, y, expDuration, expRadius);
		causeDamage(exp);
		explosionList.append(exp); 
		missile.x = x;
		missile.y = y;
		missile.finalFrame = true;
		return;
	}

	function updateExplosion(explosion) {
		explosion.deplete();
	}

	function causeDamage(explosion) {
		hr = terrain.heightArray;
		var x = Math.floor(explosion.x);
		var y = Math.floor(explosion.y);
		for (i = 0; i < expRadius; i++) {
			var sub = Math.sqrt((expRadius*expRadius)-(i*i));
			var bottom = y-sub;
			if (x+i < gameWidth && bottom < hr[x+i]) {
				hr[x+i] -= Math.min(sub*2, hr[x+i]-bottom);
			}
			if (x-i >= 0 && i != 0 && bottom < hr[x-i]) {
				hr[x-i] -= Math.min(sub*2, hr[x-i]-bottom);
			}
		}
		terrain.notifyMod(x-expRadius, x+expRadius);
		playerList.forEach(function(p) {damagePlayer(explosion, p)});
	}

	function damagePlayer(explosion, player) {
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
		var viewWidth = Math.min(window.innerWidth - 10, gameWidth);
		var viewHeight = Math.min(window.innerHeight - 10, gameHeight);
		var ratioHeight = viewWidth * (gameHeight/gameWidth);
		var ratioWidth = viewHeight * (gameWidth/gameHeight);
		if (ratioHeight > viewHeight) {
			viewWidth = ratioWidth;
		} else {
			viewHeight = ratioHeight;
		}
		if (viewWidth > gameWidth && viewHeight > gameHeight) {
			viewWidth = gameWidth;
			viewHeight = gameHeight;
		}
		tankCanvas.style.width = viewWidth;
		tankCanvas.style.height = viewHeight;
		missileCanvas.style.width = viewWidth;
		missileCanvas.style.height = viewHeight;
		terrainCanvas.style.width = viewWidth;
		terrainCanvas.style.height = viewHeight;
		bgCanvas.style.width = viewWidth;
		bgCanvas.style.height = viewHeight;
		$('game-div').style.height = viewHeight;
		$('game-div').style.width = viewWidth;
	}
}
