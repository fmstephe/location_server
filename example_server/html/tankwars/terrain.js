function generateTerrain(width, height) {
	var wave = makeWave(width, width/r(10), width/r(2), height/r(5))
	for (var i = 0; i < 10; i++) {
		var divisor = r(i)
		var oWave = makeWave(width, width/r(10), width/(divisor*2.5), height/(divisor*5))
		addWave(wave, oWave);
	}
	normaliseWave(wave, 700, 50)
	return wave
}

function flatTerrain(width, height) {
	var terrain = new Array(width);
	for (var i = 0; i < width; i++) {
		terrain[i] = height/2;
	}
	return terrain;
}


function makeWave(len, offset, waveWidth, waveHeight) {
	var wave = new Array(len)
	var widthMult = (Math.PI*2)/waveWidth
	for (i = 0; i < len; i++) {
		x = (i+offset) * widthMult
		y = Math.floor(Math.cos(x)*waveHeight)
		wave[i] = y
	}
	return wave
}

function addWave(wave1, wave2) {
	for (i = 0; i < wave1.length; i++) {
		wave1[i] = wave1[i]+wave2[i]
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
		var inMagnitude = (wave[i]-min) / (max-min)
		wave[i] = (inMagnitude * outMagnitude) + lowerLim
	}
}
