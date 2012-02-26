function Terrain(w, h, divs) {
	this.heightArray = generateTerrain(w, h, divs);
	this.w = w;
	this.h = h;
	this.regionList = new LinkedList();
	this.notifyMod = notifyModTerrain;
	this.clearMods = clearModsTerrain;
	this.setClear = setClearTerrain;
	this.render = renderTerrain;
	this.notifyMod(0,w);
}

Terrain.prototype.notifyMod = function(from, to) {
	this.regionList.append(new Region(from,to));
}

Terrain.prototype.clearMods = function() {
	this.regionList.clear();
}

Terrain.prototype.setClear = function(ctxt, hgt) {
	this.regionList.forEach(function(r) {doClearTerrain(ctxt,r,hgt);});
}

function doClearTerrain(ctxt, region, hgt) {
	var x = region.from;
	var y = 0;
	var w = region.to - region.from;
	var h = hgt;
	ctxt.clearRect(x,y,w,h);
}

Terrain.prototype.render = function(ctxt, hgt) {
	this.regionList.forEach(function(r) {doRenderTerrain(ctxt,r,hgt);});
}

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

function generateTerrain(width, height) {
	var divs = genDivisors();
	return generateTerrain(width, height, divs);
}

function generateTerrain(width, height, divs) {
	var wave;
	for (i in divs) {
		var offset = divs[i][0];
		var waveWidth = divs[i][1];
		var waveHeight = divs[i][2];
		var oWave = makeWave(width, width/(offset*1.5), width/(waveWidth*2.5), height/(waveHeight*5));
		if (!wave) {
			wave = oWave;
		} else {
			addWave(wave, oWave);
		}
	}
	normaliseWave(wave, height*0.7, height*0.2);
	return wave
}

function makeGameDef(width) {
	return {div: genDivisors(), pos: positionPair(width)};
}

function genDivisors() {
	var divs = new Array();
	for (var i = 0; i < 10; i++) {
		var offset = r(10);
		var waveWidth = r(2);
		var waveHeight = r(5);
		divs[i] = [offset, waveWidth, waveHeight];
	}
	return divs;
}

function positionPair(width) {
	var one = Math.floor((Math.random()*width)/3);
	var two = Math.floor((Math.random()*width)/3 + 2*width/3);
	console.log(one+two);
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
