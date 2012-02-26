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
