function Terrain(width, height, divs, tanks) {
	this.heightArray = generateTerrain(width, height, divs);
	this.width = width;
	this.height = height;
	this.regionList = new LinkedList();
	this.notifyMod(0,width);
	var terrain = this;
	if (tanks) tanks.forEach(function(x) {terrain.flatten(x,10)});
}

Terrain.prototype.flatten = function(x,size) {
	var height = this.heightArray[x];
	var min = Math.max(0,x-size);
	var max = Math.min(this.width,x+size);
	for (var i = min; i <= max; i++) {
		this.heightArray[i] = height;
	}
}

Terrain.prototype.notifyMod = function(from, to) {
	this.regionList.append(new Region(from,to));
}

Terrain.prototype.clearMods = function() {
	this.regionList.clear();
}

Terrain.prototype.setClear = function(ctxt) {
	var terrain = this;
	this.regionList.forEach(function(region) {terrain.setClearRegion(ctxt, region);});
}

Terrain.prototype.setClearRegion = function(ctxt, region) {
	var x = region.from;
	var y = 0;
	var width = region.to - region.from;
	ctxt.clearRect(x,y,width,this.height);
}

Terrain.prototype.render = function(ctxt) {
	var trn = this;
	this.regionList.forEach(function(region) {trn.renderRegion(ctxt, region);});
}

Terrain.prototype.renderRegion = function(ctxt, region) {
	ctxt.beginPath();
	ctxt.moveTo(region.from, this.height);
	for (x = region.from; x <= region.to; x++) {
		ctxt.lineTo(x, this.height - this.heightArray[x]);
	}
	ctxt.lineTo(region.to, this.height);
	ctxt.closePath();
	ctxt.fill();
}

function Region(from, to) {
	this.from = from;
	this.to = to;
}

function generateTerrain(width, height, divs) {
	var wave;
	for (i in divs) {
		var offset = divs[i][0];
		var waveWidth = divs[i][1];
		var waveHeight = divs[i][2];
		var oWave = makeWave(width, width/offset, width/waveWidth, height/waveHeight);
		if (!wave) {
			wave = oWave;
		} else {
			addWave(wave, oWave);
		}
	}
	normaliseWave(wave, height*0.9, height*0.05);
	return wave
}

function genDivisors() {
	var divs = new Array();
	divs[0] = [r(10),r(2),r(5)];
	var widthConst = r(9)
		for (var i = 1; i < 10; i++) {
			var mult = r(i);
			var offset = r(10);
			var waveWidth = r(mult * widthConst);
			var waveHeight = r(mult * 5);
			divs[i] = [offset, waveWidth, waveHeight];
		}
	return divs;
}

function positionPair(width) {
	var one = Math.max(50, Math.floor((Math.random()*width)/3));
	var two = Math.min(width-50, Math.floor((Math.random()*width)/3 + 2*width/3));
	return [one,two];
}

function flatTerrain(width, height) {
	var terrain = new Array(width);
	for (var i = 0; i < width; i++) {
		terrain[i] = height/2;
	}
	return terrain;
}

function makeWave(len, offset, waveWidth, waveHeight) {
	var wave = new Array(len);
	var widthMult = (Math.PI*2)/waveWidth;
	for (i = 0; i < len; i++) {
		x = (i+offset) * widthMult;
		y = Math.floor(Math.cos(x)*waveHeight);
		wave[i] = y
	}
	return wave
}

function addWave(wave1, wave2) {
	for (i = 0; i < wave1.length; i++) {
		wave1[i] = wave1[i]+wave2[i];
	}
}

function normaliseWave(wave, upperLim, lowerLim) {
	var max = wave[0];
	for (i = 0; i < wave.length; i++) {
		if (wave[i] > max)
			max = wave[i];
	}
	var min = wave[0];
	for (i = 0; i < wave.length; i++) {
		if (wave[i] < min)
			min = wave[i];
	}
	var outMagnitude = upperLim - lowerLim;
	for (i = 0; i < wave.length; i++) {
		var inMagnitude = (wave[i]-min) / (max-min);
		wave[i] = (inMagnitude * outMagnitude) + lowerLim;
	}
}

function r(lim) {
	return Math.floor(Math.random()*lim)+1;
}
