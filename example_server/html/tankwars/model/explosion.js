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

function renderExplosion(ctxt, hgt) {
	var x = Math.floor(this.x);
	var y = Math.floor(this.y);
	ctxt.beginPath();
	ctxt.arc(x, hgt-y, this.radius, 0, 2*Math.PI, true);
	ctxt.closePath();
	ctxt.fill();
}
