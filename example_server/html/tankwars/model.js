// Interface for all model elements
//
// Used to clear out previous frame's drawing, hgt is the height of the context we are drawing to
// setClear(ctxt, hgt)
// Asks this entity to draw itself to the context provided, hgt is the height of the context we are drawing to
// render(ctxt, hgt)
//
//
function KeyBindings(upKey, downKey, leftKey, rightKey, firingKey) {
	this.upKey = upKey;
	this.downKey = downKey;
	this.leftKey = leftKey;
	this.rightKey = rightKey;
	this.firingKey = firingKey;
	this.up = false;
	this.down = false;
	this.left = false;
	this.right = false;
	this.firing = false;
	this.reset = resetKeyBindings;
}

function resetKeyBindings() {
	this.up = false;
	this.down = false;
	this.left = false;
	this.right = false;
	this.firing = false;
}

function Player(x, name, turretLength, initPower, minPower, maxPower, powerInc, health, keyBindings) {
	this.x = x;
	this.y = 0; // This gets set automatically by the physics
	this.name = name;
	this.arc = 0;
	this.power = initPower;
	this.minPower = minPower;
	this.maxPower = maxPower;
	this.powerInc = powerInc;
	this.health = health;
	this.turretLength = turretLength;
	this.keyBindings = keyBindings;
	this.incPower = incPowerPlayer;
	this.decPower = decPowerPlayer;
	this.setClear = setClearPlayer;
	this.shouldRemove = shouldRemovePlayer;
	this.render = renderPlayer;
}

function incPowerPlayer() {
	this.power += this.powerInc;
	this.power = Math.min(this.power, this.maxPower);
}

function decPowerPlayer() {
	this.power -= this.powerInc;
	this.power = Math.max(this.power, this.minPower);
}

function setClearPlayer(ctxt, hgt) {
	var x = this.x-this.turretLength;
	var y = hgt - (this.y + this.turretLength);
	var w = this.turretLength*6; // This is a cludge value to allow for clearing power % text
	var h = this.turretLength*2;
	ctxt.clearRect(x, y, w, h);
}

function shouldRemovePlayer() {
	return false;
}

// ctxt.fillStyle = "rgba(255,30,40,1.0)";
// ctxt.strokeStyle = "rgba(255,255,255,1.0)";
// ctxt.lineWidth = 5;
function renderPlayer(ctxt, hgt) {
	if (this.health > 0) {
		ctxt.beginPath();
		ctxt.arc(this.x, hgt-this.y, 10, 0, 2*Math.PI, true);
		ctxt.closePath();
		ctxt.fill();
		turretX = this.x+this.turretLength*Math.sin(this.arc);
		turretY = hgt-(this.y+(this.turretLength*Math.cos(this.arc)));
		ctxt.beginPath();
		ctxt.moveTo(this.x, hgt-this.y);
		ctxt.lineTo(turretX,turretY);
		ctxt.closePath();
		ctxt.stroke();
		var powerP = Math.round((this.power/maxPower)*100);
		ctxt.font = "20pt Calibri-bold";
		ctxt.fillText(powerP+"%",this.x+this.turretLength, hgt-this.y);
	}
}

function Missile(player, gravity) {
	this.pushX = (player.power*Math.sin(player.arc));
	this.pushY = (player.power*Math.cos(player.arc));
	this.x = player.x+(player.turretLength*Math.sin(player.arc));
	this.y = player.y+(player.turretLength*Math.cos(player.arc));
	this.pX = this.x;
	this.pY = this.y;
	this.player = player;
	this.gravity = gravity;
	this.removed = false;
	this.setClear = setClearMissile;
	this.remove = removeMissile;
	this.shouldRemove = shouldRemoveMissile;
	this.render = renderMissile;
	this.advance = advance;
}

function setClearMissile(ctxt, hgt) {
	var x = Math.min(this.pX,this.x)-10;
	var y = hgt - (Math.max(this.pY,this.y)+10);
	var width = Math.abs(this.pX-this.x)+20;
	var h = Math.abs(this.pY-this.y)+20;
	ctxt.clearRect(x,y,width,h);
}

