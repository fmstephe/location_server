function Missile(power, arc, x, y, gravity) {
	this.pushX = (power*Math.sin(arc));
	this.pushY = (power*Math.cos(arc));
	this.x = x;
	this.y = y;
	this.pX = this.x;
	this.pY = this.y;
	this.gravity = gravity;
	this.removed = false;
}

Missile.prototype.setClear = function(ctxt, hgt) {
	var x = Math.min(this.pX,this.x)-10;
	var y = hgt - (Math.max(this.pY,this.y)+10);
	var width = Math.abs(this.pX-this.x)+20;
	var h = Math.abs(this.pY-this.y)+20;
	ctxt.clearRect(x,y,width,h);
}

Missile.prototype.remove = function() {
	this.removed = true;
}

Missile.prototype.shouldRemove = function() {
	return this.removed;
}

Missile.prototype.render = function(ctxt, hgt) {
	if (!this.removed) {
		var pX = this.pX;
		var pY = hgt - this.pY;
		var x = this.x;
		var y = hgt - this.y;
		ctxt.strokeStyle = ctxt.createLinearGradient(Math.floor(pX),Math.floor(pY),Math.floor(x),Math.floor(y));
		ctxt.strokeStyle.addColorStop(0,"rgba(255,255,255,0.1)");
		ctxt.strokeStyle.addColorStop(1,"rgba(255,255,255,1)");
		//ctxt.strokeStyle = "rgba(255,255,255,1.0)";
		ctxt.beginPath();
		ctxt.moveTo(pX,pY);
		ctxt.lineTo(x,y);
		ctxt.closePath();
		ctxt.stroke();
	}
}

Missile.prototype.advance = function() {
	this.pX = this.x;
	this.pY = this.y;
	this.x += this.pushX;
	this.pushY -= this.gravity;
	this.y += this.pushY;
}
