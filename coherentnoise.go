package gah

import (
	opensimplex "github.com/ojrac/opensimplex-go"
)

// CoherentNoise provides automatic layering of opensimplex noise using parameters
type CoherentNoise struct {
	Noise       opensimplex.Noise // open simplex noise generator
	Scale       float64           // number that determines at what distance to view the noisemap, smaller is closer
	Octaves     int               // the number of levels of detail you want you perlin noise to have, higher gives more possible detail
	Lacunarity  float64           // number that determines how much detail is added or removed at each octave (adjusts frequency), higher gives less blending of octaves
	Persistence float64           // number that determines how much each octave contributes to the overall shape (adjusts amplitude), higher makes rougher
}

// NewCoherentNoise returns a CoherentNoise structure with the given parameters
func NewCoherentNoise(seed int64, scale float64, octaves int, lacunarity float64, persistence float64) *CoherentNoise {
	return &CoherentNoise{opensimplex.New(seed), scale, octaves, lacunarity, persistence}
}

// GetParamSignature returns a byte slice containing all relevant unique parameters
func (cnoise *CoherentNoise) GetParamSignature() (signature []byte) {
	signature = append(signature, Float64ToBytes(cnoise.Scale)...)
	signature = append(signature, IntToBytes(cnoise.Octaves)...)
	signature = append(signature, Float64ToBytes(cnoise.Lacunarity)...)
	signature = append(signature, Float64ToBytes(cnoise.Persistence)...)
	return signature
}

// GetEvalRange returns the min and max values that can be expected from the Eval2
func (cnoise *CoherentNoise) GetEvalRange() (outMin float64, outMax float64) {
	return -1, 1
}

//TODO any way to reduce redundancy here?

// eval1 works as Eval1 does on opensimplex.noise but applies layering through the CoherentNoise parameters
func (cnoise *CoherentNoise) Eval1(x float64) float64 {
	var maxAmp float64 = 0
	var amp float64 = 1
	var freq float64 = cnoise.Scale
	var noiseSample float64 = 0
	// add successively smaller, higher-frequency terms
	for i := 0; i < cnoise.Octaves; i++ {
		noiseSample += cnoise.Noise.Eval2(x*freq, 0) * amp
		maxAmp += amp
		amp *= cnoise.Persistence
		freq *= cnoise.Lacunarity
	}
	noiseSample /= maxAmp // take the average value of the iterations
	return noiseSample
}

// eval2 works as Eval2 does on opensimplex.noise but applies layering through the CoherentNoise parameters
func (cnoise *CoherentNoise) Eval2(x, y float64) float64 {
	var maxAmp float64 = 0
	var amp float64 = 1
	var freq float64 = cnoise.Scale
	var noiseSample float64 = 0
	// add successively smaller, higher-frequency terms
	for i := 0; i < cnoise.Octaves; i++ {
		noiseSample += cnoise.Noise.Eval2(x*freq, y*freq) * amp
		maxAmp += amp
		amp *= cnoise.Persistence
		freq *= cnoise.Lacunarity
	}
	noiseSample /= maxAmp // take the average value of the iterations
	return noiseSample
}

// eval3 works as Eval3 does on opensimplex.noise but applies layering through the CoherentNoise parameters
func (cnoise *CoherentNoise) Eval3(x, y, z float64) float64 {
	var maxAmp float64 = 0
	var amp float64 = 1
	var freq float64 = cnoise.Scale
	var noiseSample float64 = 0
	// add successively smaller, higher-frequency terms
	for i := 0; i < cnoise.Octaves; i++ {
		noiseSample += cnoise.Noise.Eval3(x*freq, y*freq, z*freq) * amp
		maxAmp += amp
		amp *= cnoise.Persistence
		freq *= cnoise.Lacunarity
	}
	noiseSample /= maxAmp // take the average value of the iterations
	return noiseSample
}

// eval4 works as Eval4 does on opensimplex.noise but applies layering through the CoherentNoise parameters
func (cnoise *CoherentNoise) Eval4(x, y, z, w float64) float64 {
	var maxAmp float64 = 0
	var amp float64 = 1
	var freq float64 = cnoise.Scale
	var noiseSample float64 = 0
	// add successively smaller, higher-frequency terms
	for i := 0; i < cnoise.Octaves; i++ {
		noiseSample += cnoise.Noise.Eval4(x*freq, y*freq, z*freq, w*freq) * amp
		maxAmp += amp
		amp *= cnoise.Persistence
		freq *= cnoise.Lacunarity
	}
	noiseSample /= maxAmp // take the average value of the iterations
	return noiseSample
}