function removeMissile() {
	this.removed = true;
}

function shouldRemoveMissile() {
	return this.removed;	
}

// ctxt.lineWidth = 5;
function renderMissile(ctxt, hgt) {
	if (!this.removed) {
		var pX = this.pX;
		var pY = hgt - this.pY;
		var x = this.x;
		var y = hgt - this.y;
		ctxt.strokeStyle = ctxt.createLinearGradient(Math.floor(pX),Math.floor(pY),Math.floor(x),Math.floor(y));
		ctxt.strokeStyle.addColorStop(0,"rgba(255,255,255,0.1)");
		ctxt.strokeStyle.addColorStop(1,"rgba(255,255,255,1)");
		ctxt.beginPath();
		ctxt.moveTo(pX,pY);
		ctxt.lineTo(x,y);
		ctxt.closePath();
		ctxt.stroke();
	}
}

function advance() {
	this.ppX = this.pX;
	this.ppY = this.pY;
	this.pX = this.x;
	this.pY = this.y;
	this.x += this.pushX;
	this.pushY -= this.gravity;
	this.y += this.pushY;
}

function Explosion(x, y, life, radius) {
	this.x = x;
	this.y = y;
	this.life = life;
	this.radius = radius;
	this.shouldRender = true;
	this.shouldRemove = false;
	this.setClear = setClearExplosion;
	this.deplete = depleteExplosion;
	this.shouldRemove = shouldRemoveExplosion;
	this.render = renderExplosion;
}

function setClearExplosion(ctxt, hgt) {
	var x = this.x - this.radius-2;
	var y = hgt - (this.y + this.radius + 2);
	var w = this.radius*2 + 4;
	var h = this.radius*2 + 4;
	ctxt.clearRect(x,y,w,h);
}

function depleteExplosion() {
	this.life--;
}

function shouldRemoveExplosion() {
	return this.life <= 0;
}

// fgCtxt.fillStyle = "rgba(255,30,30,1.0)";
function renderExplosion(ctxt, hgt) {
	var x = Math.floor(this.x);
	var y = Math.floor(this.y);
	ctxt.beginPath();
	ctxt.arc(x, hgt-y, this.radius, 0, 2*Math.PI, true);
	ctxt.closePath();
	ctxt.fill();
}

function Terrain(w, h) {
	this.heightArray = generateTerrain(w, h);
	this.w = w;
	this.h = h;
	this.regionList = new LinkedList();
	this.notifyMod = notifyModTerrain;
	this.clearMods = clearModsTerrain;
	this.setClear = setClearTerrain;
	this.render = renderTerrain;
	this.notifyMod(0,w);
}

function notifyModTerrain(from, to) {
	this.regionList.append(new Region(from,to));
}

function clearModsTerrain() {
	this.regionList.clear();
}

function setClearTerrain(ctxt, hgt) {
	this.regionList.forEach(function(r) {doClearTerrain(ctxt,r,hgt);});
}

function doClearTerrain(ctxt, region, hgt) {
	var x = region.from;
	var y = 0;
	var w = region.to - region.from;
	var h = hgt;
	ctxt.clearRect(x,y,w,h);
}

function renderTerrain(ctxt, hgt) {
	this.regionList.forEach(function(r) {doRenderTerrain(ctxt,r,hgt);});
}

// bgCtxt.fillStyle = "rgba(100,100,100,1.0)";
function doRenderTerrain(ctxt, region, hgt) {
	ctxt.beginPath();
	ctxt.moveTo(region.from,hgt);
	for (x = region.from; x <= region.to; x++) {
		ctxt.lineTo(x, hgt - terrain.heightArray[x]);
	}
	ctxt.lineTo(region.to,hgt);
	ctxt.closePath();
	ctxt.fill();
}

function Region(from, to) {
	this.from = from;
	this.to = to;
}
