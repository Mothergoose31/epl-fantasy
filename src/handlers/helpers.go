package handlers

import (
	"math"
)

// hhamilton.typepad.com/files/pythag_mit_sa_2010.pdf
// https://blogs.salford.ac.uk/business-school/wp-content/uploads/sites/7/2016/09/paper.pdf
//

//  using Weibull distrubition first paper is describes win perentage as a function of goals scored and goals allowed

func PythagoreanExpectation(goalsScored, goalsAllowed float64, pythagoreanExponent float64) float64 {
	gamma := pythagoreanExponent
	kappa := math.Gamma(1 + 1/gamma)
	N := 10

	winProbability := math.Pow(goalsScored, gamma) / (math.Pow(goalsScored, gamma) + math.Pow(goalsAllowed, gamma))

	drawProbability := 0.0
	for c := 0; c <= N; c++ {
		cFloat := float64(c)
		gsProb := math.Exp(-math.Pow(kappa*(cFloat+1)/goalsScored, gamma)) - math.Exp(-math.Pow(kappa*cFloat/goalsScored, gamma))
		gaProb := math.Exp(-math.Pow(kappa*(cFloat+1)/goalsAllowed, gamma)) - math.Exp(-math.Pow(kappa*cFloat/goalsAllowed, gamma))
		drawProbability += gsProb * gaProb
	}

	expectedPoints := 3*winProbability + drawProbability

	return expectedPoints
}

//  create function to  account for weather at the stadium
